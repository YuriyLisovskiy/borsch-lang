package src

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"io/ioutil"
	"os"
	"strconv"
)

type Interpreter struct {
	scope map[string]string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		scope: map[string]string{},
	}
}

func (e *Interpreter) runNumbers(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode,
) (float64, float64, error) {
	leftStr, err := e.executeNode(leftNode)
	if err != nil {
		return 0, 0, err
	}

	leftNumber, err := strconv.ParseFloat(leftStr, 64)
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("Помилка виконання: %s", err.Error()))
	}

	rightStr, err := e.executeNode(rightNode)
	if err != nil {
		return 0, 0, err
	}

	rightNumber, err := strconv.ParseFloat(rightStr, 64)
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("Помилка виконання: %s", err.Error()))
	}

	return leftNumber, rightNumber, nil
}

func (e *Interpreter) executeNode(rootNode ast.ExpressionNode) (string, error) {
	switch node := rootNode.(type) {
	case ast.IncludeDirectiveNode:
		return "", e.ExecuteFile(node.FilePath)

	case ast.FunctionCallNode:
		var args []string
		for _, arg := range node.Args {
			sArg, err := e.executeNode(arg)
			if err != nil {
				return "", err
			}

			args = append(args, sArg)
		}

		function, found := builtin.FunctionsList[node.FunctionName.Text]
		if !found {
			return "", errors.New(
				fmt.Sprintf("Помилка виконання: функцію з назвою '%s' не знайдено", node.FunctionName.Text),
			)
		}

		res, err := function(args...)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Помилка виконання: %s", err.Error()))
		}

		return res, nil

	case ast.BinOperationNode:
		switch node.Operator.Type.Name {
		case models.Add:
			left, right, err := e.runNumbers(node.LeftNode, node.RightNode)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%f", left+right), nil

		case models.Sub:
			left, right, err := e.runNumbers(node.LeftNode, node.RightNode)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%f", left-right), nil

		case models.Mul:
			left, right, err := e.runNumbers(node.LeftNode, node.RightNode)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%f", left*right), nil

		case models.Div:
			left, right, err := e.runNumbers(node.LeftNode, node.RightNode)
			if err != nil {
				return "", err
			}

			if right == 0 {
				return "", errors.New(fmt.Sprintf("Помилка виконання: ділення на нуль"))
			}

			return fmt.Sprintf("%f", left/right), nil

		case models.Assign:
			result, err := e.executeNode(node.RightNode)
			if err != nil {
				return "", err
			}

			variableNode := node.LeftNode.(ast.VariableNode)
			e.scope[variableNode.Variable.Text] = result
			return result, nil
		}

	case ast.NumberNode:
		return node.Number.Text, nil

	case ast.StringNode:
		return node.Content.Text, nil

	case ast.VariableNode:
		if val, ok := e.scope[node.Variable.Text]; ok {
			return val, nil
		}

		return "", errors.New(fmt.Sprintf(
			"Помилка виконання: змінну з назвою '%s' не знайдено", node.Variable.Text,
		))
	}

	return "", errors.New(fmt.Sprintf("Помилка виконання: невідома помилка"))
}

func (e *Interpreter) executeAST(file string, tree *ast.AST) error {
	for _, codeRow := range tree.CodeRows {
		_, err := e.executeNode(codeRow)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n%s",
				file, codeRow.RowNumber(), codeRow.String(), err.Error(),
			))
		}
	}

	return nil
}

func (e *Interpreter) Execute(file string, code string) error {
	lexer := NewLexer(code)
	tokens, err := lexer.Lex()
	if err != nil {
		return err
	}

	parser, err := NewParser(file, tokens)
	if err != nil {
		return err
	}

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
