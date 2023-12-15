/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"mcavazotti/git-go/internal/repo"
	"os"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status.",

	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		r, err := repo.FindRepo(wd)
		if err != nil {
			panic(err)
		}

		branch, attached, err := r.GetActiveBranch()
		if err != nil {
			panic(err)
		}
		if attached {
			fmt.Printf("On branch %s\n", branch)
		} else {
			fmt.Printf("HEAD detached at %s\n", branch)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
