package commands

import (
	"fmt"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/ini.v1"
)

func InitRepo(args []string) {
	path := "."
	if len(args) >= 1 {
		path = args[0]
	}
	createRepository(path)
}

func createRepository(repoPath string) {
	var err error
	wd, _ := os.Getwd()

	err = os.Mkdir(repoPath+"/.git", os.ModePerm)
	shared.VerbosePrint("CREATE " + filepath.Clean(path.Join(wd, repoPath, ".git")))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(repoPath+"/.git/objects", os.ModePerm)
	shared.VerbosePrint("CREATE " + filepath.Clean(path.Join(wd, repoPath, ".git", "objects")))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(repoPath+"/.git/refs/heads", os.ModePerm)
	shared.VerbosePrint("CREATE " + filepath.Clean(path.Join(wd, repoPath, ".git", "refs", "heads")))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(repoPath+"/.git/refs/tags", os.ModePerm)
	shared.VerbosePrint("CREATE " + filepath.Clean(path.Join(wd, repoPath, ".git", "refs", "tags")))
	if err != nil {
		panic(err)
	}
	head, err := os.Create(repoPath + "/.git/HEAD")
	shared.VerbosePrint("CREATE " + filepath.Clean(path.Join(wd, repoPath, ".git", "HEAD")))
	if err != nil {
		panic(err)
	}
	description, err := os.Create(repoPath + "/.git/description")
	shared.VerbosePrint("CREATE " + filepath.Clean(path.Join(wd, repoPath, ".git", "description")))
	if err != nil {
		panic(err)
	}

	// .git/description
	shared.VerbosePrint("WRITE TO " + filepath.Clean(path.Join(wd, repoPath, ".git", "description")))
	_, err = fmt.Fprint(description, "Unnamed repository; edit this file 'description' to name the repository.\n")
	if err != nil {
		panic(err)
	}

	// .git/HEAD
	shared.VerbosePrint("WRITE TO " + filepath.Clean(path.Join(wd, repoPath, ".git", "HEAD")))
	_, err = fmt.Fprint(head, "ref: refs/heads/master\n")
	if err != nil {
		panic(err)
	}

	// .git/config
	shared.VerbosePrint("CREATE " + filepath.Clean(wd+"/"+repoPath+"/.git/config"))
	createConfig(wd + "/" + repoPath + "/.git/config")

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
