package objects

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"mcavazotti/git-go/internal/repo"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
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
	shared.VerbosePrintln("ReadTree>", sha)
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return TreeObject{}, err
	}
	tree := TreeObject{}
	r := bytes.NewReader(obj.Data)
	// shared.VerbosePrintln(string(obj.Data))

	entry, err := readTreeEntry(r)
	for ; err != io.EOF; entry, err = readTreeEntry(r) {
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
	shared.VerbosePrintln("readTreeEntry")
	entry := TreeEntry{}
	var mode string
	b, err := reader.ReadByte()
	for ; b != ' '; b, err = reader.ReadByte() {

		if err != nil {
			return TreeEntry{}, err
		}
		mode += string(b)
	}
	modeVal, _ := strconv.ParseInt(string(mode), 8, 32)
	entry.Mode = os.FileMode(uint32(modeVal))
	shared.VerbosePrintf("MODE %s\n", entry.Mode.String())

	pathBuffer := bytes.NewBufferString("")
	b, err = reader.ReadByte()
	for ; b != 0x00; b, err = reader.ReadByte() {
		if err != nil {
			return TreeEntry{}, err
		}
		pathBuffer.WriteByte(b)
	}
	entry.Path = pathBuffer.String()
	shared.VerbosePrintf("PATH %s\n", entry.Path)

	sha := make([]byte, 20)
	_, err = reader.Read(sha)
	if err != nil {
		return TreeEntry{}, err
	}
	entry.Sha = sha

	return entry, nil
}

func FlattenTree(repository *repo.Repository, sha string) (map[string]string, error) {
	shared.VerbosePrintln("FlatenTree> ", sha)
	tree, err := ReadTree(repository, sha)

	if err != nil {
		return nil, err
	}

	entries := make(map[string]string)

	for _, e := range tree {
		if e.Mode.IsDir() {
			subTreeEntries, err := FlattenTree(repository, hex.EncodeToString(e.Sha))
			if err != nil {
				return nil, err
			}

			for k, v := range subTreeEntries {
				entries[path.Join(e.Path, k)] = v
			}

		} else {
			entries[e.Path] = hex.EncodeToString(e.Sha)
		}
	}

	return entries, nil
}
