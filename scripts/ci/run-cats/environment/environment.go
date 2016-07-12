package environment

import (
	"fmt"
	"os"
)

type environment struct{}

func New() *environment {
	return &environment{}
}

func (e *environment) GetBoolean(varName string) (bool, error) {
	switch os.Getenv(varName) {
	case "true":
		return true, nil
	case "false", "":
		return false, nil
	default:
		return false, fmt.Errorf("Invalid environment variable: '%s' must be a boolean 'true' or 'false'", varName)
	}
}
