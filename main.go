package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var REGEX_PATTERN_FOR_TODO string = `\/\/\s*TODO:\s*(.*)`

type TodoScanner struct {
}

func (scanner *TodoScanner) scanAllFiles(directoryPath string) error {
	archives, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, archive := range archives {
		if archive.IsDir() {
			scanner.scanAllFiles(archive.Name())
		}
		scanner.getAllTodosFromFile(archive.Name())
		fmt.Println(archive.Name())
	}

	return nil
}

func (scanner *TodoScanner) getAllTodosFromFile(fileName string) error {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line, err := getTodoFromLine(fileScanner.Text())
		if err != nil {
			return fmt.Errorf("Unable to read TODO from line")
		}
		fmt.Println(line)
	}

	return nil
}

func getTodoFromLine(line string) (*Todo, error) {
	regex := regexp.MustCompile(REGEX_PATTERN_FOR_TODO)
	matches := regex.FindSubmatch([]byte(line))
	for _, match := range matches {
		fmt.Println(string(match))
	}
	return nil, nil
}

func main() {
	var err error

	todoScanner := TodoScanner{}
	err = todoScanner.scanAllFiles(".")
	if err != nil {
		log.Fatal(err)
	}
}
