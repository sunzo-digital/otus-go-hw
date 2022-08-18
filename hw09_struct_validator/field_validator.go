package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Rules map[string][]string

type validator func(value reflect.Value) error

func validateField(fieldInfo fieldInfo) error {
	rules, err := parseRules(fieldInfo.tag)
	if err != nil {
		return err
	}

	validators, err := validatorsFromRules(rules, fieldInfo.value)
	if err != nil {
		return err
	}

	validationErrors := make(ValidationErrors, 0)

	for _, v := range validators {
		err = v(fieldInfo.value)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldInfo.name,
				Err:   err,
			})
		}
	}

	return validationErrors
}

func parseRules(row string) (Rules, error) {
	rulePairs := strings.Split(row, "|")
	rules := make(Rules, len(rulePairs))
	for _, pair := range rulePairs {
		keyValues := strings.Split(pair, ":")
		if len(keyValues) != 2 {
			return nil, errors.New("invalid rules")
		}
		rules[keyValues[0]] = strings.Split(keyValues[1], ",")
	}
	return rules, nil
}

func validatorsFromRules(rules Rules, value reflect.Value) ([]validator, error) {
	validators := make([]validator, 0, len(rules))

	for name, params := range rules {
		ruleInfo, ok := supportedRules[name]

		if !ok {
			return nil, fmt.Errorf("unsuported rule: %s", name)
		}

		// если требуемое число параметров = 0, то можно передать неограниченное количество параметров, но не меньше 1
		if ruleInfo.paramsCount == 0 {
			if len(params) < 1 {
				return nil, fmt.Errorf("%s rule expect at least 1 parameter", name)
			}
		} else if len(params) != ruleInfo.paramsCount {
			return nil, fmt.Errorf(
				"%d parameters are expected for %s rule, %d were passed",
				ruleInfo.paramsCount,
				name,
				len(params),
			)
		}

		if value.Kind() != ruleInfo.kind {
			if value.Kind() != reflect.Slice {
				return nil, errors.New("unsupported kind")
			}

			switch value.Interface().(type) {
			case []string:
			case []int:
				func() {
					// допустимые типы, ничего не делаем
				}()
			default:
				return nil, fmt.Errorf("unsupported slice type: %T", value.Interface())
			}
		}

		validator, err := ruleInfo.makeValidator(value, params)
		if err != nil {
			return nil, fmt.Errorf("invalid params: %w", err)
		}

		validators = append(validators, validator)
	}

	return validators, nil
}
