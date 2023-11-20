package cmd

import (
	"mcavazotti/git-go/internal/shared"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-go",
	Short: "A simple clone of Git written in GO",
	Long: `This program has the capabilities to:
	- create a repository
	- commit files (to be implemented)
	- push to GitHub (to be implemented)
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			shared.VerboseMode()
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Turn on verbose mode")
}
