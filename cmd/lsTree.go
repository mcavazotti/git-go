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
		var t string
		switch fmt.Sprintf("%06o", entry.Mode)[:2] {
		case "04":
			t = "tree"
		case "10":
			t = "blob"
		case "12":
			t = "blob"
		case "16":
			t = "commit"
		default:
			panic("Unknown tree leaf mode")
		}

		if t == "tree" && recursive {
			listTree(r, hex.EncodeToString(entry.Sha), entry.Path)
		} else {
			fmt.Printf("%06o %s %s\t%s\n", entry.Mode, t, hex.EncodeToString(entry.Sha), entry.Path)
		}

	}
	return nil
}
