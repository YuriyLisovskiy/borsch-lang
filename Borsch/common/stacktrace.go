package common

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type TraceRow struct {
	pos       lexer.Position
	statement string
	place     string
}

func NewTraceRow(pos lexer.Position, statement, place string) *TraceRow {
	return &TraceRow{
		pos:       pos,
		statement: statement,
		place:     place,
	}
}

func (e *TraceRow) String(place string) string {
	return fmt.Sprintf(
		"  Файл \"%s\", рядок %d, у %s\n    %s",
		e.pos.Filename,
		e.pos.Line,
		place,
		e.statement,
	)
}

type StackTrace []*TraceRow

func (st *StackTrace) Push(row *TraceRow) {
	if row == nil {
		panic("stack trace row is nil")
	}

	*st = append(*st, row)
}

func (st *StackTrace) Pop() {
	stLen := len(*st)
	if stLen == 0 {
		panic("stack trace is empty")
	}

	*st = (*st)[:stLen-1]
}

func (st StackTrace) String(err error) string {
	traceLen := len(st)
	var rows []string
	if traceLen == 1 {
		rows = append(rows, st[0].String(st[0].place))
	} else if traceLen > 1 {
		for i := traceLen - 2; i >= 0; i-- {
			// if st[i].pos.Line != st[i+1].pos.Line {
			rows = append(rows, st[i].String(st[i+1].place))
			// }
		}
	}

	if err != nil {
		rows = append(rows, err.Error())
	}

	return strings.Join(rows, "\n")
}
