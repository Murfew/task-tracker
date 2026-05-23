package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer file.Close()
}
