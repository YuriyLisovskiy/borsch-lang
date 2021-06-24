package interpreter

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"github.com/YuriyLisovskiy/borsch/src/parser"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
	"path/filepath"
)

const (
	// math
	exponentOp = iota
	moduloOp
	sumOp
	subOp
	mulOp
	divOp
	unaryMinus
	unaryPlus

	// logical
	andOp
	orOp
	notOp

	// conditional
	equalsOp
	notEqualsOp
	greaterOp
	greaterOrEqualsOp
	lessOp
	lessOrEqualsOp
)

type Operator int

var opTypeNames = map[Operator]string{
	exponentOp:        "піднесення до степеня",
	moduloOp:          "остачі від ділення",
	sumOp:             "додавання",
	subOp:             "віднімання",
	mulOp:             "множення",
	divOp:             "ділення",
	unaryMinus:        "унарного мінуса",
	unaryPlus:         "унарного плюса",
	andOp:             "логічного 'і'",
	orOp:              "логічного 'або'",
	notOp:             "логічного заперечення",
	equalsOp:          "рівності",
	notEqualsOp:       "нерівності",
	greaterOp:         "'більше'",
	greaterOrEqualsOp: "'більше або дорівнює'",
	lessOp:            "'менше'",
	lessOrEqualsOp:    "'менше або дорівнює'",
}

func (op Operator) Description() string {
	if op >= 0 && int(op) < len(opTypeNames) {
		return opTypeNames[op]
	}

	panic(fmt.Sprintf(
		"Unable to retrieve description for operator '%d', please add it to 'opTypeNames' map first",
		op,
	))
}

type Interpreter struct {
	stdRoot          string
	scopes           map[string][]map[string]types.ValueType
	currentPackage   string
	parentPackage    string
	includedPackages map[string]types.ValueType
}

func NewInterpreter(stdRoot, currentPackage, parentPackage string) *Interpreter {
	return &Interpreter{
		stdRoot:          stdRoot,
		currentPackage:   currentPackage,
		parentPackage:    parentPackage,
		scopes:           map[string][]map[string]types.ValueType{},
		includedPackages: map[string]types.ValueType{},
	}
}

func (i *Interpreter) pushScope(packageName string, scope map[string]types.ValueType) {
	i.scopes[packageName] = append(i.scopes[packageName], scope)
}

func (i *Interpreter) popScope(packageName string) map[string]types.ValueType {
	if scopes, ok := i.scopes[packageName]; ok {
		if len(scopes) == 0 {
			panic("Not enough scopes")
		}

		scope := scopes[len(scopes)-1]
		i.scopes[packageName] = scopes[:len(scopes)-1]
		return scope
	}

	panic(fmt.Sprintf("Scopes for '%s' package does not exist", packageName))
}

func (i *Interpreter) getVar(packageName string, name string) (types.ValueType, error) {
	if scopes, ok := i.scopes[packageName]; ok {
		lastScopeIdx := len(scopes) - 1
		for idx := lastScopeIdx; idx >= 0; idx-- {
			if val, ok := scopes[idx][name]; ok {
				return val, nil
			}
		}

		return nil, util.RuntimeError(fmt.Sprintf(
			"ідентифікатор '%s' не визначений", name,
		))
	}

	panic(fmt.Sprintf("Scopes for '%s' package does not exist", packageName))
}

func (i *Interpreter) setVar(packageName, name string, value types.ValueType) error {
	if scopes, ok := i.scopes[packageName]; ok {
		scopesLen := len(scopes)
		for idx := 0; idx < scopesLen; idx++ {
			if oldValue, ok := scopes[idx][name]; ok {
				if oldValue.TypeHash() != value.TypeHash() {
					if scopesLen == 1 {
						return util.RuntimeError(fmt.Sprintf(
							"неможливо записати значення типу '%s' у змінну '%s' з типом '%s'",
							value.TypeName(), name, oldValue.TypeName(),
						))
					}

					// TODO: надрукувати нормальне попередження!
					fmt.Println(fmt.Sprintf(
						"Увага: несумісні типи даних '%s' та '%s', змінна '%s' стає недоступною в поточному полі видимості",
						value.TypeName(), oldValue.TypeName(), name,
					))
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

	panic(fmt.Sprintf("Scopes for '%s' package does not exist", packageName))
}

// executeNode
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) executeNode(
	rootNode ast.ExpressionNode, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, bool, error) {
	switch node := rootNode.(type) {
	case ast.ImportNode:
		if node.IsStd {
			node.FilePath = filepath.Join(i.stdRoot, node.FilePath)
		} else if !filepath.IsAbs(node.FilePath) {
			node.FilePath = filepath.Join(rootDir, node.FilePath)
		}

		if node.FilePath == parentPackage {
			return nil, false, util.RuntimeError("циклічний імпорт заборонений")
		}

		pkg, ok := i.includedPackages[node.FilePath]
		if !ok {
			var err error
			fileContent, err := util.ReadFile(node.FilePath)
			if err != nil {
				return nil, false, err
			}

			pkg, err = i.ExecuteFile(node.FilePath, thisPackage, fileContent, node.IsStd)
			if err != nil {
				return nil, false, err
			}

			i.includedPackages[node.FilePath] = pkg
		}

		if node.Name != "" {
			err := i.setVar(thisPackage, node.Name, pkg)
			if err != nil {
				return nil, false, err
			}
		}

		return pkg, false, nil

	case ast.FunctionDefNode:
		functionDef := types.NewFunctionType(
			node.Name.Text, node.Arguments, node.ReturnType,
			func(_ []types.ValueType, kwargs map[string]types.ValueType) (types.ValueType, error) {
				res, _, err := i.executeBlock(kwargs, node.Body, thisPackage, parentPackage)
				return res, err
			},
		)
		return functionDef, false, i.setVar(thisPackage, node.Name.Text, functionDef)

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
		operand, _, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, false, err
		}

		switch node.Operator.Type.Name {
		case models.NotOp:
			switch operandVal := operand.(type) {
			case types.BoolType:
				return types.BoolType{Value: !operandVal.Value}, false, nil
			}

			return nil, false, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(notOp).Description(), operand.TypeName(),
			))
		case models.Add:
			switch operandVal := operand.(type) {
			case types.IntegerType, types.RealType:
				return operandVal, false, nil
			case types.BoolType:
				if operandVal.Value {
					return types.IntegerType{Value: 1}, false, nil
				}

				return types.IntegerType{Value: 0}, false, nil
			}

			return nil, false, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(unaryPlus).Description(), operand.TypeName(),
			))
		case models.Sub:
			switch operandVal := operand.(type) {
			case types.IntegerType:
				operandVal.Value = -operandVal.Value
				return operandVal, false, nil
			case types.RealType:
				operandVal.Value = -operandVal.Value
				return operandVal, false, nil
			case types.BoolType:
				if operandVal.Value {
					return types.IntegerType{Value: -1}, false, nil
				}

				return types.IntegerType{Value: 0}, false, nil
			}

			return nil, false, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(unaryMinus).Description(), operand.TypeName(),
			))
		}

	case ast.BinOperationNode:
		switch node.Operator.Type.Name {
		case models.ExponentOp:
			res, err := i.executeArithmeticOp(node.LeftNode, node.RightNode, exponentOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.ModuloOp:
			res, err := i.executeArithmeticOp(node.LeftNode, node.RightNode, moduloOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.Add:
			res, err := i.executeArithmeticOp(node.LeftNode, node.RightNode, sumOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.Sub:
			res, err := i.executeArithmeticOp(node.LeftNode, node.RightNode, subOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.Mul:
			res, err := i.executeArithmeticOp(node.LeftNode, node.RightNode, mulOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.Div:
			res, err := i.executeArithmeticOp(node.LeftNode, node.RightNode, divOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.AndOp:
			res, err := i.executeLogicalOp(node.LeftNode, node.RightNode, andOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.OrOp:
			res, err := i.executeLogicalOp(node.LeftNode, node.RightNode, orOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.EqualsOp:
			res, err := i.executeComparisonOp(node.LeftNode, node.RightNode, equalsOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.NotEqualsOp:
			res, err := i.executeComparisonOp(node.LeftNode, node.RightNode, notEqualsOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.GreaterOp:
			res, err := i.executeComparisonOp(node.LeftNode, node.RightNode, greaterOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.GreaterOrEqualsOp:
			res, err := i.executeComparisonOp(node.LeftNode, node.RightNode, greaterOrEqualsOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.LessOp:
			res, err := i.executeComparisonOp(node.LeftNode, node.RightNode, lessOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.LessOrEqualsOp:
			res, err := i.executeComparisonOp(node.LeftNode, node.RightNode, lessOrEqualsOp, rootDir, thisPackage, parentPackage)
			return res, false, err
		case models.Assign:
			result, _, err := i.executeNode(node.RightNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			switch assignmentNode := node.LeftNode.(type) {
			case ast.VariableNode:
				return nil, false, i.setVar(thisPackage, assignmentNode.Variable.Text, result)
			case ast.CallOpNode:
				return nil, false, util.RuntimeError("неможливо присвоїти значення виклику функції")
			case ast.RandomAccessOperationNode:
				variable, _, err := i.executeNode(assignmentNode.Operand, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, false, err
				}

				variable, err = i.executeRandomAccessSetOp(
					assignmentNode.Index, variable, result, rootDir, thisPackage, parentPackage,
				)
				if err != nil {
					return nil, false, err
				}

				operand := assignmentNode.Operand
				for {
					switch external := operand.(type) {
					case ast.RandomAccessOperationNode:
						opVar, _, err := i.executeNode(external.Operand, rootDir, thisPackage, parentPackage)
						if err != nil {
							return nil, false, err
						}

						variable, err = i.executeRandomAccessSetOp(
							external.Index, opVar, variable, rootDir, thisPackage, parentPackage,
						)
						if err != nil {
							return nil, false, err
						}

						operand = external.Operand
						continue
					case ast.VariableNode:
						err = i.setVar(thisPackage, external.Variable.Text, variable)
					}

					break
				}

				return variable, false, nil
			case ast.AttrOpNode:
				variable, _, err := i.executeNode(assignmentNode.Expression, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, false, err
				}

				switch attr := assignmentNode.Attr.(type) {
				case ast.VariableNode:
					variable, err = variable.SetAttr(attr.Variable.Text, result)
					if err != nil {
						return nil, false, err
					}
				case ast.CallOpNode:
					return nil, false, util.RuntimeError("неможливо присвоїти значення виклику функції")
				default:
					panic("fatal error")
				}

				if assignmentNode.Base != nil {
					return variable, false, i.setVar(thisPackage, assignmentNode.Base.Text, variable)
				}

				return variable, false, nil
			default:
				panic("fatal error")
			}
		}

	case ast.RandomAccessOperationNode:
		res, err := i.executeRandomAccessGetOp(node.Operand, node.Index, rootDir, thisPackage, parentPackage)
		return res, false, err

	case ast.ListSlicingNode:
		container, _, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, false, err
		}

		if container.TypeHash() == types.ListTypeHash {
			fromIdx, _, err := i.executeNode(node.LeftIndex, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			if fromIdx.TypeHash() == types.IntegerTypeHash {
				toIdx, _, err := i.executeNode(node.RightIndex, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, false, err
				}

				if toIdx.TypeHash() == types.IntegerTypeHash {
					res, err := container.(types.ListType).Slice(
						fromIdx.(types.IntegerType).Value, toIdx.(types.IntegerType).Value,
					)
					return res, false, err
				}

				return nil, false, util.RuntimeError("правий індекс має бути цілого типу")
			}

			return nil, false, util.RuntimeError("лівий індекс має бути цілого типу")
		}

		return nil, false, util.RuntimeError(fmt.Sprintf(
			"неможливо застосувати оператор відсікання списку до об'єкта з типом '%s'",
			container.TypeName(),
		))

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
		res, err := types.NewRealType(node.Value.Text)
		return res, false, err

	case ast.IntegerTypeNode:
		res, err := types.NewIntegerType(node.Value.Text)
		return res, false, err

	case ast.StringTypeNode:
		res := types.StringType{Value: node.Value.Text}
		return res, false, nil

	case ast.BoolTypeNode:
		res, err := types.NewBoolType(node.Value.Text)
		return res, false, err

	case ast.ListTypeNode:
		list := types.NewListType()
		for _, valueNode := range node.Values {
			value, _, err := i.executeNode(valueNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			list.Values = append(list.Values, value)
		}

		return list, false, nil

	case ast.DictionaryTypeNode:
		dict := types.NewDictionaryType()
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

	case ast.VariableNode:
		val, err := i.getVar(thisPackage, node.Variable.Text)
		if err != nil {
			return nil, false, err
		}

		return val, false, nil

	case ast.AttrOpNode:
		val, _, err := i.executeNode(node.Expression, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, false, err
		}

		if node.Attr != nil {
			switch attr := node.Attr.(type) {
			case ast.VariableNode:
				//val, _, err := i.executeNode(node.Expression, rootDir, thisPackage, parentPackage)
				//if err != nil {
				//	return nil, false, err
				//}

				val, err = val.GetAttr(attr.Variable.Text)
				if err != nil {
					return nil, false, err
				}

				return val, false, nil
			case ast.CallOpNode:
				//val, _, err := i.executeNode(attr.Parent, rootDir, thisPackage, parentPackage)
				//if err != nil {
				//	return nil, false, err
				//}

				//val, err = i.executeCallOp(&attr, rootDir, thisPackage, parentPackage)
				//return res, false, err

				val, err = val.GetAttr(attr.CallableName.Text)
				if err != nil {
					return nil, false, err
				}

				val, err = i.executeCallOp(&attr, val, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, false, err
				}

				return val, false, nil
			}
		} else {
			panic("unknown error")
		}

		//return val, false, nil
	}

	return nil, false, util.RuntimeError("невідома помилка")
}

// executeAST
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) executeAST(
	scope map[string]types.ValueType, thisPackage, parentPackage string, tree *ast.AST,
) (types.ValueType, map[string]types.ValueType, bool, error) {
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

	var result types.ValueType = nil
	forceReturn := false
	i.pushScope(thisPackage, scope)
	for _, node := range tree.Nodes {
		result, forceReturn, err = i.executeNode(node, dir, thisPackage, parentPackage)
		if err != nil {
			return nil, scope, false, errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n%s",
				filePath, node.RowNumber(), node.String(), err.Error(),
			))
		}

		if forceReturn {
			break
		}
	}

	scope = i.popScope(thisPackage)
	return result, scope, forceReturn, nil
}

// executeBlock
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) executeBlock(
	scope map[string]types.ValueType, tokens []models.Token, thisPackage, parentPackage string,
) (types.ValueType, bool, error) {
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

// Execute
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) Execute(
	thisPackage, parentPackage string, scope map[string]types.ValueType, code string,
) (types.ValueType, map[string]types.ValueType, error) {
	lexer := src.NewLexer(thisPackage, code)
	tokens, err := lexer.Lex()
	if err != nil {
		return nil, scope, err
	}

	p := parser.NewParser(thisPackage, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, scope, err
	}

	i.pushScope(thisPackage, builtin.GlobalScope)
	var result types.ValueType
	result, scope, _, err = i.executeAST(scope, thisPackage, parentPackage, asTree)
	if err != nil {
		return nil, scope, err
	}

	return result, scope, nil
}

// ExecuteFile
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) ExecuteFile(
	packageName, parentPackage string, content []byte, isBuiltin bool,
) (types.ValueType, error) {
	_, scope, err := i.Execute(packageName, parentPackage, map[string]types.ValueType{}, string(content))
	if err != nil {
		return nil, err
	}

	return types.NewPackageType(isBuiltin, packageName, parentPackage, scope), nil
}
