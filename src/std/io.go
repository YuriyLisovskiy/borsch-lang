package std

import (
	"fmt"
	"strings"
)

func Print(args... string) (string, error) {
	fmt.Print(strings.Join(args, " "))
	return "", nil
}

func Println(args... string) (string, error) {
	fmt.Println(strings.Join(args, " "))
	return "", nil
}
