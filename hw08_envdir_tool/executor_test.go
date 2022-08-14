package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const goldenFilePath = "testdata/expected"

func TestRunCmd(t *testing.T) {
	b := new(bytes.Buffer)
	outStream = b

	currentDir, _ := os.Getwd()
	env, err := ReadDir(currentDir + "/testdata/env")
	require.NoError(t, err)

	returnCode := RunCmd([]string{"testdata/echo.sh"}, env)
	expected, _ := ioutil.ReadFile(goldenFilePath)

	require.Equal(t, 0, returnCode)
	require.Equal(t, string(expected), b.String())
}
