package objects

import (
	"mcavazotti/git-go/internal/core"
	"mcavazotti/git-go/internal/repo"
	"path"
)

func WriteObject(repository *repo.Repository, data *[]byte, objType string) (string, error) {

	folder := path.Join(repository.GitDir, "objects")
	return core.WriteObject(folder, data, objType)
}

func ReadObject(repository *repo.Repository, sha string) (core.GitObject, error) {
	objPath, err := repository.FindObject(sha)
	if err != nil {
		return core.GitObject{}, err
	}
	return core.ReadObject(objPath)
}

var HashFile = core.HashFile
