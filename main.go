package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var REGEX_PATTERN_FOR_TODO string = `TODO:\s*(.*)`

type TodoScanner struct {
	todos []Todo
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
		todo := scanner.getTodoFromLine(fileScanner.Text())
		if todo == nil {
			continue
		}
		scanner.todos = append(scanner.todos, *todo)
	}

	return nil
}

func (scanner *TodoScanner) getTodoFromLine(line string) *Todo {
	regex := regexp.MustCompile(REGEX_PATTERN_FOR_TODO)
	matches := regex.FindSubmatch([]byte(line))
	if len(matches) < 2 {
		return nil
	}
	for _, match := range matches {
		fmt.Println(string(match))
	}
	todoTitle := strings.TrimSpace(string(matches[1]))
	if len(todoTitle) == 0 {
		return nil
	}
	return newTodo(todoTitle)
}

func main() {
	var err error

	todoScanner := TodoScanner{}
	err = todoScanner.scanAllFiles(".")
	if err != nil {
		log.Fatal(err)
	}
}
