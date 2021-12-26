package interpreter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/parser"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Interpreter struct {
	stdRoot          string
	scopes           map[string][]map[string]types.Type
	currentPackage   string
	parentPackage    string
	includedPackages map[string]types.Type
}

func NewInterpreter(stdRoot, currentPackage, parentPackage string) *Interpreter {
	return &Interpreter{
		stdRoot:          stdRoot,
		currentPackage:   currentPackage,
		parentPackage:    parentPackage,
		scopes:           map[string][]map[string]types.Type{},
		includedPackages: map[string]types.Type{},
	}
}

func (i *Interpreter) pushScope(packageName string, scope map[string]types.Type) {
	i.scopes[packageName] = append(i.scopes[packageName], scope)
}

func (i *Interpreter) popScope(packageName string) map[string]types.Type {
	if scopes, ok := i.scopes[packageName]; ok {
		if len(scopes) == 0 {
			panic("fatal: not enough scopes")
		}

		scope := scopes[len(scopes)-1]
		i.scopes[packageName] = scopes[:len(scopes)-1]
		return scope
	}

	panic(fmt.Sprintf("fatal: scopes for '%s' package does not exist", packageName))
}

func (i *Interpreter) getVar(packageName string, name string) (types.Type, error) {
	if scopes, ok := i.scopes[packageName]; ok {
		lastScopeIdx := len(scopes) - 1
		for idx := lastScopeIdx; idx >= 0; idx-- {
			if val, ok := scopes[idx][name]; ok {
				return val, nil
			}
		}

		return nil, util.RuntimeError(fmt.Sprintf("ідентифікатор '%s' не визначений", name))
	}

	panic(fmt.Sprintf("fatal: scopes for '%s' package does not exist", packageName))
}

func (i *Interpreter) setVar(packageName, name string, value types.Type) error {
	if scopes, ok := i.scopes[packageName]; ok {
		scopesLen := len(scopes)
		for idx := 0; idx < scopesLen; idx++ {
			if oldValue, ok := scopes[idx][name]; ok {
				if oldValue.GetTypeHash() != value.GetTypeHash() {
					if scopesLen == 1 {
						return util.RuntimeError(
							fmt.Sprintf(
								"неможливо записати значення типу '%s' у змінну '%s' з типом '%s'",
								value.GetTypeName(), name, oldValue.GetTypeName(),
							),
						)
					}

					// TODO: надрукувати нормальне попередження!
					fmt.Println(
						fmt.Sprintf(
							"Увага: несумісні типи даних '%s' та '%s', змінна '%s' стає недоступною в поточному полі видимості",
							value.GetTypeName(), oldValue.GetTypeName(), name,
						),
					)
					break
				}

				scopes[idx][name] = value
				i.scopes[packageName] = scopes
				return nil
			}
		}

		scopes[scopesLen-1][name] = value
		i.scopes[packageName] = scopes
		return nil
	}

	panic(fmt.Sprintf("fatal: scopes for '%s' package does not exist", packageName))
}

// TODO: pass parent type from what call operation was performed, nil otherwise;
//  this is useful for calling methods of custom classes.
func (i *Interpreter) executeNode(
	rootNode ast.ExpressionNode, rootDir string, thisPackage, parentPackage string,
) (types.Type, bool, error) {
	switch node := rootNode.(type) {
	case ast.ImportNode:
		res, err := i.executeImport(&node, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.FunctionDefNode:
		functionPackage := types.NewPackageInstance(
			false,
			thisPackage,
			parentPackage,
			map[string]types.Type{}, // TODO: set attributes
		)
		functionDef := types.NewFunctionInstance(
			node.Name.Text, node.Arguments,
			func(_ *[]types.Type, kwargs *map[string]types.Type) (types.Type, error) {
				res, _, err := i.executeBlock(*kwargs, node.Body, thisPackage, parentPackage)
				return res, err
			},
			node.ReturnType,
			functionPackage,
			"", // TODO: set doc
		)
		return functionDef, false, i.setVar(thisPackage, node.Name.Text, functionDef)

	case ast.ClassDefNode:
		classPackage := types.NewPackageInstance(
			false,
			thisPackage,
			parentPackage,
			map[string]types.Type{}, // TODO: set attributes
		)
		// TODO: set doc
		// TODO: set attributes

		attributes := map[string]types.Type{}
		for _, attributeNode := range node.Attributes {
			switch attribute := attributeNode.(type) {
			case ast.BinOperationNode:
				if attribute.Operator.Type.Name == models.Assign {
					switch variableNode := attribute.LeftNode.(type) {
					case ast.VariableNode:
						res, _, err := i.executeNode(attribute.RightNode, rootDir, thisPackage, parentPackage)
						if err != nil {
							return nil, false, err
						}

						attributes[variableNode.Variable.Text] = res
						continue
					}
				}
			case ast.FunctionDefNode:
				functionDef := types.NewFunctionInstance(
					attribute.Name.Text, attribute.Arguments,
					func(_ *[]types.Type, kwargs *map[string]types.Type) (types.Type, error) {
						res, _, err := i.executeBlock(*kwargs, attribute.Body, thisPackage, parentPackage)
						return res, err
					},
					attribute.ReturnType,
					nil,
					"", // TODO: set doc
				)
				attributes[attribute.Name.Text] = functionDef
			}

			_, _, err := i.executeNode(attributeNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}
		}

		classDef := types.NewClass(node.Name.Text, classPackage, attributes, node.Doc.Text)
		return classDef, false, i.setVar(thisPackage, node.Name.Text, classDef)

	case ast.CallOpNode:
		callable, err := i.getVar(thisPackage, node.CallableName.Text)
		if err != nil {
			return nil, false, err
		}

		res, err := i.executeCallOp(&node, callable, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.ReturnNode:
		res, _, err := i.executeNode(node.Value, rootDir, thisPackage, parentPackage)
		return res, true, err

	case ast.UnaryOperationNode:
		res, err := i.executeUnaryOp(&node, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.BinOperationNode:
		res, err := i.executeBinaryOp(&node, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.RandomAccessOperationNode:
		res, err := i.executeRandomAccessGetOp(node.Operand, node.Index, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.ListSlicingNode:
		res, err := i.executeListSlicing(&node, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.IfNode:
		return i.executeIfSequence(node.Blocks, node.ElseBlock, rootDir, thisPackage, parentPackage)

	case ast.ForEachNode:
		container, _, err := i.executeNode(node.Container, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, false, err
		}

		return i.executeForEachLoop(
			node.IndexVar, node.ItemVar, container, node.Body, thisPackage, parentPackage,
		)

	case ast.RealTypeNode:
		res, err := types.NewRealInstanceFromString(node.Value.Text)
		return res, false, err

	case ast.IntegerTypeNode:
		res, err := types.NewIntegerInstanceFromString(node.Value.Text)
		return res, false, err

	case ast.StringTypeNode:
		res := types.NewStringInstance(node.Value.Text)
		return res, false, nil

	case ast.BoolTypeNode:
		res, err := types.NewBoolInstanceFromString(node.Value.Text)
		return res, false, err

	case ast.ListTypeNode:
		list := types.NewListInstance()
		for _, valueNode := range node.Values {
			value, _, err := i.executeNode(valueNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			list.Values = append(list.Values, value)
		}

		return list, false, nil

	case ast.DictionaryTypeNode:
		dict := types.NewDictionaryInstance()
		for keyNode, valueNode := range node.Map {
			key, _, err := i.executeNode(keyNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			value, _, err := i.executeNode(valueNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			err = dict.SetElement(key, value)
			if err != nil {
				return nil, false, err
			}
		}

		return dict, false, nil

	case ast.NilTypeNode:
		return types.NewNilInstance(), false, nil

	case ast.VariableNode:
		val, err := i.getVar(thisPackage, node.Variable.Text)
		if err != nil {
			return nil, false, err
		}

		return val, false, nil

	case ast.AttrAccessOpNode:
		res, err := i.executeAttrAccessOp(&node, rootDir, thisPackage, parentPackage)
		return res, false, err
	}

	return nil, false, util.RuntimeError("невідома помилка")
}

func (i *Interpreter) executeAST(
	scope map[string]types.Type, thisPackage, parentPackage string, tree *ast.AST,
) (types.Type, map[string]types.Type, bool, error) {
	var filePath string
	var dir string
	var err error
	if thisPackage == "<стдввід>" {
		filePath = "<стдввід>"
		dir, err = os.Getwd()
		if err != nil {
			return nil, scope, false, util.InternalError(err.Error())
		}
	} else {
		filePath, err = filepath.Abs(thisPackage)
		if err != nil {
			return nil, scope, false, util.InternalError(err.Error())
		}

		dir = filepath.Dir(filePath)
	}

	var result types.Type = nil
	forceReturn := false
	i.pushScope(thisPackage, scope)
	for _, node := range tree.Nodes {
		result, forceReturn, err = i.executeNode(node, dir, thisPackage, parentPackage)
		if err != nil {
			return nil, scope, false, errors.New(
				fmt.Sprintf(
					"  Файл \"%s\", рядок %d\n    %s\n%s",
					filePath, node.RowNumber(), node.String(), err.Error(),
				),
			)
		}

		if forceReturn {
			break
		}
	}

	scope = i.popScope(thisPackage)
	return result, scope, forceReturn, nil
}

func (i *Interpreter) executeBlock(
	scope map[string]types.Type, tokens []models.Token, thisPackage, parentPackage string,
) (types.Type, bool, error) {
	p := parser.NewParser(thisPackage, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, false, err
	}

	result, _, forceReturn, err := i.executeAST(scope, thisPackage, parentPackage, asTree)
	if err != nil {
		return nil, false, err
	}

	return result, forceReturn, nil
}

func (i *Interpreter) Execute(
	thisPackage, parentPackage string, scope map[string]types.Type, code string,
) (types.Type, map[string]types.Type, error) {
	lexer := borsch.NewLexer(thisPackage, code)
	tokens, err := lexer.Lex()
	if err != nil {
		return nil, scope, err
	}

	p := parser.NewParser(thisPackage, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, scope, err
	}

	i.pushScope(thisPackage, builtin.RuntimeObjects)
	var result types.Type
	result, scope, _, err = i.executeAST(scope, thisPackage, parentPackage, asTree)
	if err != nil {
		return nil, scope, err
	}

	return result, scope, nil
}

func (i *Interpreter) ExecuteFile(
	packageName, parentPackage string, content []byte, isBuiltin bool,
) (types.Type, error) {
	_, scope, err := i.Execute(packageName, parentPackage, map[string]types.Type{}, string(content))
	if err != nil {
		return nil, err
	}

	return types.NewPackageInstance(isBuiltin, packageName, parentPackage, scope), nil
}
