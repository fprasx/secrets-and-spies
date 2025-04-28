package utils

import (
	"github.com/charmbracelet/log"
)

func Assert(cond bool, msg string) {
	if !cond {
		log.Fatalf("Assertion failed: %v", msg)
	}
}
