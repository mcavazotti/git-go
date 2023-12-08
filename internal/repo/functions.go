package repo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func (r Repository) FindObject(object string) (string, error) {
	p := r.RepoPath("objects", object[:2], object[2:])
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

func (r Repository) Resolve(s string) (string, error) {
	if s == "HEAD" {
		b, err := os.ReadFile(r.RepoPath("HEAD"))
		if err != nil {
			return "", err
		}
		s = string(b[5:])
	}

	if _, err := hex.DecodeString(s); err != nil {
		return r.ResolveRef(s)
	}

	return s, nil
}

func (r Repository) ResolveRef(ref string) (string, error) {
	path := r.RepoPath(ref)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	strData := string(data)
	strData = strData[:len(strData)-1]
	if strData[0] == 'r' {
		return r.ResolveRef(strData[5:])
	}
	return strData, nil
}
