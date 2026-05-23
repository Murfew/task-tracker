package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func checkNumberArgs(expected ...int) bool {
	if slices.Contains(expected, len(os.Args)) {
		return true
	}
	fmt.Fprintln(os.Stderr, "Error: incorrect number of arguments.")
	os.Exit(1)
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: not enough arguments")
		os.Exit(1)
	}

	switch cmd := strings.ToLower(os.Args[1]); cmd {
	case "add":
		if checkNumberArgs(3) {
			//TODO
		}

	case "update":
		if checkNumberArgs(4) {
			//TODO
		}

	case "delete":
		if checkNumberArgs(3) {
			//TODO
		}

	case "mark-in-progress":
		if checkNumberArgs(3) {
			//TODO
		}

	case "mark-done":
		if checkNumberArgs(3) {
			//TODO
		}

	case "list":
		if checkNumberArgs(2, 3) {
			//TODO
		}

	default:
		fmt.Fprintln(os.Stderr, "Incorrect command. Try 'add', 'update', 'delete', 'mark-in-progress', 'mark-done' or 'list'.")
	}

	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
}
