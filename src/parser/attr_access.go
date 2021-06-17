package parser

import (
	"errors"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (p *Parser) parseAttrAccess(parent ast.ExpressionNode) (ast.ExpressionNode, error) {
	if name := p.match(models.TokenTypesList[models.Name]); name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			var err error
			parent, err = p.parseFunctionCall(name, parent)
			if err != nil {
				return nil, err
			}

			randomAccessOp, err := p.parseRandomAccessOperation(name, parent)
			if err != nil {
				return nil, err
			}

			if randomAccessOp != nil {
				parent = randomAccessOp
			}
		}

		var baseToken *models.Token = nil
		switch node := parent.(type) {
		case ast.AttrOpNode:
			baseToken = &node.Attr
		case ast.VariableNode:
			baseToken = &node.Variable
		}

		parent = ast.NewGetAttrOpNode(baseToken, parent, *name)
		randomAccessOp, err := p.parseRandomAccessOperation(baseToken, parent)
		if err != nil {
			return nil, err
		}

		if randomAccessOp != nil {
			parent = randomAccessOp
		}

		if dot := p.match(models.TokenTypesList[models.AttrAccessOp]); dot != nil {
			return p.parseAttrAccess(parent)
		}

		return parent, nil
	}

	return nil, errors.New("очікується назва атрибута")
}
