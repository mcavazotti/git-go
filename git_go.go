package main

import (
	"flag"
	"fmt"
	"mcavazotti/git-go/internal/commands"
	"mcavazotti/git-go/internal/shared"
	"os"

	"golang.org/x/exp/slices"
)

func main() {
	if len(os.Args) < 2 {
		commands.HelpMessage()
		return
	}

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	if slices.Contains(os.Args, "--verbose") {
		shared.VerboseMode()
	}

	switch os.Args[1] {
	case "help":
		commands.HelpMessage()
	case "init":
		commands.InitRepo(initCmd.Args())
	default:
		fmt.Fprintf(os.Stderr, "Invalid command: %s", os.Args[1])
		commands.HelpMessage()
		os.Exit(1)
	}

}
