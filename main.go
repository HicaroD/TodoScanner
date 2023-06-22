package main

import "log"

func getRepositoryLinkFromCommandLine() (string, error) {
	return "", nil
}

func main() {
	var err error

	repositoryLink, err := getRepositoryLinkFromCommandLine()
	if err != nil {
		log.Fatal(err)
	}

	todoScanner := newScanner(repositoryLink)
	err = todoScanner.scanAllFiles(".")
	if err != nil {
		log.Fatal(err)
	}
}
