/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"mcavazotti/git-go/internal/objects"
	"mcavazotti/git-go/internal/repo"
	"os"

	"github.com/spf13/cobra"
)

var (
	annotate bool
	force    bool
	message  string
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag [-a] [-m <message>] <tag> [<object>|<commit>]",
	Short: "A brief description of your command",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		r, err := repo.FindRepo(wd)
		if err != nil {
			panic(err)
		}
		var ref string

		if len(args) != 2 {
			ref = "HEAD"
		} else {
			ref = args[1]
		}
		if err := objects.CreateTag(&r, args[0], ref, force, annotate, message); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	tagCmd.Flags().BoolVarP(&annotate, "annotate", "a", false, "Make an unsigned, annotated tag object")
	tagCmd.Flags().BoolVarP(&force, "force", "f", false, "Replace an existing tag with the given name (instead of failing)")
	tagCmd.Flags().StringVarP(&message, "message", "m", "", "Use the given tag message")

}
