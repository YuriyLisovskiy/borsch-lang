package models

type Token struct {
	Type TokenType
	Text string
	Pos  int
	Row  int
}

func (t *Token) String() string {
	return t.Text
}
