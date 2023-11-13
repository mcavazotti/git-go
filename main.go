package gitgo

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("aqui")
	if len(os.Args) < 2 {
		helpMessage()
		return
	}

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	for _, v := range os.Args {
		if v == "--verbose" {
			setVerbose(true)
		}
	}

	switch os.Args[1] {
	case "help":
		helpMessage()
	case "init":
		initDir := "."
		if len(initCmd.Args()) > 0 {
			initDir = initCmd.Arg(0)
		}
		if err := createRepository(initDir); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid command: %s", os.Args[1])
		helpMessage()
		os.Exit(1)
	}

}

func helpMessage() {
	fmt.Printf("\nUSAGE:\n")
	fmt.Printf("> %s <commands>", os.Args[0])
}
