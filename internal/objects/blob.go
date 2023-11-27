package objects

import (
	"mcavazotti/git-go/internal/repo"
	"os"
)

func WriteBlob(repository *repo.Repository, filePath string) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return WriteObject(repository, &data, "blob")
}

func ReadBlob(repository *repo.Repository, sha string) ([]byte, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return []byte{}, err
	}

	return obj.data, nil
}
