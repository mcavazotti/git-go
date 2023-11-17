package repo

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

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
func RepoPath(r *Repository, pathSegments ...string) string {
	p := []string{r.GitDir}
	p = append(p, pathSegments...)
	return path.Join(p...)
}
