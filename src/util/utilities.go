package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

func CalcHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

func ReadFile(filePath string) (content []byte, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = RuntimeError(fmt.Sprintf("файл з ім'ям '%s' не існує", filePath))
		return
	}

	content, err = ioutil.ReadFile(filePath)
	return
}
