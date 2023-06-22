package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var REGEX_PATTERN_FOR_TODO string = `TODO:\s*(.*)`

type TodoScanner struct {
	repositoryLink string
	todos          []Todo
}

func newScanner(repositoryLink string) *TodoScanner {
	return &TodoScanner{repositoryLink, make([]Todo, 0)}
}

func (scanner *TodoScanner) scanAllFiles(directoryPath string) error {
	archives, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, archive := range archives {
		if archive.IsDir() {
			scanner.scanAllFiles(archive.Name())
		}
		scanner.getAllTodosFromFile(archive.Name())
	}
	scanner.uploadTodos()
	return nil
}

func (scanner *TodoScanner) getAllTodosFromFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		todo := scanner.getTodoFromLine(fileScanner.Text())
		if todo == nil {
			continue
		}
		if scanner.userWantsToUploadTodo(todo) {
			scanner.todos = append(scanner.todos, *todo)
		}
	}

	return nil
}

func (scanner *TodoScanner) userWantsToUploadTodo(todo *Todo) bool {
	var answer string
	fmt.Printf("\nDo you want to upload the TODO below? [y/n]\n- %s\n", todo.title)
	fmt.Scanln(&answer)

	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)
	return answer == "y"
}

func (scanner *TodoScanner) getTodoFromLine(line string) *Todo {
	regex := regexp.MustCompile(REGEX_PATTERN_FOR_TODO)
	matches := regex.FindSubmatch([]byte(line))
	if len(matches) < 2 {
		return nil
	}
	todoTitle := strings.TrimSpace(string(matches[1]))
	if len(todoTitle) == 0 {
		return nil
	}
	return newTodo(todoTitle)
}

func (scanner *TodoScanner) uploadTodos() {

}
