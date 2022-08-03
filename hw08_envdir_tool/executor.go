package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

// Не придумал, как перехватывать поток вывода в тестах без этого костыля.
var (
	inStream  io.ReadWriter = os.Stdin
	outStream io.ReadWriter = os.Stdout
	errStream io.ReadWriter = os.Stderr
)

//nolint:gosec
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if err := env.Set(); err != nil {
		log.Fatal(err.Error())
	}

	var command *exec.Cmd

	if len(cmd) > 1 {
		command = exec.Command(cmd[0], cmd[1:]...)
	} else {
		command = exec.Command(cmd[0])
	}

	command.Stdin = inStream
	command.Stdout = outStream
	command.Stderr = errStream

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		log.Fatal(err.Error())
	}

	return 0
}
