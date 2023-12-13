/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"mcavazotti/git-go/internal/objects"
	"mcavazotti/git-go/internal/repo"
	"os"

	"github.com/spf13/cobra"
)

// lsFilesCmd represents the lsFiles command
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "Show information about files in the index and the working tree",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		r, err := repo.FindRepo(wd)
		if err != nil {
			panic(err)
		}

		indexObj, err := objects.ReadIndex(&r)
		if err != nil {
			panic(err)
		}

		for _, entry := range indexObj.Entries {
			fmt.Println(entry.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsFilesCmd)
}
