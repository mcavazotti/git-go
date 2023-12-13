package objects

import (
	"mcavazotti/git-go/internal/repo"
	"mcavazotti/git-go/internal/shared"
	"os"
)

func WriteBlob(repository *repo.Repository, filePath string) error {
	shared.VerbosePrint("Write Blob: " + filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	_, err = WriteObject(repository, &data, "blob")
	return err
}

func ReadBlob(repository *repo.Repository, sha string) ([]byte, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return []byte{}, err
	}

	return obj.Data, nil
}
