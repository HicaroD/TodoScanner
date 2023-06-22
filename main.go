package main

import "log"

func main() {
	var err error
	github, err := newGitHub()
	if err != nil {
		log.Fatal(err)
	}
	todoScanner := newScanner(github)

	err = todoScanner.scanAllFiles(".")
	if err != nil {
		log.Fatal(err)
	}
}
