package cmd

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/interpreter"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"github.com/peterh/liner"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	historyFile = filepath.Join(os.TempDir(), ".borsch_interactive_console_history")
	keywords    []string
)

func init() {
	for name, id := range builtin.RegisteredIdentifiers {
		if id != builtin.ConstantKeywordId && name != "інакше" {
			keywords = append(keywords, name+"(")
		} else {
			keywords = append(keywords, name)
		}
	}
}

func inputToHistory(editor *liner.State, prompt string) (fragment string, quit bool) {
	var err error
	if fragment, err = editor.Prompt(prompt); err == nil {
		if fragment != "" {
			editor.AppendHistory(fragment)
		}

		return fragment, false
	} else if err == liner.ErrPromptAborted {
		return "", true
	} else {
		return "", true
	}
}

func getPromptText(iteration int) string {
	if iteration > 0 {
		return "... "
	}

	return ">>> "
}

func runInteractiveConsole(interpreterInstance *interpreter.Interpreter) {
	editor := liner.NewLiner()
	defer func() {
		if err := editor.Close(); err != nil {
			panic(err)
		}
	}()

	editor.SetCtrlCAborts(true)

	nameEndingRegex := regexp.MustCompile("(" + models.RawNameRegex + "$)")
	editor.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		head = string([]rune(line)[:pos])
		tail = string([]rune(line)[pos:])
		matches := nameEndingRegex.FindAllString(head, -1)
		if len(matches) > 0 {
			lastMatch := matches[len(matches)-1]
			head = strings.TrimSuffix(head, lastMatch)
			for _, keyword := range keywords {
				if strings.HasPrefix(keyword, strings.ToLower(lastMatch)) {
					completions = append(completions, keyword)
				}
			}
		}

		return
	})

	if file, err := os.Open(historyFile); err == nil {
		_, err = editor.ReadHistory(file)
		if err != nil {
			panic(util.InternalError(err.Error()))
		}

		if err = file.Close(); err != nil {
			panic(util.InternalError(err.Error()))
		}
	}

	scope := map[string]types.ValueType{}
	var quit bool
	for {
		code := ""
		iteration := 0
		for {
			var fragment string
			fragment, quit = inputToHistory(editor, getPromptText(iteration))
			if quit || fragment == "" {
				break
			}

			code += "\n" + fragment
			if fragment == ";" {
				break
			}

			iteration++
		}

		if quit {
			break
		}

		var result types.ValueType
		var err error
		result, scope, err = interpreterInstance.Execute(
			scope, "<стдввід>", strings.TrimPrefix(code, "\n"),
		)
		if err != nil {
			fmt.Printf("Відстеження (стек викликів):\n%s\n", err.Error())
		} else if result != nil {
			if _, ok := result.(types.NoneType); !ok {
				fmt.Println(result.Representation())
			}
		}
	}

	if file, err := os.Create(historyFile); err == nil {
		_, err = editor.WriteHistory(file)
		if err != nil {
			panic(util.InternalError(err.Error()))
		}

		if err = file.Close(); err != nil {
			panic(util.InternalError(err.Error()))
		}
	}
}
