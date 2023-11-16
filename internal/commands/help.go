package commands

import (
	"fmt"
	"os"
	"strings"
)

func HelpMessage() {
	programPath := strings.Split(os.Args[0], "\\")
	if len(programPath) == 1 {
		programPath = strings.Split(os.Args[0], "/")
	}

	fmt.Printf("\nUSAGE:\n")
	fmt.Printf("> %s <commands> [<args>] [--verbose]", programPath[len(programPath)-1])
	fmt.Printf("\nCOMMANDS:\n")
	fmt.Printf("\tinit [<path>]\n")
	fmt.Printf("\thelp\n")
}
