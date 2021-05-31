package src

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/models"
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
//	notOp
)

var opTypeNames = []string{
	"додавання", "віднімання", "множення", "ділення", "логічного 'і'", "логічного 'або'",
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

func (e *Interpreter) executeCalculationOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType int, rootDir string,
) (builtin.ValueType, error) {
	left, err := e.executeNode(leftNode, rootDir)
	if err != nil {
		return builtin.NoneType{}, err
	}

	right, err := e.executeNode(rightNode, rootDir)
	if err != nil {
		return builtin.NoneType{}, err
	}

	if left.TypeHash() != right.TypeHash() {
		return builtin.NoneType{}, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opTypeNames[opType], left.TypeName(), right.TypeName(),
			),
		)
	}

	switch opType {
	case sumOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			return builtin.RealNumberType{
				Value: leftVal.Value + right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			return builtin.IntegerNumberType{
				Value: leftVal.Value + right.(builtin.IntegerNumberType).Value,
			}, nil
		case builtin.StringType:
			return builtin.StringType{
				Value: leftVal.Value + right.(builtin.StringType).Value,
			}, nil
		}

	case subOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			return builtin.RealNumberType{
				Value: leftVal.Value - right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			return builtin.IntegerNumberType{
				Value: leftVal.Value - right.(builtin.IntegerNumberType).Value,
			}, nil
		}
	case mulOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			return builtin.RealNumberType{
				Value: leftVal.Value * right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			return builtin.IntegerNumberType{
				Value: leftVal.Value * right.(builtin.IntegerNumberType).Value,
			}, nil
		}
	case divOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			rightVal := right.(builtin.RealNumberType).Value
			if rightVal == 0 {
				return builtin.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return builtin.RealNumberType{
				Value: leftVal.Value / right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			rightVal := right.(builtin.IntegerNumberType).Value
			if rightVal == 0 {
				return builtin.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return builtin.RealNumberType{
				Value: float64(leftVal.Value) / right.(builtin.RealNumberType).Value,
			}, nil
		}

	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}

	return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}

func (e *Interpreter) executeLogicalOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType int, rootDir string,
) (builtin.ValueType, error) {
	left, err := e.executeNode(leftNode, rootDir)
	if err != nil {
		return builtin.NoneType{}, err
	}

	right, err := e.executeNode(rightNode, rootDir)
	if err != nil {
		return builtin.NoneType{}, err
	}

	if left.TypeHash() != right.TypeHash() {
		return builtin.NoneType{}, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opTypeNames[opType], left.TypeName(), right.TypeName(),
			),
		)
	}

	switch opType {
	case andOp:
		switch leftVal := left.(type) {
		case builtin.BoolType:
			return builtin.BoolType{
				Value: leftVal.Value && right.(builtin.BoolType).Value,
			}, nil
		}
	case orOp:
		switch leftVal := left.(type) {
		case builtin.BoolType:
			return builtin.BoolType{
				Value: leftVal.Value || right.(builtin.BoolType).Value,
			}, nil
		}

	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}

	return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}

func (e *Interpreter) executeNode(rootNode ast.ExpressionNode, rootDir string) (builtin.ValueType, error) {
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
			sArg, err := e.executeNode(arg, rootDir)
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

	case ast.BinOperationNode:
		switch node.Operator.Type.Name {
		case models.Add:
			return e.executeCalculationOp(node.LeftNode, node.RightNode, sumOp, rootDir)

		case models.Sub:
			return e.executeCalculationOp(node.LeftNode, node.RightNode, subOp, rootDir)

		case models.Mul:
			return e.executeCalculationOp(node.LeftNode, node.RightNode, mulOp, rootDir)

		case models.Div:
			return e.executeCalculationOp(node.LeftNode, node.RightNode, divOp, rootDir)

		case models.And:
			return e.executeLogicalOp(node.LeftNode, node.RightNode, andOp, rootDir)

		case models.Or:
			return e.executeLogicalOp(node.LeftNode, node.RightNode, orOp, rootDir)

		case models.Assign:
			result, err := e.executeNode(node.RightNode, rootDir)
			if err != nil {
				return builtin.NoneType{}, err
			}

			variableNode := node.LeftNode.(ast.VariableNode)
			e.scope[variableNode.Variable.Text] = result
			return result, nil
		}

	case ast.RealNumberNode:
		return builtin.NewRealNumberType(node.Number.Text)

	case ast.IntegerNumberNode:
		return builtin.NewIntegerNumberType(node.Number.Text)

	case ast.StringNode:
		return builtin.StringType{Value: node.Content.Text}, nil

	case ast.BoolNode:
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
		_, err := e.executeNode(codeRow, dir)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n%s",
				filePath, codeRow.RowNumber(), codeRow.String(), err.Error(),
			))
		}
	}

	return nil
}

func (e *Interpreter) Execute(file string, code string) error {
	lexer := NewLexer(file, code)
	tokens, err := lexer.Lex()
	if err != nil {
		return err
	}

	parser := NewParser(file, tokens)
	asTree, err := parser.Parse()
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
