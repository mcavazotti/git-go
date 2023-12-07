/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"
	"mcavazotti/git-go/internal/objects"
	"mcavazotti/git-go/internal/repo"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var recursive = false

// lsTreeCmd represents the lsTree command
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		r, err := repo.FindRepo(wd)
		if err != nil {
			panic(err)
		}

		if err := listTree(&r, args[0], ""); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsTreeCmd)

	lsTreeCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recurse into sub-trees.")
}

func listTree(r *repo.Repository, sha string, prefix string) error {
	tree, err := objects.ReadTree(r, sha)
	if err != nil {
		return err
	}

	for _, entry := range tree {
		if entry.Mode.IsDir() {
			if recursive {
				listTree(r, hex.EncodeToString(entry.Sha), entry.Path)
			} else {
				fmt.Printf("%o tree %s\t%s\n", entry.Mode, hex.EncodeToString(entry.Sha), entry.Path)
			}
		} else {
			fmt.Printf("%o blob %s\t%s\n", entry.Mode, hex.EncodeToString(entry.Sha), path.Join(prefix, entry.Path))
		}

	}
	return nil
}
