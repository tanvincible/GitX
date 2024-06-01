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

	fmt.Println("Arguments:", args) // Print out the arguments for debugging

	if len(args) < 2 || args[0] != "gitx" {
		fmt.Println("Usage: gitx <command> [init, add, commit, status, branch, merge, squash, stash, push, pull, clone]")
		os.Exit(1)
	}

	command := args[1]
	switch command {
	case "init":
		if len(args) != 3 {
			fmt.Println("Usage: gitx init <directory>")
			os.Exit(1)
		}
		// Call InitHandler from the file_operations package
		file_operations.InitHandler(args[2])
	case "add":
		if len(args) != 3 {
			fmt.Println("Usage: gitx add <file-path>")
			os.Exit(1)
		}
		file_operations.AddHandler(args[2], stagingArea)
	case "commit":
		if len(args) != 4 {
			fmt.Println("Usage: gitx commit <file-path> <message>")
			os.Exit(1)
		}
		// Call CommitHandler from the file_operations package
		file_operations.CommitHandler(args[2], stagingArea)
	case "status":
		// Call StatusHandler from the file_operations package
		file_operations.StatusHandler(stagingArea)
	case "branch":
		if len(args) < 3 {
			fmt.Println("Usage: gitx branch <create/list/switch/delete>")
			os.Exit(1)
		}
		// Handle branch-related logic using functions from vcs_operations
		branchSubCommand := args[2]
		switch branchSubCommand {
		case "create":
			if len(args) != 4 {
				fmt.Println("Usage: gitx branch create <branch-name>")
				os.Exit(1)
			}
			vcs_operations.CreateBranch(args[3])
		case "list":
			vcs_operations.ListBranches()
		case "switch":
			if len(args) != 4 {
				fmt.Println("Usage: gitx branch switch <branch-name>")
				os.Exit(1)
			}
			vcs_operations.SwitchBranch(args[3])
		case "delete":
			if len(args) != 4 {
				fmt.Println("Usage: gitx branch delete <branch-name>")
				os.Exit(1)
			}
			vcs_operations.DeleteBranch(args[3])
		default:
			fmt.Printf("gitx branch: %s is not a valid subcommand\n", branchSubCommand)
			os.Exit(1)
		}
	case "merge":
		if len(args) != 3 {
			fmt.Println("Usage: gitx merge <branch-name>")
			os.Exit(1)
		}
		// Call MergeBranch function from the vcs_operations package
		if err := vcs_operations.MergeBranch(args[2]); err != nil {
			fmt.Printf("Error merging branch: %v\n", err)
			os.Exit(1)
		}
	case "squash":
		if len(args) != 4 {
			fmt.Println("Usage: gitx squash <base-commit> <target-commit>")
			os.Exit(1)
		}
		// fmt.Printf("Squashing commits from %s to %s\n", args[2], args[3]) // Debug print
		if err := vcs_operations.SquashCommits(args[2], args[3]); err != nil {
			fmt.Printf("Error squashing commits: %v\n", err)
			os.Exit(1)
		}
	case "stash":
		if err := vcs_operations.Stash(); err != nil {
			fmt.Printf("Error stashing changes: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Changes stashed successfully")
	case "log":
		// Call LogHandler from the vcs_operations package
		vcs_operations.LogHandler()
	case "cat-file":
		// Call the CatFile function from the vcs_operations package
		if err := vcs_operations.CatFile(args[1]); err != nil {
			fmt.Printf("Error executing cat-file: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("gitx: %s is not a valid command\n", command)
		os.Exit(1)
	}
}
