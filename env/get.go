package env

import (
	"errors"
	"fmt"
	"os"
)

func Get(name string) (string, error) {
	val, ok := os.LookupEnv(name)
	if ok {
		return val, nil
	} else {
		return "", errors.New(fmt.Sprintf("Environment variable `%s` not found", name))
	}
}
