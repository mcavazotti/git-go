package repo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"mcavazotti/git-go/internal/core"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (r Repository) FindObject(object string) (string, error) {
	sha, err := r.Resolve(object)
	shared.VerbosePrint("Object SHA: " + sha)
	if err != nil {
		return "", err
	}

	p := r.RepoPath("objects", sha[:2], sha[2:])
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		shared.VerbosePrint("Error> " + err.Error())
		return "", fmt.Errorf("Not a valid object name %s", sha)
	}
	shared.VerbosePrint("Found object: " + p)
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
	shared.VerbosePrint("RepoPath> " + path.Join(p...))
	return path.Join(p...)
}

func (r Repository) Resolve(s string) (string, error) {
	shared.VerbosePrint("Resolving: " + s)
	if s == "HEAD" {
		b, err := os.ReadFile(r.RepoPath("HEAD"))
		if err != nil {
			return "", err
		}
		s = string(b[5:])
	}

	if _, err := hex.DecodeString(s); err != nil {
		tag, errTag := r.ResolveRef(path.Join("refs", "tags", s))
		branch, errBranch := r.ResolveRef(path.Join("refs", "heads", s))
		ref, err := r.ResolveRef(s)

		if errTag != nil && errBranch != nil && err != nil {
			return "", err
		}

		if (tag != "" && branch != "") || (tag != "" && ref != "") || (branch != "" && ref != "") {
			return "", fmt.Errorf("refname '%s' is ambiguous.", s)
		}

		if errTag == nil {
			p, err := r.FindObject(tag)
			if err != nil {
				return "", err
			}
			obj, err := core.ReadObject(p)
			if err != nil {
				return "", err
			}
			if obj.ObjType == "tag" {
				tag = string(obj.Data[7:47])
			}
			return tag, nil
		} else if errBranch == nil {
			return branch, nil
		}

		return ref, err
	}

	dir, err := os.ReadDir(r.RepoPath("objects", s[:2]))

	if err != nil {
		return "", err
	}

	var candidates []string

	for _, f := range dir {
		if s[2:] == f.Name()[:len(s[2:])] {
			candidates = append(candidates, s[:2]+f.Name())
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("Failed to resolve '%s' as a valid ref.", s)
	}

	if len(candidates) > 1 {
		return "", fmt.Errorf("Short object ID 0a4b is ambiguous.\nFailed to resolve '%s' as a valid ref.", s)
	}

	return candidates[0], nil
}

func (r Repository) ResolveRef(ref string) (string, error) {
	shared.VerbosePrint("Resolving reference: " + ref)
	path := r.RepoPath(ref)

	data, err := os.ReadFile(path)
	if err != nil {
		shared.VerbosePrint("Error> " + err.Error())
		return "", err
	}
	strData := string(data)
	if strData[0] == 'r' {
		return r.ResolveRef(strData[5:])
	}
	shared.VerbosePrint("Resolved: " + strData)
	return strData, nil
}

func (r Repository) GetActiveBranch() (string, bool, error) {
	data, err := os.ReadFile(path.Join(r.GitDir, "HEAD"))
	if err != nil {
		return "", false, err
	}
	if strings.HasPrefix(string(data), "ref: refs/heads/") {
		return string(data)[16:], true, nil
	} else {
		return string(data), false, nil
	}
}
