/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tag called")
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	tagCmd.Flags().BoolVarP(&annotate, "annotate", "a", false, "Make an unsigned, annotated tag object")
	tagCmd.Flags().BoolVarP(&force, "force", "f", false, "Replace an existing tag with the given name (instead of failing)")
	tagCmd.Flags().StringVarP(&message, "message", "m", "", "Use the given tag message")

}
