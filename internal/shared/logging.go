package shared

import (
	"fmt"
)

var verbose bool = false

func VerboseMode() {
	verbose = true
}

func VerbosePrint(message string) {
	if verbose {
		fmt.Println(message)
	}
}
