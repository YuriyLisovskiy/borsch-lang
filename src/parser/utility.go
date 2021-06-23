package parser

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

//func (p *Parser) readFunctionScope(retType int) ([]models.Token, error) {
//	_, err := p.require(models.TokenTypesList[models.LCurlyBracket])
//	if err != nil {
//		return nil, err
//	}
//
//	fromPos := p.pos
//	unreachableCount := 0
//	gotRet := false
//	for p.pos < len(p.tokens) {
//		if p.match(models.TokenTypesList[models.RCurlyBracket]) != nil {
//			p.pos--
//			break
//		}
//
//		if tok := p.current(); tok != nil && tok.Type.Name == models.Return {
//			gotRet = true
//		}
//
//		if gotRet {
//			unreachableCount++
//		} else {
//			p.pos++
//		}
//	}
//
//	_, err = p.require(models.TokenTypesList[models.RCurlyBracket])
//	if err != nil {
//		return nil, err
//	}
//
//	var result []models.Token
//	if fromPos < p.pos {
//		result = p.tokens[fromPos : p.pos-unreachableCount-1]
//	}
//
//	if retType != types.NoneTypeHash {
//
//	}
//
//	return result, err
//}

func (p *Parser) readScope() ([]models.Token, error) {
	_, err := p.require(models.TokenTypesList[models.LCurlyBracket])
	if err != nil {
		return nil, err
	}

	lCurlyBracketsCount := 0
	fromPos := p.pos
	for p.pos < len(p.tokens) {
		tok := p.match(models.TokenTypesList[models.LCurlyBracket], models.TokenTypesList[models.RCurlyBracket])
		if tok != nil {
			if tok.Type.Name == models.LCurlyBracket {
				lCurlyBracketsCount++
			} else if tok.Type.Name == models.RCurlyBracket {
				if lCurlyBracketsCount == 0 {
					p.pos--
					break
				}

				lCurlyBracketsCount--
			}
		}

		p.pos++
	}

	_, err = p.require(models.TokenTypesList[models.RCurlyBracket])
	if err != nil {
		return nil, err
	}

	var result []models.Token
	if fromPos < p.pos {
		result = p.tokens[fromPos : p.pos-1]
	}

	return result, err
}
