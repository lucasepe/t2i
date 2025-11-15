package cmd

import (
	"fmt"
	"os"
)

// CheckErr check for an error and eventually exit
func CheckErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}
