package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)

func (s Status) IsValid() bool {
	switch s {
	case Todo, InProgress, Done:
		return true
	}
	return false

}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func checkNumberArgs(expected ...int) bool {
	if slices.Contains(expected, len(os.Args)) {
		return true
	}
	fmt.Fprintln(os.Stderr, "Error: incorrect number of arguments.")
	os.Exit(1)
	return false
}

func findTaskById(tasks []Task, arg string) (Task, int) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)

	}

	index := slices.IndexFunc(tasks, func(task Task) bool {
		return task.ID == id
	})

	return tasks[index], index
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: not enough arguments")
		os.Exit(1)
	}

	var tasks []Task
	data, err := os.ReadFile("tasks.json")
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading tasks: %v\n", err)
			os.Exit(1)
		}
	}

	counter := 1
	for _, task := range tasks {
		if task.ID >= counter {
			counter = task.ID + 1
		}
	}

	save := func() {
		data, err := json.MarshalIndent(tasks, "", " ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile("tasks.json", data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	switch cmd := strings.ToLower(os.Args[1]); cmd {
	case "add":
		if checkNumberArgs(3) {
			task := Task{ID: counter, Description: os.Args[2], Status: Todo, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			tasks = append(tasks, task)
			save()
			fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
		}

	case "update":
		if checkNumberArgs(4) {
			task, index := findTaskById(tasks, os.Args[2])
			task.Description = os.Args[3]
			tasks[index] = task
			save()
			fmt.Printf("Task (ID: %d) updated successfully\n", task.ID)
		}

	case "delete":
		if checkNumberArgs(3) {
			task, index := findTaskById(tasks, os.Args[2])
			tasks = slices.Delete(tasks, index, index+1)
			save()
			fmt.Printf("Task (ID: %d) deleted successfully\n", task.ID)
		}

	case "mark-in-progress":
		if checkNumberArgs(3) {
			task, index := findTaskById(tasks, os.Args[2])
			task.Status = InProgress
			tasks[index] = task
			save()
			fmt.Printf("Task (ID: %d) marked as in progress\n", task.ID)
		}

	case "mark-done":
		if checkNumberArgs(3) {
			task, index := findTaskById(tasks, os.Args[2])
			task.Status = Done
			tasks[index] = task
			save()
			fmt.Printf("Task (ID: %d) marked as done\n", task.ID)
		}

	case "list":
		if checkNumberArgs(2, 3) {
			switch len(os.Args) {
			case 2:
				for _, task := range tasks {
					fmt.Printf("%+v\n", task)
				}
			case 3:
				status := strings.ToLower(os.Args[2])
				if !Status(status).IsValid() {
					fmt.Fprintf(os.Stderr, "Invalid status: %v. Try 'done', 'in-progress' or 'todo'.\n", status)
					os.Exit(1)
				}

				for _, task := range tasks {
					if task.Status == Status(status) {
						fmt.Printf("%+v\n", task)
					}
				}
			}
		}

	default:
		fmt.Fprintln(os.Stderr, "Incorrect command. Try 'add', 'update', 'delete', 'mark-in-progress', 'mark-done' or 'list'.")
	}

}
