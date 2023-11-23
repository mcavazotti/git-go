package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
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
	return GitObject{data: uncompressedData[separatorIdx+1:], objType: string(header[:spaceIdx])}, nil
}
