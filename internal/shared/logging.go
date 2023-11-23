package shared

import (
	"fmt"
)

var Verbose bool = false

func VerboseMode() {
	Verbose = true
}

func VerbosePrint(message string) {
	if Verbose {
		fmt.Println(message)
	}
}
