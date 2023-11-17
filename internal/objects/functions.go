package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

func HashFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:]), nil
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
