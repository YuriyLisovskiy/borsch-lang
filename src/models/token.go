package models

import "fmt"

type Token struct {
	Type TokenType
	Text string
	Pos  int
}

func (t *Token) String() string {
	return fmt.Sprintf("[%s | %s | %d]", t.Type.String(), t.Text, t.Pos)
}
