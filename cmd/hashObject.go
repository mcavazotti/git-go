package cmd

import (
	"fmt"
	"mcavazotti/git-go/internal/objects"
	"mcavazotti/git-go/internal/repo"
	"mcavazotti/git-go/internal/shared"
	"os"

	"github.com/spf13/cobra"
)

var objType = "blob"
var write = false

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [-w] [-t <type>] <file>",
	Short: "Compute object ID and optionally creates a blob from a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shared.VerbosePrint("Write flag: " + fmt.Sprint(write))
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		r, err := repo.FindRepo(wd)
		if err != nil {
			panic(err)
		}

		sha, err := objects.HashFile(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println(sha)

		if write {
			shared.VerbosePrint("Write")
			switch objType {
			case "blob":
				shared.VerbosePrint("blob")
				if err := objects.WriteBlob(&r, args[0]); err != nil {
					panic(err)
				}
			default:
				panic("Unknown type (maybe not implemented)")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "Actually write the object into the database")
	hashObjectCmd.Flags().StringVarP(&objType, "type", "t", "blob", "Specify the type")
}
