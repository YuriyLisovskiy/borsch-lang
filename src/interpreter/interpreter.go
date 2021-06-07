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
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// math
	sumOp = iota
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
	stdRoot string
	scopes  []map[string]types.ValueType
}

func NewInterpreter(stdRoot string) *Interpreter {
	return &Interpreter{
		stdRoot: stdRoot,
		scopes:  []map[string]types.ValueType{},
	}
}

func (i *Interpreter) pushScope(scope map[string]types.ValueType) {
	i.scopes = append(i.scopes, scope)
}

func (i *Interpreter) popScope() {
	if len(i.scopes) == 0 {
		panic("Not enough scopes")
	}

	i.scopes = i.scopes[:len(i.scopes)-1]
}

func (i *Interpreter) getVar(name string) (types.ValueType, error) {
	lastScopeIdx := len(i.scopes) - 1
	for idx := lastScopeIdx; idx >= 0; idx-- {
		if val, ok := i.scopes[idx][name]; ok {
			return val, nil
		}
	}

	return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"змінну з назвою '%s' не знайдено", name,
	))
}

func (i *Interpreter) setVar(name string, value types.ValueType) error {
	scopesLen := len(i.scopes)
	for idx := 0; idx < scopesLen; idx++ {
		if oldValue, ok := i.scopes[idx][name]; ok {
			if oldValue.TypeHash() != value.TypeHash() {
				if scopesLen == 1 {
					return util.RuntimeError(fmt.Sprintf(
						"неможливо записати значення типу '%s' у змінну '%s' з типом '%s'",
						value.TypeName(), name, oldValue.TypeName(),
					))
				}

				// TODO: надрукувати нормальне попередження!
				fmt.Println(fmt.Sprintf(
					"Попередження: несумісні типи даних '%s' та '%s', змінна '%s' стає недоступною в поточному полі видимості",
					value.TypeName(), oldValue.TypeName(), name,
				))
				break
			}

			i.scopes[idx][name] = value
			return nil
		}
	}

	i.scopes[scopesLen-1][name] = value
	return nil
}

func (i *Interpreter) executeNode(
	rootNode ast.ExpressionNode, rootDir string, currentFile string,
) (types.ValueType, error) {
	switch node := rootNode.(type) {
	case ast.IncludeDirectiveNode:
		if node.IsStd {
			node.FilePath = filepath.Join(i.stdRoot, node.FilePath)
		} else if !filepath.IsAbs(node.FilePath) {
			node.FilePath = filepath.Join(rootDir, node.FilePath)
		}

		return nil, i.ExecuteFile(node.FilePath)

	case ast.FunctionCallNode:
		var args []types.ValueType
		for _, arg := range node.Args {
			sArg, err := i.executeNode(arg, rootDir, currentFile)
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
		operand, err := i.executeNode(node.Operand, rootDir, currentFile)
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
		case models.Add:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, sumOp, rootDir, currentFile)

		case models.Sub:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, subOp, rootDir, currentFile)

		case models.Mul:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, mulOp, rootDir, currentFile)

		case models.Div:
			return i.executeArithmeticOp(node.LeftNode, node.RightNode, divOp, rootDir, currentFile)

		case models.AndOp:
			return i.executeLogicalOp(node.LeftNode, node.RightNode, andOp, rootDir, currentFile)

		case models.OrOp:
			return i.executeLogicalOp(node.LeftNode, node.RightNode, orOp, rootDir, currentFile)

		case models.EqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, equalsOp, rootDir, currentFile)

		case models.NotEqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, notEqualsOp, rootDir, currentFile)

		case models.GreaterOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, greaterOp, rootDir, currentFile)

		case models.GreaterOrEqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, greaterOrEqualsOp, rootDir, currentFile)

		case models.LessOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, lessOp, rootDir, currentFile)

		case models.LessOrEqualsOp:
			return i.executeComparisonOp(node.LeftNode, node.RightNode, lessOrEqualsOp, rootDir, currentFile)

		case models.Assign:
			result, err := i.executeNode(node.RightNode, rootDir, currentFile)
			if err != nil {
				return nil, err
			}

			switch assignmentNode := node.LeftNode.(type) {
			case ast.VariableNode:
				return nil, i.setVar(assignmentNode.Variable.Text, result)
			case ast.RandomAccessSetOperationNode:
				variable, err := i.getVar(assignmentNode.Variable.Text)
				if err != nil {
					return nil, err
				}

				variable, err = i.executeRandomAccessSetOp(assignmentNode.Index, variable, result, rootDir, currentFile)
				if err != nil {
					return nil, err
				}

				return nil, i.setVar(assignmentNode.Variable.Text, variable)
			default:
				// TODO: обробити помилку
			}
		}

	case ast.RandomAccessGetOperationNode:
		return i.executeRandomAccessGetOp(node.Operand, node.Index, rootDir, currentFile)

	case ast.IfSequenceNode:
		return i.executeIfSequence(node.Blocks, node.ElseBlock, rootDir, currentFile)

	case ast.ForEachNode:
		container, err := i.executeNode(node.Container, rootDir, currentFile)
		if err != nil {
			return nil, err
		}

		return i.executeForEachLoop(node.IndexVar, node.ItemVar, container, node.Body, currentFile)

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
			value, err := i.executeNode(valueNode, rootDir, currentFile)
			if err != nil {
				return nil, err
			}

			list.Values = append(list.Values, value)
		}

		return list, nil

	case ast.VariableNode:
		val, err := i.getVar(node.Variable.Text)
		if err != nil {
			return types.NoneType{}, err
		}

		return val, nil
	}

	return nil, util.RuntimeError("невідома помилка")
}

func (i *Interpreter) executeAST(
	scope map[string]types.ValueType, file string, tree *ast.AST,
) (types.ValueType, error) {
	var filePath string
	var dir string
	var err error
	if file == "<стдввід>" {
		filePath = "<стдввід>"
		dir, err = os.Getwd()
		if err != nil {
			return nil, util.InternalError(err.Error())
		}
	} else {
		filePath, err = filepath.Abs(file)
		if err != nil {
			return nil, util.InternalError(err.Error())
		}

		dir = filepath.Dir(filePath)
	}

	i.pushScope(scope)
	for _, node := range tree.CodeNodes {
		_, err := i.executeNode(node, dir, file)
		if err != nil {
			return nil, errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n%s",
				filePath, node.RowNumber(), node.String(), err.Error(),
			))
		}

		//if result != nil {
		//	return result, nil
		//}
	}

	i.popScope()
	return nil, nil
}

func (i *Interpreter) executeBlock(
	scope map[string]types.ValueType, tokens []models.Token, currentFile string,
) (types.ValueType, error) {
	p := parser.NewParser(currentFile, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, err
	}

	result, err := i.executeAST(scope, currentFile, asTree)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *Interpreter) Execute(file string, code string) error {
	lexer := src.NewLexer(file, code)
	tokens, err := lexer.Lex()
	if err != nil {
		return err
	}

	p := parser.NewParser(file, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return err
	}

	_, err = i.executeAST(map[string]types.ValueType{}, file, asTree)
	if err != nil {
		return err
	}

	return nil
}

func (i *Interpreter) ExecuteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return util.RuntimeError(fmt.Sprintf("файл з ім'ям '%s' не існує", filePath))
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return i.Execute(filePath, string(content))
}
