package strepeat

import (
	"strings"
)

type StringToRepeat struct {
	str         string
	repeatCount int
}

func New(str string, repeatCount int) StringToRepeat {
	return StringToRepeat{
		str:         str,
		repeatCount: repeatCount,
	}
}

func (rs StringToRepeat) ToString() string {
	if rs.repeatCount == 0 {
		return ""
	}

	return strings.Repeat(rs.str, rs.repeatCount)
}

func (rs *StringToRepeat) SetRepeatCount(count int) {
	rs.repeatCount = count
}
