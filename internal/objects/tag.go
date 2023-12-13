package objects

import (
	"bufio"
	"fmt"
	"mcavazotti/git-go/internal/repo"
	"os"
	"path"
	"strings"
)

type TagObject struct {
	object  string
	objType string
	tag     string
	tagger  string
	gpgsig  string
	message string
}

func TagToString(tag TagObject) string {
	var tagStr string

	tagStr += "object " + tag.object + "\n"
	tagStr += "type " + tag.objType + "\n"
	tagStr += "tag " + tag.tag + "\n"
	tagStr += "tagger " + tag.tagger + "\n"
	if tag.gpgsig != "" {
		tagStr += "gpgsig " + tag.gpgsig + "\n"
	}
	tagStr += "\n"
	tagStr += tag.message
	return tagStr
}

func WriteTag(repository *repo.Repository, tag TagObject) (string, error) {
	commitData := []byte(TagToString(tag))
	return WriteObject(repository, &commitData, "tag")
}

func ReadTag(repository *repo.Repository, sha string) (TagObject, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return TagObject{}, err
	}

	tagData := string(obj.Data)

	var tag TagObject

	scanner := bufio.NewScanner(strings.NewReader(tagData))
	readingMessage := false
	readingSignature := false
	for scanner.Scan() {
		if scanner.Text() == "" {
			readingMessage = true
			continue
		}
		if !readingMessage && !readingSignature {
			idx := strings.Index(scanner.Text(), " ")
			if idx != -1 {
				switch scanner.Text()[:idx] {
				case "object":
					tag.object = scanner.Text()[idx+1:]
				case "type":
					tag.objType = scanner.Text()[idx+1:]
				case "tag":
					tag.tag = scanner.Text()[idx+1:]
				case "tagger":
					tag.tagger = scanner.Text()[idx+1:]
				case "gpgsig":
					tag.gpgsig += scanner.Text()[idx+1:] + "\n"
					readingSignature = true
				}
			} else {
				readingMessage = true
			}
		} else {
			if readingSignature {
				tag.gpgsig += scanner.Text() + "\n"
				readingSignature = scanner.Text() != " -----END PGP SIGNATURE-----"
				readingMessage = scanner.Text() == " -----END PGP SIGNATURE-----"
			} else {
				tag.message += scanner.Text() + "\n"
			}
		}

	}
	if err = scanner.Err(); err != nil {
		return TagObject{}, err
	}
	return tag, nil
}

func CreateTag(repository *repo.Repository, name string, ref string, force bool, createObj bool, message string) error {
	tagsFolder := repository.RepoPath("refs/tags")
	files, err := os.ReadDir(tagsFolder)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.Name() == name && !force {
			return fmt.Errorf("tag '%s' already exists", name)
		}
	}

	sha, err := repository.Resolve(ref)
	if err != nil {
		return err
	}

	if !createObj {
		return os.WriteFile(path.Join(tagsFolder, name), []byte(sha), os.ModePerm)
	}

	gitObj, err := ReadObject(repository, sha)
	if err != nil {
		return fmt.Errorf("cannot update ref: 'refs/tags/%s': %s", name, err.Error())
	}

	tag := TagObject{
		object:  sha,
		objType: gitObj.ObjType,
		tag:     name,
		tagger:  "Todo",
		gpgsig:  "",
		message: message,
	}

	sha, err = WriteTag(repository, tag)

	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(tagsFolder, name), []byte(sha), os.ModePerm)
}
