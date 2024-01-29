package objects

import (
	"bufio"
	"mcavazotti/git-go/internal/repo"
	"mcavazotti/git-go/internal/shared"
	"os"
	"path"
	"path/filepath"
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
	shared.VerbosePrintln("ReadGitIgnore")
	var ignore GitIgnore

	excludeFile := repository.RepoPath(path.Join("info", "exclude"))
	shared.VerbosePrintln("Parsing>", excludeFile)

	if _, err := os.Stat(excludeFile); err == nil {
		rules, err := parseIgnoreFile(excludeFile)
		if err != nil {
			return ignore, err
		}

		ignore.absolute = append(ignore.absolute, rules...)
	}

	var configHome string
	shared.VerbosePrintln("Finding config home...")
	shared.VerbosePrintln("Looking for ENV variable...")
	if val, exists := os.LookupEnv("XDG_CONFIG_HOME"); exists && val != "" {
		shared.VerbosePrintln("Found ENV variable")
		configHome = val
	} else {
		shared.VerbosePrintln("ENV variable not found")
		shared.VerbosePrintln("Looking for user home")
		home, err := os.UserHomeDir()
		if err != nil {
			return ignore, err
		}
		configHome = path.Join(home, ".config")
	}
	shared.VerbosePrintln("Found config home")

	globalIgnore := path.Join(configHome, "git", "ignore")

	shared.VerbosePrintln("Exists?>", globalIgnore)
	if _, err := os.Stat(globalIgnore); err == nil {
		rules, err := parseIgnoreFile(globalIgnore)
		if err != nil {
			return ignore, err
		}

		ignore.absolute = append(ignore.absolute, rules...)
	}

	shared.VerbosePrintln("Reading index...")
	index, err := ReadIndex(repository)
	shared.VerbosePrintln("Read index")

	if err != nil {
		return ignore, err
	}

	shared.VerbosePrintln("Looking for repo's ignore files")
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
	shared.VerbosePrintln("parseIgnoreFile>", p)
	file, err := os.Open(p)
	if err != nil {
		return []GitIgnoreRule{}, err
	}

	scanner := bufio.NewScanner(file)

	var rules []GitIgnoreRule
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\n")
		shared.VerbosePrintln("line", line)
		if line == "" || line[0] == '#' {
			shared.VerbosePrintln("ignore line")
			continue
		} else if line[1] == '!' {
			rules = append(rules, GitIgnoreRule{rule: filepath.FromSlash(path.Join(p[:len(p)-10], line[1:])), exclude: false})
		} else if line[1] == '\\' {
			rules = append(rules, GitIgnoreRule{rule: filepath.FromSlash(path.Join(p[:len(p)-10], line[1:])), exclude: true})
		} else {
			rules = append(rules, GitIgnoreRule{rule: filepath.FromSlash(path.Join(p[:len(p)-10], line)), exclude: true})
		}
	}
	shared.VerbosePrintln("Num rules", len(rules))
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
