package interpreter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/parser"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Interpreter struct {
	stdRoot          string
	scopes           map[string][]map[string]types.Type
	context          *Context
	includedPackages map[string]types.Type
}

func NewInterpreter(stdRoot, currentPackage, parentPackage string) *Interpreter {
	context := &Context{
		rootDir:           stdRoot,
		package_:          types.NewPackageInstance(false, currentPackage, parentPackage, map[string]types.Type{}),
		parentPackageName: parentPackage,
	}
	context.parentObject = context.package_
	return &Interpreter{
		stdRoot:          stdRoot,
		scopes:           map[string][]map[string]types.Type{},
		context:          context,
		includedPackages: map[string]types.Type{},
	}
}

func (i *Interpreter) GetContext() *Context {
	return i.context
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

func (i *Interpreter) setVar(package_ *types.PackageInstance, name string, value types.Type) error {
	if package_ == nil {
		// return errors.New("setVar: package is nil")
		return nil
	}

	if scopes, ok := i.scopes[package_.Name]; ok {
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
				i.scopes[package_.Name] = scopes
				return nil
			}
		}

		scopes[scopesLen-1][name] = value
		i.scopes[package_.Name] = scopes
		return nil
	}

	panic(fmt.Sprintf("fatal: scopes for '%s' package does not exist", package_.Name))
}

func (i *Interpreter) executeNode(ctx *Context, rootNode ast.ExpressionNode) (types.Type, bool, error) {
	switch node := rootNode.(type) {
	case ast.ImportNode:
		return nil, false, i.executeImport(ctx, &node)

	case ast.FunctionDefNode:
		functionDef := types.NewFunctionInstance(
			node.Name.Text, node.Arguments,
			func(_ *[]types.Type, kwargs *map[string]types.Type) (types.Type, error) {
				res, _, err := i.executeBlock(ctx, *kwargs, node.Body)
				return res, err
			},
			node.ReturnType,
			false,
			ctx.package_,
			node.Doc,
		)
		return functionDef, false, i.setVar(ctx.GetPackageFromParent(), node.Name.Text, functionDef)

	case ast.ClassDefNode:
		// TODO: set attributes

		attributes := map[string]types.Type{}
		classContext := ctx.WithParent(nil)
		for _, attributeNode := range node.Attributes {
			switch attribute := attributeNode.(type) {
			case ast.BinOperationNode:
				if attribute.Operator.Type.Name == models.Assign {
					switch variableNode := attribute.LeftNode.(type) {
					case ast.VariableNode:
						res, _, err := i.executeNode(classContext, attribute.RightNode)
						if err != nil {
							return nil, false, err
						}

						if _, ok := res.(*types.FunctionInstance); ok {
							// Error is muted.
							_, _ = res.SetAttribute(ops.DocAttributeName, types.NewStringInstance(variableNode.Doc))
						}

						attributes[variableNode.Variable.Text] = res
						continue
					}
				}

			case ast.FunctionDefNode:
				functionDef := types.NewFunctionInstance(
					attribute.Name.Text, attribute.Arguments,
					func(_ *[]types.Type, kwargs *map[string]types.Type) (types.Type, error) {
						res, _, err := i.executeBlock(ctx, *kwargs, attribute.Body)
						return res, err
					},
					attribute.ReturnType,
					true,
					nil,
					attribute.Doc,
				)
				attributes[attribute.Name.Text] = functionDef
			default:
				_, _, err := i.executeNode(classContext, attributeNode)
				if err != nil {
					return nil, false, err
				}
			}
		}

		classDef := types.NewClass(node.Name.Text, ctx.package_, attributes, node.Doc)
		return nil, false, i.setVar(ctx.GetPackageFromParent(), node.Name.Text, classDef)

	case ast.CallOpNode:
		callable, err := i.getVar(ctx.package_.Name, node.CallableName.Text)
		if err != nil {
			return nil, false, err
		}

		res, err := i.executeCallOp(ctx, &node, callable)
		return res, false, err

	case ast.ReturnNode:
		res, _, err := i.executeNode(ctx, node.Value)
		return res, true, err

	case ast.UnaryOperationNode:
		res, err := i.executeUnaryOp(ctx, &node)
		return res, false, err

	case ast.BinOperationNode:
		res, err := i.executeBinaryOp(ctx, &node)
		return res, false, err

	case ast.RandomAccessOperationNode:
		res, err := i.executeRandomAccessGetOp(ctx, node.Operand, node.Index)
		return res, false, err

	case ast.ListSlicingNode:
		res, err := i.executeListSlicing(ctx, &node)
		return res, false, err

	case ast.IfNode:
		return i.executeIfSequence(ctx, node.Blocks, node.ElseBlock)

	case ast.ForEachNode:
		container, _, err := i.executeNode(ctx, node.Container)
		if err != nil {
			return nil, false, err
		}

		return i.executeForEachLoop(ctx, node.IndexVar, node.ItemVar, container, node.Body)
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
			value, _, err := i.executeNode(ctx, valueNode)
			if err != nil {
				return nil, false, err
			}

			list.Values = append(list.Values, value)
		}

		return list, false, nil

	case ast.DictionaryTypeNode:
		dict := types.NewDictionaryInstance()
		for keyNode, valueNode := range node.Map {
			key, _, err := i.executeNode(ctx, keyNode)
			if err != nil {
				return nil, false, err
			}

			value, _, err := i.executeNode(ctx, valueNode)
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
		val, err := i.getVar(ctx.package_.Name, node.Variable.Text)
		if err != nil {
			return nil, false, err
		}

		return val, false, nil

	case ast.AttrAccessOpNode:
		res, err := i.executeAttrAccessOp(ctx, &node)
		return res, false, err
	}

	return nil, false, util.RuntimeError("невідома помилка")
}

func (i *Interpreter) executeAST(
	ctx *Context, scope map[string]types.Type, tree *ast.AST,
) (types.Type, map[string]types.Type, bool, error) {
	var filePath string
	var err error
	if ctx.package_.Name == builtin.RootPackageName {
		filePath = builtin.RootPackageName
		ctx.rootDir, err = os.Getwd()
		if err != nil {
			return nil, scope, false, util.InternalError(err.Error())
		}
	} else {
		filePath, err = filepath.Abs(ctx.package_.Name)
		if err != nil {
			return nil, scope, false, util.InternalError(err.Error())
		}

		ctx.rootDir = filepath.Dir(filePath)
	}

	var result types.Type = nil
	forceReturn := false
	i.pushScope(ctx.package_.Name, scope)
	for _, node := range tree.Nodes {
		result, forceReturn, err = i.executeNode(ctx, node)
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

	scope = i.popScope(ctx.package_.Name)
	return result, scope, forceReturn, nil
}

func (i *Interpreter) executeBlock(ctx *Context, scope map[string]types.Type, tokens []models.Token) (
	types.Type,
	bool,
	error,
) {
	p := parser.NewParser(ctx.package_.Name, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, false, err
	}

	result, _, forceReturn, err := i.executeAST(ctx, scope, asTree)
	if err != nil {
		return nil, false, err
	}

	return result, forceReturn, nil
}

func (i *Interpreter) Execute(ctx *Context, scope map[string]types.Type, code string) (
	types.Type,
	map[string]types.Type,
	error,
) {
	lexer := borsch.NewLexer(ctx.package_.Name, code)
	tokens, err := lexer.Lex()
	if err != nil {
		return nil, scope, err
	}

	p := parser.NewParser(ctx.package_.Name, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, scope, err
	}

	i.pushScope(ctx.package_.Name, builtin.RuntimeObjects)
	var result types.Type
	result, scope, _, err = i.executeAST(ctx, scope, asTree)
	if err != nil {
		return nil, scope, err
	}

	return result, scope, nil
}

func (i *Interpreter) ExecuteFile(ctx *Context, content []byte) error {
	_, scope, err := i.Execute(ctx, map[string]types.Type{}, string(content))
	if err != nil {
		return err
	}

	ctx.package_.Attributes = scope
	return nil
}
