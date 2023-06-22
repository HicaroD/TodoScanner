package main

type Todo struct {
	Title string
}

func newTodo(title string) *Todo {
	return &Todo{title}
}
