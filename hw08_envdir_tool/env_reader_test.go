package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expected := make(Environment, 5)
	expected["BAR"] = EnvValue{Value: "bar"}
	expected["EMPTY"] = EnvValue{}
	expected["FOO"] = EnvValue{Value: "   foo\nwith new line"}
	expected["HELLO"] = EnvValue{Value: "\"hello\""}
	expected["UNSET"] = EnvValue{NeedRemove: true}

	actual, err := ReadDir("testdata/env")

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
