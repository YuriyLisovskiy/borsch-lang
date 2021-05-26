package main

import (
	"bufio"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src"
	"os"
	"strings"
)

func main() {
	interpreter := src.NewInterpreter()
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		err := interpreter.ExecuteFile(filePath)
		if err != nil {
			fmt.Println(err)
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
			if code == "exit()" {
				break
			}

			err = interpreter.Execute(code)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
