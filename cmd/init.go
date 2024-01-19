package cmd

import (
	"fmt"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [<directory>]",
	Short: "Initialize a new, empty repository.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shared.VerbosePrintln("Run command Init")
		repoPath := "."
		if len(args) != 0 {
			repoPath = args[0]
		}
		createRepository(repoPath)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func createRepository(repoPath string) {
	var err error
	wd, _ := os.Getwd()

	err = os.Mkdir(path.Join(repoPath, ".git"), os.ModePerm)
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git")))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(path.Join(repoPath, ".git", "objects"), os.ModePerm)
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git", "objects")))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(path.Join(repoPath, ".git", "refs", "heads"), os.ModePerm)
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git", "refs", "heads")))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(path.Join(repoPath, ".git", "refs", "tags"), os.ModePerm)
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git", "refs", "tags")))
	if err != nil {
		panic(err)
	}
	head, err := os.Create(path.Join(repoPath, ".git", "HEAD"))
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git", "HEAD")))
	if err != nil {
		panic(err)
	}
	description, err := os.Create(path.Join(repoPath, ".git", "description"))
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git", "description")))
	if err != nil {
		panic(err)
	}

	// .git/description
	shared.VerbosePrintln("WRITE TO", filepath.Clean(path.Join(wd, repoPath, ".git", "description")))
	_, err = fmt.Fprint(description, "Unnamed repository; edit this file 'description' to name the repository.\n")
	if err != nil {
		panic(err)
	}

	// .git/HEAD
	shared.VerbosePrintln("WRITE TO", filepath.Clean(path.Join(wd, repoPath, ".git", "HEAD")))
	_, err = fmt.Fprint(head, "ref: refs/heads/master\n")
	if err != nil {
		panic(err)
	}

	// .git/config
	shared.VerbosePrintln("CREATE", filepath.Clean(path.Join(wd, repoPath, ".git", "config")))
	createConfig(path.Join(wd, repoPath, ".git", "config"))

}

func createConfig(path string) {
	iniData := ini.Empty()
	sec, err := iniData.NewSection("core")
	if err != nil {
		panic(err)
	}

	_, err = sec.NewKey("repositoryformatversion", "0")
	if err != nil {
		panic(err)
	}

	_, err = sec.NewKey("filemode", "false")
	if err != nil {
		panic(err)
	}

	_, err = sec.NewKey("bare", "false")
	if err != nil {
		panic(err)
	}

	err = iniData.SaveToIndent(path, "\t")
	if err != nil {
		panic(err)
	}
}
