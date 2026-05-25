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

func (t Task) String() string {
	return fmt.Sprintf("[%d] %-30s %-12s created %s", t.ID, t.Description, t.Status, t.CreatedAt.Format("2006-01-02"))
}

func requireArgs(expected ...int) {
	if slices.Contains(expected, len(os.Args)) {
		return
	}
	fmt.Fprintln(os.Stderr, "error: incorrect number of arguments")
	os.Exit(1)
}

func findTaskById(tasks []Task, arg string) (Task, int) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)

	}

	index := slices.IndexFunc(tasks, func(task Task) bool {
		return task.ID == id
	})

	if index == -1 {
		fmt.Fprintf(os.Stderr, "error: task with ID %d not found\n", id)
		os.Exit(1)
	}

	return tasks[index], index
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: not enough arguments")
		os.Exit(1)
	}

	var tasks []Task
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error reading tasks.json: %v\n", err)
			os.Exit(1)
		}
	} else if len(data) > 0 {
		if err := json.Unmarshal(data, &tasks); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing tasks.json: %v\n", err)
			os.Exit(1)
		}
	}

	nextId := 1
	for _, task := range tasks {
		if task.ID >= nextId {
			nextId = task.ID + 1
		}
	}

	save := func() {
		encoding, werr := json.MarshalIndent(tasks, "", "  ")
		if werr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", werr)
			os.Exit(1)
		}

		if werr := os.WriteFile("tasks.json", encoding, 0644); werr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", werr)
			os.Exit(1)
		}
	}

	switch cmd := strings.ToLower(os.Args[1]); cmd {
	case "add":
		requireArgs(3) 
			task := Task{ID: nextId, Description: os.Args[2], Status: Todo, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			tasks = append(tasks, task)
			save()
			fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
		

	case "update":
		 requireArgs(4)
			task, index := findTaskById(tasks, os.Args[2])
			task.Description = os.Args[3]
			task.UpdatedAt = time.Now()
			tasks[index] = task
			save()
			fmt.Printf("Task (ID: %d) updated successfully\n", task.ID)
		

	case "delete":
		requireArgs(3) 
			task, index := findTaskById(tasks, os.Args[2])
			tasks = slices.Delete(tasks, index, index+1)
			save()
			fmt.Printf("Task (ID: %d) deleted successfully\n", task.ID)

	case "mark-in-progress":
		requireArgs(3) 
			task, index := findTaskById(tasks, os.Args[2])
			task.Status = InProgress
			task.UpdatedAt = time.Now()
			tasks[index] = task
			save()
			fmt.Printf("Task (ID: %d) marked as in progress\n", task.ID)

	case "mark-done":
		requireArgs(3) 
			task, index := findTaskById(tasks, os.Args[2])
			task.Status = Done
			task.UpdatedAt = time.Now()
			tasks[index] = task
			save()
			fmt.Printf("Task (ID: %d) marked as done\n", task.ID)

	case "list":
		requireArgs(2, 3) 
			switch len(os.Args) {
			case 2:
				for _, task := range tasks {
					fmt.Println(task)
				}
			case 3:
				status := strings.ToLower(os.Args[2])
				if !Status(status).IsValid() {
					fmt.Fprintf(os.Stderr, "invalid status: %v\n", status)
					os.Exit(1)
				}

				for _, task := range tasks {
					if task.Status == Status(status) {
						fmt.Println(task)
					}
				}
			}

	default:
		fmt.Fprintf(os.Stderr, "Incorrect command: %v\n", cmd)
	}

}
