package shared

import (
	"fmt"
)

var Verbose bool = false

func VerboseMode() {
	Verbose = true
}

func VerbosePrintf(format string, a ...any) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}
func VerbosePrintln(a ...any) {
	if Verbose {
		fmt.Println(a...)
	}
}
