package objects

import (
	"bufio"
	"mcavazotti/git-go/internal/repo"
	"os"
	"os/user"
	"path"
	"strings"
)

type GitIgnoreRule struct {
	rule    string
	exclude bool
}

type GitIgnore struct {
	scoped   []GitIgnoreRule
	absolute []GitIgnoreRule
}

func ReadGitIgnore(repository *repo.Repository) (GitIgnore, error) {
	var ignore GitIgnore

	excludeFile := repository.RepoPath(path.Join("info", "exclude"))

	if _, err := os.Stat(excludeFile); err == nil {
		rules, err := parseIgnoreFile(excludeFile)
		if err != nil {
			return ignore, err
		}

		ignore.absolute = append(ignore.absolute, rules...)
	}

	var configHome string
	if val, exists := os.LookupEnv("XDG_CONFIG_HOME"); exists && val != "" {
		configHome = val
	} else {
		usr, err := user.Current()
		if err != nil {
			return ignore, err
		}
		configHome = path.Join(usr.HomeDir, ".config")
	}

	globalIgnore := path.Join(configHome, "git", "ignore")

	if _, err := os.Stat(globalIgnore); err == nil {
		rules, err := parseIgnoreFile(globalIgnore)
		if err != nil {
			return ignore, err
		}

		ignore.absolute = append(ignore.absolute, rules...)
	}

	index, err := ReadIndex(repository)

	if err != nil {
		return ignore, err
	}

	for _, entry := range index.Entries {
		if strings.HasSuffix(entry.Name, ".gitignore") {
			rules, err := parseIgnoreFile(path.Join(repository.WorkTree, entry.Name))
			if err != nil {
				return ignore, err
			}
			ignore.scoped = append(ignore.scoped, rules...)
		}
	}
	return ignore, nil
}

func parseIgnoreFile(p string) ([]GitIgnoreRule, error) {
	file, err := os.Open(p)
	if err != nil {
		return []GitIgnoreRule{}, err
	}

	scanner := bufio.NewScanner(file)

	var rules []GitIgnoreRule
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n")
		if line == "" || line[1] == '#' {
			continue
		} else if line[1] == '!' {
			rules = append(rules, GitIgnoreRule{rule: path.Join(p[:len(p)-10], line[1:]), exclude: false})
		} else if line[1] == '\\' {
			rules = append(rules, GitIgnoreRule{rule: path.Join(p[:len(p)-10], line[1:]), exclude: true})
		} else {
			rules = append(rules, GitIgnoreRule{rule: path.Join(p[:len(p)-10], line), exclude: true})
		}
	}
	return rules, nil
}

func (g GitIgnore) IgnoreFile(p string) bool {
	for _, r := range g.scoped {
		if strings.Contains(p, r.rule) {
			return r.exclude
		}
	}

	for _, r := range g.absolute {
		if strings.Contains(p, r.rule) {
			return r.exclude
		}
	}

	return false
}
