package main

import (
	"strings"
)

func process_rules(message string, rlog string) bool {
	return !strings.Contains(message, "AH01276")
	//return true
}
