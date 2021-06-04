package parser

import "github.com/YuriyLisovskiy/borsch/src/models"

func (p *Parser) readScope() ([]models.Token, error) {
	_, err := p.require(models.TokenTypesList[models.LCurlyBracket])
	if err != nil {
		return nil, err
	}

	fromPos := p.pos
	for p.pos < len(p.tokens) {
		if p.match(models.TokenTypesList[models.RCurlyBracket]) != nil {
			p.pos--
			break
		}

		p.pos++
	}

	_, err = p.require(models.TokenTypesList[models.RCurlyBracket])
	if err != nil {
		return nil, err
	}

	var result []models.Token
	if fromPos < p.pos {
		result = p.tokens[fromPos:p.pos-1]
	}

	return result, err
}
