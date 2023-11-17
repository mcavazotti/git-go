package objects

import (
	"mcavazotti/git-go/internal/repo"
	"os"
	"path"
)

func WriteBlob(repository repo.Repository, filePath string) error {
	hash, err := HashFile(filePath)

	if err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	folder := path.Join(repository.GitDir, "objects", hash[:2])
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	compressedObj, err := CreateObjectData(&data, "blob")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(folder, hash[2:]), compressedObj, os.ModePerm)
	return err
}
