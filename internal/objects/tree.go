package objects

import (
	"bytes"
	"io"
	"mcavazotti/git-go/internal/repo"
)

type TreeEntry struct {
	mode []byte
	path string
	sha  []byte
}

type TreeObject []TreeEntry

func WriteTree(repository *repo.Repository, tree TreeObject) error {
	var treeData []byte

	for _, entry := range tree {
		treeData = append(treeData, entry.mode...)
		treeData = append(treeData, ' ')
		treeData = append(treeData, []byte(entry.path)...)
		treeData = append(treeData, 0x00)
		treeData = append(treeData, entry.sha...)
	}

	return WriteObject(repository, &treeData, "tree")
}

func ReadTree(repository *repo.Repository, sha string) (TreeObject, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return TreeObject{}, err
	}

	tree := TreeObject{}
	r := bytes.NewReader(obj.data)

	for entry, err := readEntry(r); err != io.EOF; {
		if err != nil {
			return TreeObject{}, err
		}
		tree = append(tree, entry)
	}
	return tree, nil
}

func readEntry(reader *bytes.Reader) (TreeEntry, error) {
	entry := TreeEntry{}

	for b, err := reader.ReadByte(); b != ' '; {
		if err != nil {
			return TreeEntry{}, err
		}
		entry.mode = append(entry.mode, b)
	}

	pathBuffer := bytes.NewBufferString("")
	for b, err := reader.ReadByte(); b != 0x00; {
		if err != nil {
			return TreeEntry{}, err
		}
		pathBuffer.WriteByte(b)
	}
	entry.path = pathBuffer.String()

	sha := make([]byte, 20)
	_, err := reader.Read(sha)
	if err != nil {
		return TreeEntry{}, err
	}
	entry.sha = sha

	return entry, nil
}
