package main

type Todo struct {
	title string
}

func newTodo(title string) *Todo {
	return &Todo{title}
}
