package main

import (
	"bufio"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	interpreter2 "github.com/YuriyLisovskiy/borsch/src/interpreter"
	"os"
)

func main() {
	stdRoot := os.Getenv("BORSCH_STD")
	interpreter := interpreter2.NewInterpreter(stdRoot)
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		err := interpreter.ExecuteFile(filePath)
		if err != nil {
			fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		scope := map[string]types.ValueType{}
		var err error
		for {
			fmt.Print(">>> ")
			code := ""
			for {
				fragment, err := reader.ReadString('\n')
				if err != nil {
					panic(err)
				}

				if fragment == "\n" {
					break
				} else {
					code += fragment
					fmt.Print("... ")
				}
			}

			scope, err = interpreter.Execute(scope, "<стдввід>", code)
			if err != nil {
				fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
			}
		}
	}
}
