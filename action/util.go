package action

import (
	"fmt"
	"unicode/utf8"
)

func FixedLengthString(s string, length int) string {
	if utf8.RuneCountInString(s) >= length {
		return s[:length]
	}

	return fmt.Sprintf("%-*s", length, s)
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {

		if !encountered[elements[v]] {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func TruncateString(s string, length int) string {
	if utf8.RuneCountInString(s) >= length {
		return s[:length] + "..."
	}

	return s
}
