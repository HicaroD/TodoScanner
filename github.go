package main

import (
	"encoding/json"
	"os"
)

type GitHub struct {
	Repository  string   `json:"repo"`
	GithubToken string   `json:"token"`
	Extensions  []string `json:"extensions"`
}

func newGitHub() (*GitHub, error) {
	file, err := os.ReadFile("todo.json")
	if err != nil {
		return nil, err
	}

	var github GitHub
	err = json.Unmarshal(file, &github)
	if err != nil {
		return nil, err
	}

	return &github, nil
}
