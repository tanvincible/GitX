package main

import (
	"GitX/utils/file_operations"
	"GitX/utils/vcs_operations"
	"fmt"
	"os"
)

// Staging area to hold files for next commit
var stagingArea = make(map[string]string)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: gitx <command> [init, add, commit, status, branch, merge, squash, stash]")
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
	case "branch":
		if len(args) != 2 {
			fmt.Println("Usage: gitx branch <create/list/switch/delete>")
			os.Exit(1)
		}
		// Handle branch-related logic using functions from vcs_operations
		switch args[1] {
		case "create":
			// Call CreateBranch function
			vcs_operations.CreateBranch(args[2])
		case "list":
			// Call ListBranches function
			vcs_operations.ListBranches()
		case "switch":
			// Call SwitchBranch function
			vcs_operations.SwitchBranch(args[2])
		case "delete":
			// Call DeleteBranch function
			vcs_operations.DeleteBranch(args[2])
		default:
			fmt.Printf("gitx branch: %s is not a valid subcommand\n", args[1])
			os.Exit(1)
		}
	case "merge":
		if len(args) != 2 {
			fmt.Println("Usage: gitx merge <branch-name>")
			os.Exit(1)
		}
		// Call MergeBranch function
		if err := vcs_operations.MergeBranch(args[1]); err != nil {
			fmt.Printf("Error merging branch: %v\n", err)
			os.Exit(1)
		}
	case "squash":
		if len(args) != 3 {
			fmt.Println("Usage: gitx squash <base-commit> <target-commit>")
			os.Exit(1)
		}
		fmt.Printf("Squashing commits from %s to %s\n", args[1], args[2]) // Debug print
		if err := vcs_operations.SquashCommits(args[1], args[2]); err != nil {
			fmt.Printf("Error squashing commits: %v\n", err)
			os.Exit(1)
		}
	case "stash":
		if err := vcs_operations.Stash(); err != nil {
			fmt.Printf("Error stashing changes: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Changes stashed successfully")
	default:
		fmt.Printf("gitx: %s is not a valid command\n", command)
		os.Exit(1)
	}
}
