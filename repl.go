package main

import (
	"strings"
)

func cleanInput(text string) []string {
	var output []string
	lowerText := strings.ToLower(text)
	output = strings.Fields(lowerText)
	return output
}
