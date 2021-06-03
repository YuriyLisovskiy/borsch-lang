package interpreter

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
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
	sumOp: "додавання",
	subOp: "віднімання",
	mulOp: "множення",
	divOp: "ділення",
	andOp: "логічного 'і'",
	orOp: "логічного 'або'",
	notOp: "логічного заперечення",
	equalsOp: "рівності",
	notEqualsOp: "нерівності",
	greaterOp: "'більше'",
	greaterOrEqualsOp: "'більше або дорівнює'",
	lessOp: "'менше'",
	lessOrEqualsOp: "'менше або дорівнює'",
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
	scope   map[string]builtin.ValueType
}

func NewInterpreter(stdRoot string) *Interpreter {
	return &Interpreter{
		stdRoot: stdRoot,
		scope: map[string]builtin.ValueType{},
	}
}

func (e *Interpreter) executeNode(
	rootNode ast.ExpressionNode, rootDir string, currentFile string,
) (builtin.ValueType, error) {
	switch node := rootNode.(type) {
	case ast.IncludeDirectiveNode:
		if node.IsStd {
			node.FilePath = filepath.Join(e.stdRoot, node.FilePath)
		} else if !filepath.IsAbs(node.FilePath) {
			node.FilePath = filepath.Join(rootDir, node.FilePath)
		}

		return builtin.NoneType{}, e.ExecuteFile(node.FilePath)

	case ast.FunctionCallNode:
		var args []builtin.ValueType
		for _, arg := range node.Args {
			sArg, err := e.executeNode(arg, rootDir, currentFile)
			if err != nil {
				return builtin.NoneType{}, err
			}

			args = append(args, sArg)
		}

		function, found := builtin.FunctionsList[node.FunctionName.Text]
		if !found {
			return builtin.NoneType{}, util.RuntimeError(
				fmt.Sprintf("функцію з назвою '%s' не знайдено", node.FunctionName.Text),
			)
		}

		res, err := function(args...)
		if err != nil {
			return builtin.NoneType{}, err
		}

		return res, nil

	case ast.UnaryOperationNode:
		operand, err := e.executeNode(node.Operand, rootDir, currentFile)
		if err != nil {
			return builtin.NoneType{}, err
		}

		switch node.Operator.Type.Name {
		case models.NotOp:
			switch operandVal := operand.(type) {
			case builtin.BoolType:
				return builtin.BoolType{Value: !operandVal.Value}, nil
			}

			return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(notOp).Description(), operand.TypeName(),
			))
		case models.Add:
			switch operandVal := operand.(type) {
			case builtin.IntegerNumberType, builtin.RealNumberType:
				return operandVal, nil
			case builtin.BoolType:
				return builtin.CastToInt(operandVal)
			}

			return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(notOp).Description(), operand.TypeName(),
			))
		case models.Sub:
			switch operandVal := operand.(type) {
			case builtin.IntegerNumberType:
				operandVal.Value = -operandVal.Value
				return operandVal, nil
			case builtin.RealNumberType:
				operandVal.Value = -operandVal.Value
				return operandVal, nil
			case builtin.BoolType:
				return builtin.CastToInt(operandVal)
			}

			return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
				"непідтримуваний тип операнда для оператора %s: '%s'",
				Operator(notOp).Description(), operand.TypeName(),
			))
		}

	case ast.BinOperationNode:
		switch node.Operator.Type.Name {
		case models.Add:
			return e.executeArithmeticOp(node.LeftNode, node.RightNode, sumOp, rootDir, currentFile)

		case models.Sub:
			return e.executeArithmeticOp(node.LeftNode, node.RightNode, subOp, rootDir, currentFile)

		case models.Mul:
			return e.executeArithmeticOp(node.LeftNode, node.RightNode, mulOp, rootDir, currentFile)

		case models.Div:
			return e.executeArithmeticOp(node.LeftNode, node.RightNode, divOp, rootDir, currentFile)

		case models.AndOp:
			return e.executeLogicalOp(node.LeftNode, node.RightNode, andOp, rootDir, currentFile)

		case models.OrOp:
			return e.executeLogicalOp(node.LeftNode, node.RightNode, orOp, rootDir, currentFile)

		case models.EqualsOp:
			return e.executeComparisonOp(node.LeftNode, node.RightNode, equalsOp, rootDir, currentFile)

		case models.NotEqualsOp:
			return e.executeComparisonOp(node.LeftNode, node.RightNode, notEqualsOp, rootDir, currentFile)

		case models.GreaterOp:
			return e.executeComparisonOp(node.LeftNode, node.RightNode, greaterOp, rootDir, currentFile)

		case models.GreaterOrEqualsOp:
			return e.executeComparisonOp(node.LeftNode, node.RightNode, greaterOrEqualsOp, rootDir, currentFile)

		case models.LessOp:
			return e.executeComparisonOp(node.LeftNode, node.RightNode, lessOp, rootDir, currentFile)

		case models.LessOrEqualsOp:
			return e.executeComparisonOp(node.LeftNode, node.RightNode, lessOrEqualsOp, rootDir, currentFile)

		case models.Assign:
			result, err := e.executeNode(node.RightNode, rootDir, currentFile)
			if err != nil {
				return builtin.NoneType{}, err
			}

			variableNode := node.LeftNode.(ast.VariableNode)
			e.scope[variableNode.Variable.Text] = result
			return result, nil
		}

	case ast.IfSequenceNode:
		return e.executeIfSequence(node.Blocks, node.ElseBlock, rootDir, currentFile)
		
	case ast.RealTypeNode:
		return builtin.NewRealNumberType(node.Value.Text)

	case ast.IntegerTypeNode:
		return builtin.NewIntegerNumberType(node.Value.Text)

	case ast.StringTypeNode:
		return builtin.StringType{Value: node.Value.Text}, nil

	case ast.BoolTypeNode:
		return builtin.NewBoolType(node.Value.Text)

	case ast.VariableNode:
		if val, ok := e.scope[node.Variable.Text]; ok {
			return val, nil
		}

		return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"змінну з назвою '%s' не знайдено", node.Variable.Text,
		))
	}

	return builtin.NoneType{}, util.RuntimeError("невідома помилка")
}

func (e *Interpreter) executeAST(file string, tree *ast.AST) error {
	var filePath string
	var dir string
	var err error
	if file == "<стдввід>" {
		filePath = "<стдввід>"
		dir, err = os.Getwd()
		if err != nil {
			return util.InternalError(err.Error())
		}
	} else {
		filePath, err = filepath.Abs(file)
		if err != nil {
			return util.InternalError(err.Error())
		}

		dir = filepath.Dir(filePath)
	}

	for _, codeRow := range tree.CodeRows {
		_, err := e.executeNode(codeRow, dir, file)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n%s",
				filePath, codeRow.RowNumber(), codeRow.String(), err.Error(),
			))
		}
	}

	return nil
}

func (e *Interpreter) executeBlock(tokens []models.Token, currentFile string) (builtin.ValueType, error) {
	p := parser.NewParser(currentFile, tokens)
	asTree, err := p.Parse()
	if err != nil {
		return nil, err
	}

	err = e.executeAST(currentFile, asTree)
	if err != nil {
		return nil, err
	}

	return builtin.NoneType{}, nil
}

func (e *Interpreter) Execute(file string, code string) error {
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

	err = e.executeAST(file, asTree)
	if err != nil {
		return err
	}

	return nil
}

func (e *Interpreter) ExecuteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("файл з ім'ям '%s' не існує", filePath))
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return e.Execute(filePath, string(content))
}
