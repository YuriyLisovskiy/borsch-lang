package common

import (
	"errors"
	"fmt"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
)

func assertionFailed(actual, expected string) string {
	return fmt.Sprintf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", actual, expected)
}

func TestStackTrace_Empty(t *testing.T) {
	st := StackTrace{}
	expected := ""
	actual := st.String(nil)
	if actual != expected {
		t.Error(assertionFailed(actual, expected))
	}
}

func TestStackTrace_ErrorOnly(t *testing.T) {
	st := StackTrace{}
	expected := `Помилка: сталося щось страшне`
	actual := st.String(errors.New("Помилка: сталося щось страшне"))
	if actual != expected {
		t.Error(assertionFailed(actual, expected))
	}
}

func TestStackTrace_SingleRow(t *testing.T) {
	st := StackTrace{}
	st.Push(
		NewTraceRow(
			lexer.Position{
				Filename: "/Users/проект/якийсь_пакет.борщ",
				Line:     3,
			},
			"якась_функція(1, 2)",
			"<пакет>",
		),
	)
	expected := `  Файл "/Users/проект/якийсь_пакет.борщ", рядок 3, у <пакет>
    якась_функція(1, 2)`
	actual := st.String(nil)
	if actual != expected {
		t.Error(assertionFailed(actual, expected))
	}
}

func TestStackTrace_SingleRowWithError(t *testing.T) {
	st := StackTrace{}
	st.Push(
		NewTraceRow(
			lexer.Position{
				Filename: "/Users/проект/якийсь_пакет.борщ",
				Line:     3,
			},
			"якась_функція(1, 2)",
			"<пакет>",
		),
	)
	expected := `  Файл "/Users/проект/якийсь_пакет.борщ", рядок 3, у <пакет>
    якась_функція(1, 2)
Помилка: сталося щось страшне`
	actual := st.String(errors.New("Помилка: сталося щось страшне"))
	if actual != expected {
		t.Error(assertionFailed(actual, expected))
	}
}

func TestStackTrace_MultipleRowsWithError(t *testing.T) {
	st := StackTrace{}
	st.Push(
		NewTraceRow(
			lexer.Position{
				Filename: "/Users/проект/якийсь_пакет.борщ",
				Line:     5,
			},
			"щось_зробити(\"Нічого не робити.\")",
			"щось_зробити",
		),
	)
	st.Push(
		NewTraceRow(
			lexer.Position{
				Filename: "/Users/проект/якийсь_пакет.борщ",
				Line:     3,
			},
			"якась_функція(1, 2)",
			"<пакет>",
		),
	)
	expected := `  Файл "/Users/проект/якийсь_пакет.борщ", рядок 5, у <пакет>
    щось_зробити("Нічого не робити.")
Помилка: сталося щось страшне`
	actual := st.String(errors.New("Помилка: сталося щось страшне"))
	if actual != expected {
		t.Error(assertionFailed(actual, expected))
	}
}
