package builtin

import (
	"os"
)

func Exit(code int64) error {
	os.Exit(int(code))
	return nil
}
