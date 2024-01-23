package objects

import (
	"bufio"
	"mcavazotti/git-go/internal/repo"
	"strings"
)

type CommitObject struct {
	Tree      string
	Parent    []string
	Author    string
	Committer string
	Gpgsig    string
	Message   string
}

func CommitToString(commit CommitObject) string {
	var commitStr string

	commitStr += "tree " + commit.Tree + "\n"
	for _, p := range commit.Parent {
		commitStr += "parent " + p + "\n"
	}
	commitStr += "author " + commit.Author + "\n"
	commitStr += "committer " + commit.Committer + "\n"
	if commit.Gpgsig != "" {
		commitStr += "gpgsig " + commit.Gpgsig + "\n"
	}
	commitStr += "\n"
	commitStr += commit.Message
	return commitStr
}

func WriteCommit(repository *repo.Repository, commit CommitObject) error {
	commitData := []byte(CommitToString(commit))
	_, err := WriteObject(repository, &commitData, "commit")
	return err
}

func ReadCommit(repository *repo.Repository, sha string) (CommitObject, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return CommitObject{}, err
	}

	commitData := string(obj.Data)

	var commit CommitObject

	scanner := bufio.NewScanner(strings.NewReader(commitData))
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
				case "tree":
					commit.Tree = scanner.Text()[idx+1:]
				case "parent":
					commit.Parent = append(commit.Parent, scanner.Text()[idx+1:])
				case "author":
					commit.Author = scanner.Text()[idx+1:]
				case "committer":
					commit.Committer = scanner.Text()[idx+1:]
				case "gpgsig":
					commit.Gpgsig += scanner.Text()[idx+1:] + "\n"
					readingSignature = true
				}
			} else {
				readingMessage = true
			}
		} else {
			if readingSignature {
				commit.Gpgsig += scanner.Text() + "\n"
				readingSignature = scanner.Text() != " -----END PGP SIGNATURE-----"
				readingMessage = scanner.Text() == " -----END PGP SIGNATURE-----"
			} else {
				commit.Message += scanner.Text() + "\n"
			}
		}

	}
	if err = scanner.Err(); err != nil {
		return CommitObject{}, err
	}
	return commit, nil
}
