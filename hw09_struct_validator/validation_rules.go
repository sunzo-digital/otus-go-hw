package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

var supportedRules = map[string]struct {
	paramsCount   int
	kind          reflect.Kind
	makeValidator func(value reflect.Value, params []string) (validator, error)
}{
	"len": {
		paramsCount: 1,
		makeValidator: func(value reflect.Value, params []string) (validator, error) {
			expectedLen, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, err
			}

			//nolint:exhaustive
			switch value.Kind() {
			case reflect.String:
				return func(value reflect.Value) error {
					if len(value.String()) != expectedLen {
						return fmt.Errorf("length of value is not equal to %d", expectedLen)
					}

					return nil
				}, nil
			case reflect.Slice:
				return func(value reflect.Value) error {
					for i := 0; i < value.Len(); i++ {
						if len(value.Index(i).String()) != expectedLen {
							return fmt.Errorf("length of value is not equal to %d", expectedLen)
						}
					}

					return nil
				}, nil
			default:
				panic("wrong value type")
			}
		},
		kind: reflect.String,
	},
	"regexp": {
		paramsCount: 1,
		makeValidator: func(value reflect.Value, params []string) (validator, error) {
			r, err := regexp.Compile(params[0])
			if err != nil {
				return nil, err
			}

			return func(value reflect.Value) error {
				if r.Match([]byte(value.String())) {
					return nil
				}

				return errors.New("value does not match the regular expression")
			}, nil
		},
		kind: reflect.String,
	},
	"in": {
		paramsCount: 0,
		makeValidator: func(value reflect.Value, params []string) (validator, error) {
			//nolint:exhaustive
			switch value.Kind() {
			case reflect.String:
				return func(value reflect.Value) error {
					for _, acceptable := range params {
						if value.String() == acceptable {
							return nil
						}
					}

					return fmt.Errorf("value is not among the valid values: %v", params)
				}, nil
			case reflect.Int:
				return func(value reflect.Value) error {
					for _, param := range params {
						acceptable, _ := strconv.Atoi(param)
						if value.Int() == int64(acceptable) {
							return nil
						}
					}

					return fmt.Errorf("value is not among the valid values: %v", params)
				}, nil
			default:
				panic("wrong value type")
			}
		},
		kind: reflect.String,
	},
	"min": {
		paramsCount: 1,
		makeValidator: func(value reflect.Value, params []string) (validator, error) {
			min, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, errors.New("param must be numeric")
			}

			return func(value reflect.Value) error {
				if value.Int() < int64(min) {
					return fmt.Errorf("value is less than %d", min)
				}

				return nil
			}, nil
		},
		kind: reflect.Int,
	},
	"max": {
		paramsCount: 1,
		makeValidator: func(value reflect.Value, params []string) (validator, error) {
			max, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, errors.New("param must be numeric")
			}

			return func(value reflect.Value) error {
				if value.Int() > int64(max) {
					return fmt.Errorf("value is greater than %d", max)
				}

				return nil
			}, nil
		},
		kind: reflect.Int,
	},
}
