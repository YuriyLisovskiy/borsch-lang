package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/interpreter"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
	"github.com/peterh/liner"
)

var (
	historyFile = filepath.Join(os.TempDir(), ".borsch_interactive_console_history")
	keywords    []string
	// packages    map[string][]string
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

func pushKeywords(parent, name string, value types.Type) {
	switch v := value.(type) {
	case types.SequentialType, types.DictionaryInstance, types.BoolInstance, types.IntegerInstance, types.NilInstance, types.RealInstance:
		if parent != "" {
			// packages[parent] = append(packages[parent], name)
		} else {
			keywords = append(keywords, name)
		}
	case types.PackageInstance:
		keywords = append(keywords, name)
		for attrName, attrValue := range v.Object.Attributes {
			pushKeywords(name, attrName, attrValue)
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

func makeCompleter(nameEndingRegex regexp.Regexp) func(string, int) (string, []string, string) {
	return func(line string, pos int) (head string, completions []string, tail string) {
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
	}
}

func runInteractiveConsole(interpreterInstance *interpreter.Interpreter) {
	editor := liner.NewLiner()
	defer func() {
		if err := editor.Close(); err != nil {
			panic(err)
		}
	}()

	editor.SetCtrlCAborts(true)
	editor.SetWordCompleter(makeCompleter(*regexp.MustCompile("(" + models.RawNameRegex + "$)")))

	if file, err := os.Open(historyFile); err == nil {
		_, err = editor.ReadHistory(file)
		if err != nil {
			panic(util.InternalError(err.Error()))
		}

		if err = file.Close(); err != nil {
			panic(util.InternalError(err.Error()))
		}
	}

	scope := map[string]types.Type{}
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
			if fragment == ";" ||
				(!(strings.Contains(code, "{") ||
					strings.Contains(code, "}")) &&
					strings.HasSuffix(fragment, ";")) {
				break
			}

			iteration++
		}

		if quit {
			break
		}

		var result types.Type
		var err error
		result, scope, err = interpreterInstance.Execute(
			interpreterInstance.GetContext(),
			scope,
			strings.TrimPrefix(code, "\n"),
		)
		if err != nil {
			fmt.Printf("Відстеження (стек викликів):\n%s\n", err.Error())
		} else if result != nil {
			switch result.(type) {
			case types.NilInstance, types.PackageInstance:
			default:
				fmt.Println(result.Representation())
			}
		}

		for name, value := range scope {
			pushKeywords("", name, value)
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
