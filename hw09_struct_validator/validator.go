package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field \"%s\" validation error: %s", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errorMessages strings.Builder

	errorMessages.WriteString("Validation errors:\n")

	for _, err := range v {
		errorMessages.WriteString(err.Error() + "\n")
	}

	return errorMessages.String()
}

type fieldInfo struct {
	name  string
	value reflect.Value
	tag   string
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return errors.New("received entity is not a structure")
	}

	queue := make(chan fieldInfo, value.NumField())
	go fillQueue(value, queue)

	validationErrors := make(ValidationErrors, 0)

	for fieldInfo := range queue {
		err := validateField(fieldInfo)
		fieldValidationErrors, ok := err.(ValidationErrors) //nolint:errorlint
		if !ok {
			return err
		}

		if len(fieldValidationErrors) > 0 {
			validationErrors = append(validationErrors, fieldValidationErrors...)
		}
	}

	return validationErrors
}

func fillQueue(value reflect.Value, queue chan fieldInfo) {
	defer close(queue)

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		queue <- fieldInfo{
			name:  field.Name,
			value: value.Field(i),
			tag:   validateTag,
		}
	}
}
