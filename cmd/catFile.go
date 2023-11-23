package cmd

import (
	"fmt"
	"mcavazotti/git-go/internal/objects"
	"mcavazotti/git-go/internal/repo"
	"os"

	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:       "cat-file <type> <object>",
	Short:     "Provide content of repository objects",
	ValidArgs: []string{"blob", "commit", "tag", "tree"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cobra.ExactArgs(2)(cmd, args)
		}
		return cobra.OnlyValidArgs(cmd, args[:1])
	},
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		r, err := repo.FindRepo(wd)
		if err != nil {
			panic(err)
		}

		switch args[0] {
		case "blob":
			data, err := objects.ReadBlob(&r, args[1])
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", string(data))
		}

	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
