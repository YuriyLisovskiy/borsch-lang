package src

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"unicode/utf8"
)

type Lexer struct {
	code             string
	pos              int
	tokenList        []models.Token
	tokenTypesValues []models.TokenType
}

func NewLexer(code string) *Lexer {
	values := make([]models.TokenType, 0, len(models.TokenTypesList))
	for _, value := range models.TokenTypesList {
		values = append(values, value)
	}

	return &Lexer{
		code:      code,
		pos:       0,
		tokenList: []models.Token{},
		tokenTypesValues: values,
	}
}

func (l *Lexer) Lex() ([]models.Token, error) {
	for {
		hasNext, err := l.nextToken()
		if err != nil {
			return nil, err
		}

		if !hasNext {
			break
		}
	}

	var result []models.Token
	for _, token := range l.tokenList {
		switch models.TokenTypesList[token.Type.Name].Name {
		case models.Space, models.SingleLineComment:
			break
		default:
			result = append(result, token)
		}
	}

	l.tokenList = result
	return l.tokenList, nil
}

func (l *Lexer) nextToken() (bool, error) {
	if l.pos >= utf8.RuneCountInString(l.code) {
		return false, nil
	}

	for _, tokenType := range l.tokenTypesValues {
		runes := []rune(l.code)
		strToMatch := string(runes[l.pos:])
		result := tokenType.Regex.FindString(strToMatch)
		if utf8.RuneCountInString(result) > 0 {
			token := models.Token{
				Type:  tokenType,
				Text: result,
				Pos:  l.pos,
			}
			l.pos += utf8.RuneCountInString(result)
			l.tokenList = append(l.tokenList, token)
			return true, nil
		}
	}

	return false, errors.New(fmt.Sprintf("На позиції %d знайдено помилку", l.pos))
}
