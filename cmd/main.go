package main

import (
	"GitX/utils/file_operations"
	"fmt"
	"os"
)

// Staging area to hold files for next commit
var stagingArea = make(map[string]string)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: gitx <command> [init, add, commit, status]")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "init":
		if len(args) != 2 {
			fmt.Println("Usage: gitx init <directory>")
			os.Exit(1)
		}
		// Call InitHandler from the utils package
		file_operations.InitHandler(args[1])
	case "add":
		if len(args) != 2 {
			fmt.Println("Usage: gitx add <file-path>")
			os.Exit(1)
		}
		file_operations.AddHandler(args[1], stagingArea)
	case "commit":
		if len(args) != 3 {
			fmt.Println("Usage: gitx commit <file-path> <message>")
			os.Exit(1)
		}
		// Call CommitHandler from the utils package
		file_operations.CommitHandler(args[1], stagingArea)
	case "status":
		// Call StatusHandler from the utils package
		file_operations.StatusHandler(stagingArea)
	default:
		fmt.Printf("gitx: %s is not a valid command\n", command)
		os.Exit(1)
	}
}
