package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) (content []byte, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = RuntimeError(fmt.Sprintf("файл з ім'ям '%s' не існує", filePath))
		return
	}

	content, err = ioutil.ReadFile(filePath)
	return
}
