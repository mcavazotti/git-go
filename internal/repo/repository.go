package repo

import (
	"errors"
	"os"

	"gopkg.in/ini.v1"
)

type Repository struct {
	GitDir   string
	WorkTree string
	Config   *ini.File
}

func New(path string) Repository {
	workTree := path
	gitDir := workTree + "/.git"

	if _, err := os.Stat(gitDir); errors.Is(err, os.ErrNotExist) {
		panic("Not a Git repository: " + path)
	}

	config, err := ini.Load(gitDir + "/config")

	if err != nil {
		panic(err)
	}

	repo := Repository{GitDir: gitDir, WorkTree: workTree, Config: config}
	return repo
}
