package objects

import (
	"bytes"
	"fmt"
	"io"
	"mcavazotti/git-go/internal/repo"
	"os"
	"sort"
	"strconv"
)

type TreeEntry struct {
	Mode os.FileMode
	Path string
	Sha  []byte
}

type TreeObject []TreeEntry

func WriteTree(repository *repo.Repository, tree TreeObject) error {
	var treeData []byte

	for _, entry := range tree {
		treeData = append(treeData, fmt.Sprint(entry.Mode)...)
		treeData = append(treeData, ' ')
		treeData = append(treeData, entry.Path...)
		treeData = append(treeData, 0x00)
		treeData = append(treeData, entry.Sha...)
	}

	_, err := WriteObject(repository, &treeData, "tree")
	return err
}

func ReadTree(repository *repo.Repository, sha string) (TreeObject, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return TreeObject{}, err
	}
	tree := TreeObject{}
	r := bytes.NewReader(obj.Data)

	for entry, err := readTreeEntry(r); err != io.EOF; {
		if err != nil {
			return TreeObject{}, err
		}
		tree = append(tree, entry)
	}

	sort.Slice(tree, func(i, j int) bool {
		pathA := tree[i].Path
		pathB := tree[j].Path
		if tree[i].Mode.IsDir() {
			pathA += "/"
		}
		if tree[j].Mode.IsDir() {
			pathB += "/"
		}
		return pathA < pathB
	})

	return tree, nil
}

func readTreeEntry(reader *bytes.Reader) (TreeEntry, error) {
	entry := TreeEntry{}
	var mode string
	for b, err := reader.ReadByte(); b != ' '; {
		if err != nil {
			return TreeEntry{}, err
		}
		mode += string(b)
	}
	modeVal, _ := strconv.ParseInt(string(mode), 8, 32)
	entry.Mode = os.FileMode(uint32(modeVal))

	pathBuffer := bytes.NewBufferString("")
	for b, err := reader.ReadByte(); b != 0x00; {
		if err != nil {
			return TreeEntry{}, err
		}
		pathBuffer.WriteByte(b)
	}
	entry.Path = pathBuffer.String()

	sha := make([]byte, 20)
	_, err := reader.Read(sha)
	if err != nil {
		return TreeEntry{}, err
	}
	entry.Sha = sha

	return entry, nil
}
