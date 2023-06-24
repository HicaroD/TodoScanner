package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var REGEX_PATTERN_FOR_TODO string = `TODO:\s*(.*)`

type TodoScanner struct {
	Github *GitHub
	Todos  []Todo
}

func newScanner(github *GitHub) *TodoScanner {
	return &TodoScanner{github, make([]Todo, 0)}
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
		fileExtension := filepath.Ext(archive.Name())
		if scanner.fileHasAllowedExtension(fileExtension) {
			scanner.getAllTodosFromFile(archive.Name())
		}
	}
	scanner.uploadTodos()
	return nil
}

func (scanner *TodoScanner) fileHasAllowedExtension(fileExtension string) bool {
	if len(scanner.Github.Extensions) == 0 {
		return true
	}

	for _, allowedExtension := range scanner.Github.Extensions {
		if fileExtension == allowedExtension {
			return true
		}
	}
	return false
}

func (scanner *TodoScanner) getAllTodosFromFile(fileName string) error {
	fileScanner, err := scanner.getLineByLineFileScanner(fileName)
	if err != nil {
		return err
	}

	for fileScanner.Scan() {
		todo := scanner.getTodoFromLine(fileScanner.Text())
		if todo == nil {
			continue
		}
		if scanner.userWantsToUploadTodo(todo) {
			scanner.Todos = append(scanner.Todos, *todo)
		}
	}

	return nil
}

func (*TodoScanner) getLineByLineFileScanner(fileName string) (*bufio.Scanner, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	return fileScanner, nil
}

func (scanner *TodoScanner) userWantsToUploadTodo(todo *Todo) bool {
	var answer string
	fmt.Printf("\nDo you want to upload the TODO below? [y/n]\n- %s\n", todo.Title)
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

func (scanner *TodoScanner) uploadTodos() error {
	for _, todo := range scanner.Todos {
		err := scanner.makeRequestInGitHubApi(todo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (scanner *TodoScanner) makeRequestInGitHubApi(todo Todo) error {
	client := &http.Client{}

	payload := scanner.getRequestPayload(todo)
	url := scanner.getGitHubUrl()

	request, err := scanner.getRequest(url, payload)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("unable to create issue.\nreason: %s", response.Body)
	}
	fmt.Println("Issue uploaded")
	return nil
}

func (scanner *TodoScanner) getRequest(url string, payload *bytes.Reader) (*http.Request, error) {
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", scanner.Github.GithubToken))
	return request, nil
}

func (scanner *TodoScanner) getGitHubUrl() string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues", scanner.Github.Repository)
	return url
}

func (*TodoScanner) getRequestPayload(todo Todo) *bytes.Reader {
	rawPayload := []byte(fmt.Sprintf(`{"title": "%s"}`, todo.Title))
	payload := bytes.NewReader(rawPayload)
	return payload
}
