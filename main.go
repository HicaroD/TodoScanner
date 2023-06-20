package main

import (
	"fmt"
	"regexp"
)

func main() {
	regex := regexp.MustCompile(`\/\/\s*TODO:\s*(.*)`)
	matches := regex.FindSubmatch([]byte("// TODO: my TODO text"))
	if len(matches) != 2 {
		fmt.Println("No match")
		return
	}
	fmt.Println(string(matches[1]))
}
