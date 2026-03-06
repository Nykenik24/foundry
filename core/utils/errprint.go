package utils

import (
	"fmt"
	"os"
)

func PrintError(err string, args ...any) {
	fmt.Printf("\033[31mFoundry error:\033[0m\n\t%s\n", fmt.Sprintf(err, args...))
}

func PrintFatal(err string, args ...any) {
	PrintError(err, args...)
	os.Exit(1)
}
