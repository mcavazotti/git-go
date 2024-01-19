package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
)

func HashData(data *[]byte) (string, error) {
	hash := sha1.Sum(*data)
	return hex.EncodeToString(hash[:]), nil
}

func HashFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return HashData(&data)
}

func CreateObjectData(data *[]byte, objType string) ([]byte, error) {
	header := []byte(objType + " " + fmt.Sprintf("%d", len(*data)))
	header = append(header, 0)
	content := append(header, (*data)...)

	var buffer bytes.Buffer
	w := zlib.NewWriter(&buffer)

	var err error

	if _, err := w.Write(content); err != nil {
		return nil, err
	}
	err = w.Close()

	return buffer.Bytes(), err
}

func WriteObject(p string, data *[]byte, objType string) (string, error) {
	shared.VerbosePrintln("Write Object:", objType)
	hash, err := HashData(data)
	shared.VerbosePrintln("Object SHA:", hash)

	if err != nil {
		return hash, err
	}

	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return hash, err
	}

	compressedObj, err := CreateObjectData(data, objType)
	if err != nil {
		return hash, err
	}
	err = os.WriteFile(path.Join(p, hash[2:]), compressedObj, os.ModePerm)
	return hash, err
}

func ReadObject(objPath string) (GitObject, error) {

	compressedData, err := os.ReadFile(objPath)
	if err != nil {
		return GitObject{}, err
	}

	b := bytes.NewReader(compressedData)
	reader, err := zlib.NewReader(b)
	if err != nil {
		return GitObject{}, err
	}

	uncompressedData, err := io.ReadAll(reader)
	if err != nil {
		reader.Close()
		return GitObject{}, err
	}

	reader.Close()

	var separatorIdx int
	for i := 0; i < len(uncompressedData); i++ {
		if uncompressedData[i] == 0x0 {
			separatorIdx = i
			break
		}
	}
	header := uncompressedData[:separatorIdx]

	var spaceIdx int
	for i := 0; i < len(header); i++ {
		if header[i] == 0x0 {
			spaceIdx = i
			break
		}
	}
	return GitObject{Data: uncompressedData[separatorIdx+1:], ObjType: string(header[:spaceIdx])}, nil
}
