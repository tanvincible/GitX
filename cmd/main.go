package main

import (
	"GitX/utils/file_operations"
	"GitX/utils/vcs_operations"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Staging area to hold files for the next commit
var stagingArea = make(map[string]string)

func main() {
    // Define flags
    initCommand := flag.NewFlagSet("init", flag.ExitOnError)

    addCommand := flag.NewFlagSet("add", flag.ExitOnError)
    addFilePath := addCommand.String("files", "", "Space-separated list of file paths to add")

    commitCommand := flag.NewFlagSet("commit", flag.ExitOnError)
    commitMessage := commitCommand.String("message", "", "Commit message")

    branchCommand := flag.NewFlagSet("branch", flag.ExitOnError)
    branchDelete := branchCommand.Bool("d", false, "Delete branch")
    branchName := branchCommand.String("name", "", "Branch name")

    checkoutCommand := flag.NewFlagSet("checkout", flag.ExitOnError)

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
            fmt.Printf("Error getting current working directory: %v\n", err)
            os.Exit(1)
        }
        repoPath := filepath.Join(cwd, repoName)
        file_operations.InitHandler(repoPath)
        
    case "add":
        addCommand.Parse(os.Args[2:])
        file_operations.AddHandler(*addFilePath, stagingArea)

    case "commit":
        commitCommand.Parse(os.Args[2:])
        // Commit all staged changes with the provided message
        file_operations.CommitHandler(*commitMessage, stagingArea)

    case "branch":
        branchCommand.Parse(os.Args[2:])
        if *branchDelete {
            // Delete the specified branch
            vcs_operations.DeleteBranch(*branchName)
        } else if *branchName != "" {
            // Create a new branch with the specified name
            vcs_operations.CreateBranch(*branchName)
        } else {
            // List all branches if no name is provided
            vcs_operations.ListBranches()
        }

    case "checkout":
        checkoutCommand.Parse(os.Args[2:])
        if len(checkoutCommand.Args()) != 1 {
            fmt.Println("Usage: gitx checkout <branch-name>")
            os.Exit(1)
        }
        branchName := checkoutCommand.Arg(0)
        vcs_operations.SwitchBranch(branchName)

    case "log":
        vcs_operations.LogHandler()

    case "status":
        file_operations.StatusHandler(stagingArea)
    
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
