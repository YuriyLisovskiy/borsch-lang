package lang

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/models"
	"sort"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	filePath         string
	code             string
	pos              int
	tokenList        []models.Token
	tokenTypesValues []models.TokenType
}

func NewLexer(filePath string, code string) *Lexer {
	values := make([]models.TokenType, 0, len(models.TokenTypesList))
	for _, value := range models.TokenTypesList {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Name < values[j].Name
	})

	return &Lexer{
		filePath:         filePath,
		code:             code,
		pos:              0,
		tokenList:        []models.Token{},
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
	rowCounter := 1
	for _, token := range l.tokenList {
		switch token.Type.Name {
		case models.Space:
			if token.Text == "\n" {
				rowCounter++
			}
		case models.MultiLineComment, models.SingleLineComment:
			rowCounter += strings.Count(token.Text, "\n")
		default:
			token.Row = rowCounter
			result = append(result, token)
		}
	}

	l.tokenList = result
	return l.tokenList, nil
}

func (l *Lexer) nextToken() (bool, error) {
	codeSize := utf8.RuneCountInString(l.code)
	if l.pos >= codeSize {
		return false, nil
	}

	runes := []rune(l.code)
	for _, tokenType := range l.tokenTypesValues {
		strToMatch := string(runes[l.pos:])
		result := tokenType.Regex.FindString(strToMatch)
		if utf8.RuneCountInString(result) > 0 {
			token := models.Token{
				Type:            tokenType,
				Text:            result,
				Pos:             l.pos,
				IsUnaryOperator: false,
			}
			l.pos += utf8.RuneCountInString(result)
			l.tokenList = append(l.tokenList, token)
			return true, nil
		}
	}

	rowNumber := strings.Count(string(runes[:l.pos]), "\n") + 1
	leftPos := l.pos
	for leftPos > 0 {
		if runes[leftPos] == '\n' {
			break
		}

		leftPos--
	}

	rightPos := l.pos
	for rightPos < codeSize {
		if runes[rightPos] == '\n' {
			break
		}

		rightPos++
	}

	if leftPos != 0 {
		leftPos++
	}

	codeFragment := string(runes[leftPos : rightPos-1])
	underline := strings.Repeat(" ", len(runes[leftPos:l.pos])) + "^"
	return false, errors.New(fmt.Sprintf(
		"  Файл \"%s\", рядок %d\n    %s\n    %s\n%s",
		l.filePath, rowNumber, codeFragment, underline, "Синтаксичка помилка: некоректний синтаксис",
	))
}
