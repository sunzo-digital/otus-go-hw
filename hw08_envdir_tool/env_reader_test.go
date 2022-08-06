package main

import (
	"os"
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

func TestEnvironmentSetting(t *testing.T) {
	_ = os.Setenv("ORIGINAL", "MUST_STAY_ORIGINAL")
	_ = os.Setenv("BAR", "MUST_REPLACE")
	_ = os.Setenv("EMPTY", "MUST_REPLACE")
	_ = os.Setenv("UNSET", "MUST_UNSET")

	environment := make(Environment, 3)
	environment["BAR"] = EnvValue{Value: "bar"}
	environment["EMPTY"] = EnvValue{}
	environment["UNSET"] = EnvValue{NeedRemove: true}

	require.NoError(t, environment.Set())

	val, _ := os.LookupEnv("ORIGINAL")
	require.Equal(t, "MUST_STAY_ORIGINAL", val)

	val, _ = os.LookupEnv("BAR")
	require.Equal(t, "bar", val)

	val, isSet := os.LookupEnv("EMPTY")
	require.True(t, isSet)
	require.Equal(t, "", val)

	_, isSet = os.LookupEnv("UNSET")
	require.False(t, isSet)
}
