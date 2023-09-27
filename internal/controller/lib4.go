package controller

import (
	"fmt"
	"regexp"
	"strings"
)

type CustomError struct {
	Message string
	Code    int
}

// Error returns a string representation of the error.
func (ce CustomError) Error() string {
	return fmt.Sprintf("CustomError: %s (Code: %d)", ce.Message, ce.Code)
}

func removeWordFromLast(input string, wordToRemove string) string {
	// Create a regular expression pattern to match the word at the end of the string.
	pattern := regexp.MustCompile(wordToRemove + `\s*$`)

	// Replace the matched word with an empty string.
	result := pattern.ReplaceAllString(input, "")
	return result
}

func removeTextdFromLast(input string, textToRemove string) string {
	index := strings.LastIndex(input, textToRemove)

	if index != -1 {
		// Remove the text from the end of the string.
		result := input[:index] + input[index+len(textToRemove):]
		return result
	} else {
		// Text to remove not found in the string.
		return ""
	}
}
