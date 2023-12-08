/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"mcavazotti/git-go/internal/repo"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

// showRefCmd represents the showRef command
var showRefCmd = &cobra.Command{
	Use:   "show-ref",
	Short: "List references.",

	Run: func(cmd *cobra.Command, args []string) {
		refs := listRefs()

		for _, r := range refs {
			fmt.Println(r[0], r[1])
		}
	},
}

func init() {
	rootCmd.AddCommand(showRefCmd)
}

func listRefs() [][2]string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	r, err := repo.FindRepo(wd)
	if err != nil {
		panic(err)
	}

	refsPath := r.RepoPath("refs")

	var refs [][2]string

	filepath.Walk(refsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ref, err := r.ResolveRef(path)
			if err != nil {
				return err
			}
			refs = append(refs, [2]string{ref, path})
		}
		return nil
	})

	sort.Slice(refs, func(i, j int) bool {
		return refs[i][1] < refs[j][1]
	})

	return refs
}
