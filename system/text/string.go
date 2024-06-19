package text

import "strings"

type ZString struct{}

func NewString() *ZString {
	return &ZString{}
}

func (c *ZString) Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

func (c *ZString) Replace(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(origin, search, replace, n)
}

func (c *ZString) Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

func (c *ZString) Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func (c *ZString) SubStr(str string, start int, length ...int) (substr string) {
	lth := len(str)
	// Simple border checks.
	if start < 0 {
		start = 0
	}
	if start >= lth {
		// start = lth
		return ""
	}
	end := lth
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = lth
		}
	}
	if end > lth {
		end = lth
	}
	return str[start:end]
}

func (c *ZString) Join(array []string, sep string) string {
	return strings.Join(array, sep)
}

func (c *ZString) UcWords(str string) string {
	return strings.Title(str)
}

func (c *ZString) ToUpper(s string) string {
	return strings.ToUpper(s)
}
