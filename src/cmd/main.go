package main

import (
	"bufio"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src"
	"os"
	"strings"
)

func main() {
	stdRoot := os.Getenv("BORSCH_STD")
	interpreter := src.NewInterpreter(stdRoot)
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		err := interpreter.ExecuteFile(filePath)
		if err != nil {
			fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(">>> ")
			code, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			code = strings.TrimSuffix(code, "\n")
			if code == "вихід()" {
				break
			}

			err = interpreter.Execute("<стдввід>", code)
			if err != nil {
				fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
			}
		}
	}
}
