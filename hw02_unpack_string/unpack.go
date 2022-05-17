package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	strepeat "github.com/sunzo-digital/otus-go-hw/hw02_unpack_string/string-repeat"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	firstRune := firstRuneOfString(input)

	if unicode.IsDigit(firstRune) {
		return "", ErrInvalidString
	}

	stringsToRepeat, err := stringsToRepeat(input)
	if err != nil {
		return "", err
	}

	strBuilder := strings.Builder{}

	for _, strToRepeat := range stringsToRepeat {
		strBuilder.WriteString(strToRepeat.ToString())
	}

	return strBuilder.String(), nil
}

func stringsToRepeat(input string) ([]strepeat.StringToRepeat, error) {
	stringsToRepeat := make([]strepeat.StringToRepeat, 0, len([]rune(input)))

	prev := firstRuneOfString(input)

	for _, current := range input {
		if unicode.IsLetter(current) {
			strToRepeat := strepeat.New(string(current), 1)
			stringsToRepeat = append(stringsToRepeat, strToRepeat)
			prev = current
			continue
		}

		// 2 цифры подряд
		if unicode.IsDigit(prev) {
			return nil, ErrInvalidString
		}

		prevStrToRepeat := &stringsToRepeat[len(stringsToRepeat)-1]
		repeatCount, err := strconv.Atoi(string(current))
		if err != nil {
			return nil, err
		}

		prevStrToRepeat.SetRepeatCount(repeatCount)
		prev = current
	}

	return stringsToRepeat, nil
}

func firstRuneOfString(input string) rune {
	return []rune(input)[0]
}
