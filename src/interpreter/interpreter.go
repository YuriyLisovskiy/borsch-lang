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
			"змінну з назвою '%s' не знайдено", name,
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
) (types.ValueType, error) {
	switch node := rootNode.(type) {
	case ast.ImportNode:
		if node.IsStd {
			node.FilePath = filepath.Join(i.stdRoot, node.FilePath)
		} else if !filepath.IsAbs(node.FilePath) {
			node.FilePath = filepath.Join(rootDir, node.FilePath)
		}

		if node.FilePath == parentPackage {
			return nil, util.RuntimeError("циклічний імпорт заборонений")
		}

		pkg, ok := i.includedPackages[node.FilePath]
		if !ok {
			var err error
			fileContent, err := util.ReadFile(node.FilePath)
			if err != nil {
				return nil, err
			}

			pkg, err = i.ExecuteFile(node.FilePath, thisPackage, fileContent, node.IsStd)
			if err != nil {
				return nil, err
			}

			i.includedPackages[node.FilePath] = pkg
		}

		if node.Name != "" {
			err := i.setVar(thisPackage, node.Name, pkg)
			if err != nil {
				return nil, err
			}
		}

		return pkg, nil

	case ast.FunctionCallNode:
		var args []types.ValueType
		for _, arg := range node.Args {
			sArg, err := i.executeNode(arg, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			args = append(args, sArg)
		}

		function, found := builtin.FunctionsList[node.FunctionName.Text]
		if !found {
			return nil, util.RuntimeError(
				fmt.Sprintf("функцію з назвою '%s' не знайдено", node.FunctionName.Text),
			)
		}

		res, err := function(args...)
		if err != nil {
			return nil, err
		}

		return res, nil

	case ast.UnaryOperationNode:
		operand, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		switch node.Operator.Type.Name {
		case models.NotOp:
			switch operandVal := operand.(type) {
			case types.BoolType:
				return types.BoolType{Value: !operandVal.Value}, nil
			}

			return nil, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(notOp).Description(), operand.TypeName(),
			))
		case models.Add:
			switch operandVal := operand.(type) {
			case types.IntegerType, types.RealType:
				return operandVal, nil
			case types.BoolType:
				if operandVal.Value {
					return types.IntegerType{Value: 1}, nil
				}

				return types.IntegerType{Value: 0}, nil
			}

			return nil, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(unaryPlus).Description(), operand.TypeName(),
			))
		case models.Sub:
			switch operandVal := operand.(type) {
			case types.IntegerType:
				operandVal.Value = -operandVal.Value
				return operandVal, nil
			case types.RealType:
				operandVal.Value = -operandVal.Value
				return operandVal, nil
			case types.BoolType:
				if operandVal.Value {
					return types.IntegerType{Value: -1}, nil
				}

				return types.IntegerType{Value: 0}, nil
			}

			return nil, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(unaryMinus).Description(), operand.TypeName(),
			))
		}

	case ast.BinOperationNode:
		switch node.Operator.Type.Name {
		case models.ExponentOp:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, exponentOp, rootDir, thisPackage, parentPackage)
		case models.ModuloOp:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, moduloOp, rootDir, thisPackage, parentPackage)
		case models.Add:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, sumOp, rootDir, thisPackage, parentPackage)
		case models.Sub:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, subOp, rootDir, thisPackage, parentPackage)
		case models.Mul:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, mulOp, rootDir, thisPackage, parentPackage)
		case models.Div:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, divOp, rootDir, thisPackage, parentPackage)
		case models.AndOp:
			return i.executeLogicalOp(node.LeftNode, node.RightNode, andOp, rootDir, thisPackage, parentPackage)
		case models.OrOp:
			return i.executeLogicalOp(node.LeftNode, node.RightNode, orOp, rootDir, thisPackage, parentPackage)
		case models.EqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, equalsOp, rootDir, thisPackage, parentPackage)
		case models.NotEqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, notEqualsOp, rootDir, thisPackage, parentPackage)
		case models.GreaterOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, greaterOp, rootDir, thisPackage, parentPackage)
		case models.GreaterOrEqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, greaterOrEqualsOp, rootDir, thisPackage, parentPackage)
		case models.LessOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, lessOp, rootDir, thisPackage, parentPackage)
		case models.LessOrEqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, lessOrEqualsOp, rootDir, thisPackage, parentPackage)
		case models.Assign:
			result, err := i.executeNode(node.RightNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			switch assignmentNode := node.LeftNode.(type) {
			case ast.VariableNode:
				return nil, i.setVar(thisPackage, assignmentNode.Variable.Text, result)
			case ast.FunctionCallNode:
				return nil, util.RuntimeError("неможливо присвоїти значення виклику функції")
			case ast.RandomAccessOperationNode:
				variable, err := i.executeNode(assignmentNode.Operand, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, err
				}

				variable, err = i.executeRandomAccessSetOp(
					assignmentNode.Index, variable, result, rootDir, thisPackage, parentPackage,
				)
				if err != nil {
					return nil, err
				}

				operand := assignmentNode.Operand
				for {
					switch external := operand.(type) {
					case ast.RandomAccessOperationNode:
						opVar, err := i.executeNode(external.Operand, rootDir, thisPackage, parentPackage)
						if err != nil {
							return nil, err
						}

						variable, err = i.executeRandomAccessSetOp(
							external.Index, opVar, variable, rootDir, thisPackage, parentPackage,
						)
						if err != nil {
							return nil, err
						}

						operand = external.Operand
						continue
					case ast.VariableNode:
						err = i.setVar(thisPackage, external.Variable.Text, variable)
					}

					break
				}

				return variable, nil
			case ast.AttrOpNode:
				variable, err := i.executeNode(assignmentNode.Expression, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, err
				}

				variable, err = variable.SetAttr(assignmentNode.Attr.Text, result)
				if err != nil {
					return nil, err
				}

				if assignmentNode.Base != nil {
					return variable, i.setVar(thisPackage, assignmentNode.Base.Text, variable)
				}

				return variable, nil
			default:
				return nil, util.RuntimeError("неможливо присвоїти значення")
			}
		}

	case ast.RandomAccessOperationNode:
		return i.executeRandomAccessGetOp(node.Operand, node.Index, rootDir, thisPackage, parentPackage)

	case ast.ListSlicingNode:
		container, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		if container.TypeHash() == types.ListTypeHash {
			fromIdx, err := i.executeNode(node.LeftIndex, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			if fromIdx.TypeHash() == types.IntegerTypeHash {
				toIdx, err := i.executeNode(node.RightIndex, rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, err
				}

				if toIdx.TypeHash() == types.IntegerTypeHash {
					return container.(types.ListType).Slice(
						fromIdx.(types.IntegerType).Value, toIdx.(types.IntegerType).Value,
					)
				}

				return nil, util.RuntimeError("правий індекс має бути цілого типу")
			}

			return nil, util.RuntimeError("лівий індекс має бути цілого типу")
		}

		return nil, util.RuntimeError(fmt.Sprintf(
			"неможливо застосувати оператор відсікання списку до об'єкта з типом '%s'",
			container.TypeName(),
		))

	case ast.IfNode:
		return i.executeIfSequence(node.Blocks, node.ElseBlock, rootDir, thisPackage, parentPackage)

	case ast.ForEachNode:
		container, err := i.executeNode(node.Container, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		return i.executeForEachLoop(node.IndexVar, node.ItemVar, container, node.Body, thisPackage, parentPackage)

	case ast.RealTypeNode:
		return types.NewRealType(node.Value.Text)

	case ast.IntegerTypeNode:
		return types.NewIntegerType(node.Value.Text)

	case ast.StringTypeNode:
		return types.StringType{Value: node.Value.Text}, nil

	case ast.BoolTypeNode:
		return types.NewBoolType(node.Value.Text)

	case ast.ListTypeNode:
		list := types.NewListType()
		for _, valueNode := range node.Values {
			value, err := i.executeNode(valueNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			list.Values = append(list.Values, value)
		}

		return list, nil

	case ast.DictionaryTypeNode:
		dict := types.NewDictionaryType()
		for keyNode, valueNode := range node.Map {
			key, err := i.executeNode(keyNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			value, err := i.executeNode(valueNode, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			err = dict.SetElement(key, value)
			if err != nil {
				return nil, err
			}
		}

		return dict, nil

	case ast.VariableNode:
		val, err := i.getVar(thisPackage, node.Variable.Text)
		if err != nil {
			return nil, err
		}

		return val, nil
	case ast.AttrOpNode:
		parent, err := i.executeNode(node.Expression, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		val, err := parent.GetAttr(node.Attr.Text)
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	return nil, util.RuntimeError("невідома помилка")
}

// executeAST
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) executeAST(
	scope map[string]types.ValueType, thisPackage, parentPackage string, tree *ast.AST,
) (types.ValueType, map[string]types.ValueType, error) {
	var filePath string
	var dir string
	var err error
	if thisPackage == "<стдввід>" {
		filePath = "<стдввід>"
		dir, err = os.Getwd()
		if err != nil {
			return nil, scope, util.InternalError(err.Error())
		}
	} else {
		filePath, err = filepath.Abs(thisPackage)
		if err != nil {
			return nil, scope, util.InternalError(err.Error())
		}

		dir = filepath.Dir(filePath)
	}

	var result types.ValueType = nil
	i.pushScope(thisPackage, scope)
	for _, node := range tree.Nodes {
		result, err = i.executeNode(node, dir, thisPackage, parentPackage)
		if err != nil {
			return nil, scope, errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n%s",
				filePath, node.RowNumber(), node.String(), err.Error(),
			))
		}

		//if result != nil {
		//	return result, nil
		//}
	}

	scope = i.popScope(thisPackage)
	return result, scope, nil
}

// executeBlock
// 'thisPackage' is an absolute path to current package file
func (i *Interpreter) executeBlock(
	scope map[string]types.ValueType, tokens []models.Token, thisPackage, parentPackage string,
) (types.ValueType, error) {
	p := parser.NewParser(thisPackage, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, err
	}

	result, _, err := i.executeAST(scope, thisPackage, parentPackage, asTree)
	if err != nil {
		return nil, err
	}

	return result, nil
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

	var result types.ValueType
	result, scope, err = i.executeAST(scope, thisPackage, parentPackage, asTree)
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
