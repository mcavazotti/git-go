package main

import (
	"flag"
	"fmt"
	"mcavazotti/git-go/internal/commands"
	"mcavazotti/git-go/internal/shared"
	"os"
)

func main() {

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	verbose := flag.Bool("verbose mode", false, "Activate verbose messages")
	if *verbose {
		fmt.Print("verbose")
		shared.VerboseMode()
	}

	if len(os.Args) < 2 {
		commands.HelpMessage()
		return
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
