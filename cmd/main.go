package main

import (
	"GitX/utils/file_operations"
	"GitX/utils/vcs_operations"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Staging area to hold files for the next commit
var stagingArea = make(map[string]string)

func main() {
	// Define flags
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)

	commitCommand := flag.NewFlagSet("commit", flag.ExitOnError)
	commitMessage := commitCommand.String("message", "", "Commit message")

	branchCommand := flag.NewFlagSet("branch", flag.ExitOnError)
	branchDelete := branchCommand.Bool("d", false, "Delete branch")

	checkoutCommand := flag.NewFlagSet("checkout", flag.ExitOnError)
	checkoutBranch := checkoutCommand.String("b", "", "Switch to branch")

	configCommand := flag.NewFlagSet("config", flag.ExitOnError)

	// Parse command-line arguments
	flag.Parse()

	// Ensure a command is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: gitx <command> [options]")
		os.Exit(1)
	}

	// Execute the appropriate command
	switch os.Args[1] {
	case "init":
		initCommand.Parse(os.Args[2:])
		if len(initCommand.Args()) != 1 {
			fmt.Println("Usage: gitx init <repo-name>")
			os.Exit(1)
		}
		repoName := initCommand.Arg(0)
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %v\n", err)
		}
		repoPath := filepath.Join(cwd, repoName)
		file_operations.InitHandler(repoPath)

	case "config":
		// Handle config command
		configCommand.Parse(os.Args[2:])
		if len(configCommand.Args()) != 2 {
			fmt.Println("Usage: gitx config <key> <value>")
			os.Exit(1)
		}
		configKey := configCommand.Arg(0)
		configValue := configCommand.Arg(1)

		// Ensure the command looks for the correct config file
		configFilePath := filepath.Join(".gitx", "config.toml")
		file_operations.ConfigHandlerWithFilePath(configFilePath, configKey, configValue)

	case "add":
		// Handle add command
		if len(os.Args) < 3 {
			fmt.Println("Error: No file path provided for the 'add' command")
			os.Exit(1)
		}

		// Loop through all the provided file paths
		for _, filePath := range os.Args[2:] {
			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				fmt.Printf("Error getting absolute path for file '%s': %v\n", filePath, err)
				os.Exit(1)
			}
			err = file_operations.AddHandler(filepath.Join(".gitx", "INDEX"), absFilePath)
			if err != nil {
				fmt.Printf("Error adding file '%s': %v\n", filePath, err)
				os.Exit(1)
			}
		}

	case "commit":
		commitCommand.Parse(os.Args[2:])
		if *commitMessage == "" {
			fmt.Println("Error: Commit message is required for the 'commit' command")
			os.Exit(1)
		}
		file_operations.CommitHandler(*commitMessage)

	case "branch":
		// Parse the command line arguments starting from the second argument
		branchCommand.Parse(os.Args[2:])

		// Check if the delete flag is set
		if *branchDelete {
			// Delete the specified branch
			if branchCommand.NArg() == 0 {
				fmt.Println("Error: branch name is required")
				os.Exit(1)
			}
			branchName := branchCommand.Arg(0)
			err := vcs_operations.DeleteBranch(branchName)
			if err != nil {
				fmt.Println("Error deleting branch:", err)
				os.Exit(1)
			}
		} else {
			// If the delete flag is not set, check for the branch name as a positional argument
			if branchCommand.NArg() > 0 {
				branchName := branchCommand.Arg(0)
				// Create a new branch with the specified name
				err := vcs_operations.CreateBranch(branchName)
				if err != nil {
					fmt.Println("Error creating branch:", err)
					os.Exit(1)
				}
			} else {
				// List all branches if no additional arguments are provided
				vcs_operations.ListBranches()
			}
		}

	case "checkout":
		checkoutCommand.Parse(os.Args[2:])
		if *checkoutBranch != "" {
			branchName := *checkoutBranch
			// Try to create the branch if it does not exist
			err := vcs_operations.CreateBranch(branchName)
			if err != nil && !strings.Contains(err.Error(), "already exists") {
				fmt.Println("Error creating branch:", err)
				os.Exit(1)
			}
			// Switch to the branch
			err = vcs_operations.SwitchBranch(branchName)
			if err != nil {
				fmt.Println("Error switching to branch:", err)
				os.Exit(1)
			}
		} else if len(checkoutCommand.Args()) != 1 {
			fmt.Println("Usage: gitx checkout [-b] <branch-name>")
			os.Exit(1)
		} else {
			branchName := checkoutCommand.Arg(0)
			err := vcs_operations.SwitchBranch(branchName)
			if err != nil {
				fmt.Println("Error switching to branch:", err)
				os.Exit(1)
			}
		}

	case "log":
		vcs_operations.LogHandler()

	case "status":
		file_operations.StatusHandler()

	case "merge":
		// Define flags for merge command
		mergeCommand := flag.NewFlagSet("merge", flag.ExitOnError)
		mergeBranchName := mergeCommand.String("branch", "", "Branch name to merge")

		// Parse flags for merge command
		mergeCommand.Parse(os.Args[2:])
		if len(os.Args) != 4 {
			fmt.Println("Usage: gitx merge -branch <branch-name>")
			os.Exit(1)
		}
		// Call MergeBranch function from the vcs_operations package
		if err := vcs_operations.MergeBranch(*mergeBranchName); err != nil {
			fmt.Printf("Error merging branch: %v\n", err)
			os.Exit(1)
		}

	case "squash":
		// Define flags for squash command
		squashCommand := flag.NewFlagSet("squash", flag.ExitOnError)
		baseCommit := squashCommand.String("base-commit", "", "Base commit")
		targetCommit := squashCommand.String("target-commit", "", "Target commit")

		// Parse flags for squash command
		squashCommand.Parse(os.Args[2:])
		if len(os.Args) != 6 {
			fmt.Println("Usage: gitx squash -base-commit <base-commit> -target-commit <target-commit>")
			os.Exit(1)
		}
		// fmt.Printf("Squashing commits from %s to %s\n", *baseCommit, *targetCommit) // Debug print
		if err := vcs_operations.SquashCommits(*baseCommit, *targetCommit); err != nil {
			fmt.Printf("Error squashing commits: %v\n", err)
			os.Exit(1)
		}

	case "stash":
		if err := vcs_operations.Stash(); err != nil {
			fmt.Printf("Error stashing changes: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Changes stashed successfully")

	case "cat-file":
		// Define flags for cat-file command
		catFileCommand := flag.NewFlagSet("cat-file", flag.ExitOnError)
		objectType := catFileCommand.String("type", "", "Object type")

		// Parse flags for cat-file command
		catFileCommand.Parse(os.Args[2:])
		if len(os.Args) != 4 {
			fmt.Println("Usage: gitx cat-file -type <object-type>")
			os.Exit(1)
		}
		// Call the CatFile function from the vcs_operations package
		if err := vcs_operations.CatFile(*objectType); err != nil {
			fmt.Printf("Error executing cat-file: %v\n", err)
			os.Exit(1)
		}

	case "reflog":
		// Call ReflogHandler from the vcs_operations package
		vcs_operations.ReflogHandler()

	default:
		fmt.Printf("gitx: %s is not a valid command\n", os.Args[1])
		os.Exit(1)
	}
}
