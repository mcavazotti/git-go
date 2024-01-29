/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"mcavazotti/git-go/internal/objects"
	"mcavazotti/git-go/internal/repo"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/djherbis/times"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
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

		index, err := objects.ReadIndex(&r)
		if err != nil {
			panic(err)
		}

		compareIndexToHead(&r, &index)
		// fmt.Println("")
		compareIndexToWorktree(&r, &index)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func compareIndexToHead(repository *repo.Repository, index *objects.IndexObject) {
	shared.VerbosePrintln("")
	shared.VerbosePrintln("##################")
	shared.VerbosePrintln("compareIndexToHead")
	shared.VerbosePrintln("##################")
	shared.VerbosePrintln("")

	headSha, err := repository.Resolve("HEAD")
	shared.VerbosePrintln("")
	shared.VerbosePrintln("##################")
	shared.VerbosePrintln("found HEAD")
	shared.VerbosePrintf("%q\n", headSha)
	shared.VerbosePrintln("##################")
	shared.VerbosePrintln("")

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No commits yet")
		return
	}

	if err != nil {
		panic(err)
	}
	commit, err := objects.ReadCommit(repository, headSha)
	if err != nil {
		panic(err)
	}

	headEntries, err := objects.FlattenTree(repository, commit.Tree)
	if err != nil {
		panic(err)
	}

	var stagedChanges []string

	for _, indexEntry := range index.Entries {
		sha, exists := headEntries[indexEntry.Name]

		if exists {
			if sha != hex.EncodeToString(indexEntry.Sha) {
				stagedChanges = append(stagedChanges, fmt.Sprint("  modified: ", indexEntry.Name))
			}

			delete(headEntries, indexEntry.Name)
		} else {
			stagedChanges = append(stagedChanges, fmt.Sprint("  new file: ", indexEntry.Name))
		}
	}

	if len(stagedChanges) > 0 {
		fmt.Println("Changes to be commited:")
		for _, ln := range stagedChanges {
			println(ln)
		}
	}

	for filename := range headEntries {
		fmt.Println("  deleted:  ", filename)
	}
}

func compareIndexToWorktree(repository *repo.Repository, index *objects.IndexObject) {

	fmt.Println("Changes not staged for commit:")
	var allFiles []string

	filepath.Walk(repository.WorkTree, func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(path, filepath.FromSlash(repository.GitDir)) || path == repository.WorkTree || info.IsDir() {
			return nil
		}
		allFiles = append(allFiles, path)
		return nil
	})

	for _, entry := range index.Entries {
		fullPath := path.Join(repository.WorkTree, entry.Name)

		_, err := os.Stat(fullPath)

		if os.IsNotExist(err) {
			fmt.Println("  deleted:  ", entry.Name)
		} else {
			t, err := times.Stat(fullPath)

			if err != nil {
				panic(err)
			}

			if entry.Mtime_s != uint32(t.ModTime().Unix()) {
				hash, err := objects.HashFile(fullPath)
				if err != nil {
					panic(err)
				}
				if hex.EncodeToString(entry.Sha) != hash {
					fmt.Println("  modified:", entry.Name)
				}
			}
		}
		idx := slices.Index(allFiles, filepath.FromSlash(fullPath))
		if idx != -1 {
			allFiles = slices.Delete(allFiles, idx, idx+1)
		}
	}

	fmt.Println()

	ignore, err := objects.ReadGitIgnore(repository)
	if err != nil {
		panic(err)
	}

	fmt.Println("Untracked files:")
	for _, f := range allFiles {
		if !ignore.IgnoreFile(f) {
			fmt.Println(" ", filepath.ToSlash(strings.TrimPrefix(f, repository.WorkTree))[1:])
		}
	}
	fmt.Println()

}
