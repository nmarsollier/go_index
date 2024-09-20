package builtin

import "strings"

func Shorten(name string) string {
	return fromString(name).
		split(" ").
		mapNotEmpty(func(s string) string {
			return strings.ToUpper(string(s[0]))
		}).
		joinToString()
}

type shorten struct{ string }
type shortenSlice struct{ slice []string }

func fromString(value string) shorten {
	return shorten{value}
}

func (s shorten) split(separator string) shortenSlice {
	return shortenSlice{strings.Split(s.string, separator)}
}

func (values shortenSlice) mapNotEmpty(f func(string) string) shortenSlice {
	var result []string
	for _, v := range values.slice {
		if len(v) > 0 {
			result = append(result, f(v))
		}
	}
	return shortenSlice{result}
}

func (values shortenSlice) joinToString() (result string) {
	for _, v := range values.slice {
		result += v
	}
	return result
}
