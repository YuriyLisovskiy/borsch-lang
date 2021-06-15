package parser

import (
	"errors"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (p *Parser) parseGetAttr(parent ast.ExpressionNode) (ast.ExpressionNode, error) {
	if name := p.match(models.TokenTypesList[models.Name]); name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			p.pos--
			var err error
			parent, err = p.parseFunctionCall(parent)
			if err != nil {
				return nil, err
			}
		}

		parent = ast.NewGetAttrOpNode(parent, *name)
		if dot := p.match(models.TokenTypesList[models.AttrAccessOp]); dot != nil {
			return p.parseGetAttr(parent)
		}

		return parent, nil
	}

	return nil, errors.New("очікується назва атрибута")
}
