package configuration

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type EnvironmentVariable struct {
	Key       string
	Validator func(string, string) error
}

func (ev *EnvironmentVariable) Validate(value string) error {
	return ev.Validator(ev.Key, value)
}

func isNotBlank(key string, value string) error {
	if strings.Trim(value, " ") == "" {
		return errors.New(key + " variable is not set.")
	}
	return nil
}

func isNumber(key string, value string) error {
	if strings.Trim(value, " ") == "" {
		return errors.New(key + " variable is not set.")
	}
	_, err := strconv.Atoi(value)
	if err != nil {
		return errors.Wrapf(err, "Variable %s is not a number (environment value is %s)", key, value)
	}
	return nil
}

type EnvironmentReader struct {
	Errors []error
}

func (e *EnvironmentReader) Read(variable *EnvironmentVariable) string {
	value := os.Getenv(variable.Key)
	if err := variable.Validate(value); err != nil {
		e.Errors = append(e.Errors, err)
	}
	return value
}

func NewEnvReader() *EnvironmentReader {
	return &EnvironmentReader{make([]error, 0)}
}
