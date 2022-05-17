package strepeat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	tests := []struct {
		strToRepeat StringToRepeat
		expected    string
	}{
		{New("a", 1), "a"},
		{New("b", 5), "bbbbb"},
		{New("c", 0), ""},
	}

	for _, tc := range tests {
		tcName := fmt.Sprintf("%s * %d", tc.strToRepeat.str, tc.strToRepeat.repeatCount)
		t.Run(tcName, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.strToRepeat.ToString())
		})
	}
}
