package repo

import (
	"errors"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

type Repository struct {
	GitDir   string
	WorkTree string
	Config   *ini.File
}

func New(p string) Repository {
	workTree := p
	gitDir := path.Join(workTree, ".git")

	if _, err := os.Stat(gitDir); errors.Is(err, os.ErrNotExist) {
		panic("Not a Git repository: " + p)
	}

	config, err := ini.Load(path.Join(gitDir, "config"))

	if err != nil {
		panic(err)
	}

	repo := Repository{GitDir: gitDir, WorkTree: workTree, Config: config}
	return repo
}
