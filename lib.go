package gitgo

import (
	"fmt"
	"os"
)

var verbose = false

func setVerbose(v bool) {
	verbose = v
}

func verbosePrint(message string) {
	if verbose {
		fmt.Println(message)
	}
}

func createRepository(path string) error {
	var err error

	err = os.MkdirAll(path+"/.git/objects", os.ModePerm)
	verbosePrint("CREATE " + path + "/.git/objects")
	if err != nil {
		return err
	}
	err = os.MkdirAll(path+"/.git/refs/heads", os.ModePerm)
	verbosePrint("CREATE " + path + "/.git/refs/heads")
	if err != nil {
		return err
	}
	err = os.MkdirAll(path+"/.git/refs/tags", os.ModePerm)
	verbosePrint("CREATE " + path + "/.git/refs/tags")
	if err != nil {
		return err
	}
	head, err := os.Create(path + "/.git/HEAD")
	verbosePrint("CREATE " + path + "/.git/HEAD")
	if err != nil {
		return err
	}
	config, err := os.Create(path + "/.git/config")
	verbosePrint("CREATE " + path + "/.git/config")
	if err != nil {
		return err
	}
	description, err := os.Create(path + "/.git/description")
	verbosePrint("CREATE " + path + "/.git/description")
	if err != nil {
		return err
	}

	// .git/description
	verbosePrint("WRITE TO " + path + "/.git/description")
	_, err = fmt.Fprint(description, "Unnamed repository; edit this file 'description' to name the repository.\n")
	if err != nil {
		return err
	}

	// .git/HEAD
	verbosePrint("WRITE TO " + path + "/.git/HEAD")
	_, err = fmt.Fprint(head, "ref: refs/heads/master\n")
	if err != nil {
		return err
	}

	// .git/config
	verbosePrint("WRITE TO " + path + "/.git/config")
	err = createConfig(config)

	return err
}

func createConfig(file *os.File) error {
	var err error
	_, err = fmt.Fprint(file, "[core]")
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(file, "\trepositoryformatversion = 0")
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(file, "\tfilemode = false")
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(file, "\tbare = false")
	return err
}
