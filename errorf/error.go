package errorf

import (
	"fmt"
	"os"
)

func Error(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func Exit(s string, args ...interface{}) {
	Error(s, args...)
	os.Exit(1)
}
