package parser

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

func (p *Parser) parseArgument(futureClass string) (types.FunctionArgument, error) {
	paramName, err := p.require(models.TokenTypesList[models.Name])
	if err != nil {
		return types.FunctionArgument{}, errors.New(fmt.Sprintf("%s аргументу", err.Error()))
	}

	err = p.checkForKeyword(paramName.Text)
	if err != nil {
		return types.FunctionArgument{}, err
	}

	_, err = p.require(models.TokenTypesList[models.Colon])
	if err != nil {
		return types.FunctionArgument{}, err
	}

	isVariadic := p.match(models.TokenTypesList[models.TripleDot]) != nil
	typeName, err := p.require(models.TokenTypesList[models.Name])
	if err != nil {
		return types.FunctionArgument{}, errors.New(fmt.Sprintf("%s типу", err.Error()))
	}

	if !types.IsBuiltinType(typeName.Text) && futureClass != typeName.Text {
		return types.FunctionArgument{}, errors.New(fmt.Sprintf("невідомий тип '%s'", typeName.Text))
	}

	isNullable := p.match(models.TokenTypesList[models.QuestionMark]) != nil
	return types.FunctionArgument{
		TypeHash:   types.GetTypeHash(typeName.Text),
		Name:       paramName.Text,
		IsVariadic: isVariadic,
		IsNullable: isNullable,
	}, nil
}

func (p *Parser) parseFunctionDefinition(futureClass string) (ast.ExpressionNode, error) {
	if p.match(models.TokenTypesList[models.FunctionDef]) != nil {
		name, err := p.require(models.TokenTypesList[models.Name])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s функції", err.Error()))
		}

		_, err = p.require(models.TokenTypesList[models.LPar])
		if err != nil {
			return nil, err
		}

		var parameters []types.FunctionArgument
		if p.match(models.TokenTypesList[models.RPar]) == nil {
			for {
				argument, err := p.parseArgument(futureClass)
				if err != nil {
					return nil, err
				}

				parameters = append(parameters, argument)
				if p.match(models.TokenTypesList[models.Comma]) == nil {
					break
				}

				if argument.IsVariadic {
					return nil, errors.New("'...' можна використовувати лише для останнього аргумента")
				}
			}

			_, err = p.require(models.TokenTypesList[models.RPar])
			if err != nil {
				return nil, err
			}
		}

		visited := map[string]bool{}
		for _, parameter := range parameters {
			if visited[parameter.Name] {
				return nil, errors.New(fmt.Sprintf(
					"аргумент '%s' є продубльованим у визначенні функції", parameter.Name,
				))
			} else {
				visited[parameter.Name] = true
			}
		}

		retType := types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		}
		if p.match(models.TokenTypesList[models.Arrow]) != nil {
			retTypeName, err := p.require(models.TokenTypesList[models.Name])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s типу, який повертає функція", err.Error()))
			}

			if !types.IsBuiltinType(retTypeName.Text) && futureClass != retTypeName.Text {
				return nil, errors.New(fmt.Sprintf("невідомий тип '%s'", retTypeName.Text))
			}

			retType.TypeHash = types.GetTypeHash(retTypeName.Text)
			retType.IsNullable = p.match(models.TokenTypesList[models.QuestionMark]) != nil
		}

		body, err := p.readScope()
		if err != nil {
			return nil, err
		}

		functionNode := ast.NewFunctionDefNode(*name, parameters, retType, body)
		return functionNode, nil
	}

	return nil, nil
}

func (p *Parser) parseReturnStatement() (ast.ExpressionNode, error) {
	if retToken := p.match(models.TokenTypesList[models.Return]); retToken != nil {
		var value ast.ExpressionNode = nil
		if p.match(models.TokenTypesList[models.Semicolon]) != nil {
			p.pos--
		} else {
			var err error
			value, err = p.parseFormula()
			if err != nil {
				return nil, err
			}
		}

		return ast.NewReturnNode(value, retToken.Row), nil
	}

	return nil, nil
}
