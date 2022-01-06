package parser

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

func (p *Parser) parseClassDefinition(doc string) (ast.ExpressionNode, error) {
	if p.match(models.TokenTypesList[models.ClassDef]) != nil {
		name, err := p.require(models.TokenTypesList[models.Name])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s класу", err.Error()))
		}

		// TODO: parse inherited classes here

		_, err = p.require(models.TokenTypesList[models.LCurlyBracket])
		if err != nil {
			return nil, err
		}

		// TODO: parse class scope
		var attributes []ast.ExpressionNode
		for {
			p.skipSemicolons()
			doc := p.tryParseDoc()

			functionNode, err := p.parseFunctionDefinition(name.Text, false, doc)
			if err != nil {
				return nil, err
			}

			if functionNode != nil {
				attributes = append(attributes, functionNode)
			}

			variableNode, err := p.parseVariableAssignment(name.Text, doc)
			if err != nil {
				return nil, err
			}

			if variableNode != nil {
				_, err = p.require(models.TokenTypesList[models.Semicolon])
				if err != nil {
					return nil, err
				}

				attributes = append(attributes, variableNode)
			}

			endScope := p.match(models.TokenTypesList[models.RCurlyBracket])
			if endScope != nil {
				break
			}
		}

		// TODO: check for end of file

		// _, err = p.require(models.TokenTypesList[models.RCurlyBracket])
		// if err != nil {
		// 	return nil, err
		// }

		classNode := ast.NewClassDefNode(*name, doc, attributes)
		return classNode, nil
	}

	return nil, nil
}
