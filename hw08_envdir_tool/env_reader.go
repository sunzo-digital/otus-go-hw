package main

import (
	"bytes"
	"os"
	"strings"
)

var terminalZero, lineBreak []byte = []byte{0}, []byte{10}

type Environment map[string]EnvValue

func (e Environment) Set() error {
	for key, env := range e {
		if _, isSet := os.LookupEnv(key); isSet {
			if err := os.Unsetenv(key); err != nil {
				return err
			}
		}

		// если установлен этот флаг, то нужно только удалить переменную без установки нового значения
		if env.NeedRemove {
			continue
		}

		if err := os.Setenv(key, env.Value); err != nil {
			return err
		}
	}

	return nil
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment, len(dirEntries))

	for _, entry := range dirEntries {
		filename := entry.Name()

		content, err := os.ReadFile(dir + "/" + filename)
		if err != nil {
			return nil, err
		}

		if len(content) == 0 {
			environment[filename] = EnvValue{NeedRemove: true}
			continue
		}

		content = trimAllAfterLineBreak(content)

		value := strings.TrimRight(string(content), "\t ")

		environment[filename] = EnvValue{Value: value}
	}

	return environment, nil
}

func trimAllAfterLineBreak(file []byte) []byte {
	lineBreakIndex := bytes.IndexByte(file, lineBreak[0])

	if lineBreakIndex >= 0 {
		file = file[:lineBreakIndex]
	}

	return bytes.ReplaceAll(file, terminalZero, lineBreak)
}
