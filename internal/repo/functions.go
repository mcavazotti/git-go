package repo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func (r Repository) FindObject(object string) (string, error) {
	p := r.RepoPath( "objects", object[:2], object[2:])
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("Not a valid object name %s", object)
	}
	return p, nil
}

func FindRepo(currPath string) (Repository, error) {
	if _, err := os.Stat(path.Join(currPath, ".git")); errors.Is(err, os.ErrNotExist) {
		parent, err := filepath.Abs(path.Join(currPath, ".."))
		if err != nil {
			panic(err)
		}
		if parent == currPath {
			return Repository{}, errors.New("No git repository.")
		}
		return FindRepo(parent)

	} else {
		return New(currPath), nil
	}
}
func (r Repository) RepoPath(pathSegments ...string) string {
	p := []string{r.GitDir}
	p = append(p, pathSegments...)
	return path.Join(p...)
}
