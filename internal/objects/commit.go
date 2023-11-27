package objects

import (
	"bufio"
	"mcavazotti/git-go/internal/repo"
	"strings"
)

type CommitObject struct {
	tree      string
	parent    []string
	author    string
	committer string
	gpgsig    string
	message   string
}

func CommitToString(commit CommitObject) string {
	var commitStr string

	commitStr += "tree " + commit.tree + "\n"
	for _, p := range commit.parent {
		commitStr += "parent " + p + "\n"
	}
	commitStr += "author " + commit.author + "\n"
	commitStr += "committer " + commit.committer + "\n"
	if commit.gpgsig != "" {
		commitStr += "gpgsig " + commit.gpgsig + "\n"
	}
	commitStr += "\n"
	commitStr += commit.message
	return commitStr
}

func WriteCommit(repository *repo.Repository, commit CommitObject) error {
	commitData := []byte(CommitToString(commit))
	return WriteObject(repository, &commitData, "commit")
}

func ReadCommit(repository *repo.Repository, sha string) (CommitObject, error) {
	obj, err := ReadObject(repository, sha)
	if err != nil {
		return CommitObject{}, err
	}

	commitData := string(obj.data)

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
					commit.tree = scanner.Text()[idx+1:]
				case "parent":
					commit.parent = append(commit.parent, scanner.Text()[idx+1:])
				case "author":
					commit.author = scanner.Text()[idx+1:]
				case "committer":
					commit.committer = scanner.Text()[idx+1:]
				case "gpgsig":
					commit.gpgsig += scanner.Text()[idx+1:] + "\n"
					readingSignature = true
				}
			} else {
				readingMessage = true
			}
		} else {
			if readingSignature {
				commit.gpgsig += scanner.Text() + "\n"
				readingSignature = scanner.Text() != " -----END PGP SIGNATURE-----"
				readingMessage = scanner.Text() == " -----END PGP SIGNATURE-----"
			} else {
				commit.message += scanner.Text() + "\n"
			}
		}

	}
	if err = scanner.Err(); err != nil {
		return CommitObject{}, err
	}
	return commit, nil
}
