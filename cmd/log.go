/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [<commit>...]",
	Short: "Shows the commit logs.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To be implemented")
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
