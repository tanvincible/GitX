package file_operations

import (
	"GitX/internal/hash"
	"GitX/models"
	"GitX/utils/metadata_operations"
	"GitX/utils/vcs_operations"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/BurntSushi/toml"
)

// InitHandler initializes a new GitX repository by creating the necessary directories and files
func InitHandler(directory string) {
	// Create repository directory
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatalf("Error creating repository directory: %v", err)
	}

	// Create .gitx directory inside the repository directory
	gitxDir := filepath.Join(directory, ".gitx")
	if err := os.MkdirAll(gitxDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating .gitx directory: %v", err)
	}

	// Create metadata file
	metadataFile := filepath.Join(gitxDir, "metadata.json")
	if _, err := os.Create(metadataFile); err != nil {
		log.Fatalf("Error creating metadata file: %v", err)
	}

	// Create HEAD file
	headFile := filepath.Join(gitxDir, "HEAD")
	if err := os.WriteFile(headFile, []byte("refs/heads/main"), 0644); err != nil {
		log.Fatalf("Error creating HEAD file: %v", err)
	}

	// Create refs/heads directory
	refsHeadsDir := filepath.Join(gitxDir, "refs", "heads")
	if err := os.MkdirAll(refsHeadsDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating refs/heads directory: %v", err)
	}

	// Create main branch file
	mainBranchFile := filepath.Join(refsHeadsDir, "main")
	if _, err := os.Create(mainBranchFile); err != nil {
		log.Fatalf("Error creating main branch file: %v", err)
	}

	// Create commits directory
	commitsDir := filepath.Join(gitxDir, "commits")
	if err := os.MkdirAll(commitsDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating commits directory: %v", err)
	}

	// Create objects directory
	objectsDir := filepath.Join(gitxDir, "objects")
	if err := os.MkdirAll(objectsDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating objects directory: %v", err)
	}

	// Set up ignore file
	ignoreFile := filepath.Join(directory, ".gitxignore")
	if _, err := os.Create(ignoreFile); err != nil {
		log.Fatalf("Error creating ignore file: %v", err)
	}

	// Create config file with default contents in TOML format
	configFile := filepath.Join(gitxDir, "config.toml")
	config := models.GitXConfig{
		UserName:  "Your Name",
		UserEmail: "your.email@example.com",
	}
	err := UpdateConfig(configFile, &config)
	if err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}

	// Create description file
	descriptionFile := filepath.Join(gitxDir, "description")
	descriptionContent := []byte("Unnamed repository; edit this file to name the repository.\n")
	if err := os.WriteFile(descriptionFile, descriptionContent, 0644); err != nil {
		log.Fatalf("Error creating description file: %v", err)
	}

	// Create INDEX file
	indexFile := filepath.Join(gitxDir, "INDEX")
	if _, err := os.Create(indexFile); err != nil {
		log.Fatalf("Error creating INDEX file: %v", err)
	}

	fmt.Printf("Initialized empty repository in %s\n", directory)
	fmt.Println("Please configure your user information using the following commands:")
	fmt.Println("  gitx config user.name 'Your Name'")
	fmt.Println("  gitx config user.email 'your.email@example.com'")
}

// ConfigHandler reads and updates configuration settings.
func ConfigHandler(key, value string) {
	configFile := filepath.Join(".gitx", "config.toml")

	// Load existing config
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Update config based on the key
	switch key {
	case "user.name":
		config.UserName = value
	case "user.email":
		config.UserEmail = value
	default:
		log.Fatalf("Unknown config key: %s", key)
	}

	// Write updated config back to file
	err = UpdateConfig(configFile, config)
	if err != nil {
		log.Fatalf("Error updating config: %v", err)
	}

	fmt.Printf("Config updated successfully!\n")
}

// LoadConfig reads the configuration from a file.
func LoadConfig(filePath string) (*models.GitXConfig, error) {
	config := &models.GitXConfig{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If the config file does not exist, return an empty config with no error
		return config, nil
	}

	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return nil, err
	}
	return config, nil
}

// UpdateConfig writes the updated configuration back to the file.
func UpdateConfig(filePath string, config *models.GitXConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

// ConfigHandlerWithFilePath reads and updates configuration settings from the specified config file path.
func ConfigHandlerWithFilePath(configFilePath, key, value string) {
	// Load existing config
	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Update config based on the key
	switch key {
	case "user.name":
		config.UserName = value
	case "user.email":
		config.UserEmail = value
	default:
		log.Fatalf("Unknown config key: %s", key)
	}

	// Write updated config back to file
	err = UpdateConfig(configFilePath, config)
	if err != nil {
		log.Fatalf("Error updating config: %v", err)
	}

	fmt.Printf("Config updated successfully!\n")
}

// AddHandler adds a file to the index for staging, following Git conventions.
func AddHandler(indexFilePath, absFilePath string) error {
	// Calculate the SHA-1 hash of the file
	hashValue, err := hash.SHA1Hash(absFilePath)
	if err != nil {
		return fmt.Errorf("error calculating hash for file %s: %w", absFilePath, err)
	}

	// Normalize the file path to use forward slashes
	normalizedPath := filepath.ToSlash(absFilePath)

	// Open the INDEX file for appending or create it if not exists
	indexFile, err := os.OpenFile(indexFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening INDEX file: %w", err)
	}
	defer indexFile.Close()

	// Format the entry to write into the INDEX file
	entry := fmt.Sprintf("100644 %s %d\t%s\n", hashValue, 0, normalizedPath)

	// Write the entry into the INDEX file
	if _, err := indexFile.WriteString(entry); err != nil {
		return fmt.Errorf("error writing to INDEX file: %w", err)
	}

	return nil
}

// UpdateIndex updates the index file with the provided file path.
func UpdateIndex(indexFile, filePath string) error {
	// Open the index file for appending
	f, err := os.OpenFile(indexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open index file: %w", err)
	}
	defer f.Close()

	// Write the file path to the index file
	_, err = fmt.Fprintf(f, "%s\n", filePath)
	if err != nil {
		return fmt.Errorf("failed to write to index file: %w", err)
	}

	return nil
}

// CommitHandler creates a commit object, compresses the file content, stores the compressed file,
// updates metadata, and updates the HEAD reference.
func CommitHandler(message string) {
	// Create the commits directory if it doesn't exist
	commitsDir := ".gitx/commits"
	if _, err := os.Stat(commitsDir); os.IsNotExist(err) {
		os.MkdirAll(commitsDir, os.ModePerm)
	}

	headFile := ".gitx/HEAD"
	headContent, err := os.ReadFile(headFile)
	if err != nil {
		log.Fatalf("Error reading HEAD file: %v", err)
	}

	headBranch := strings.TrimSpace(strings.TrimPrefix(string(headContent), "refs/heads/"))
	if headBranch == string(headContent) {
		headBranch = strings.TrimSpace(string(headContent))
	}

	branchRefPath := filepath.Join(".gitx", "refs", "heads", headBranch)
	branchRefPath = filepath.Clean(branchRefPath) // Ensure the path is clean

	parentCommitHash, err := os.ReadFile(branchRefPath)
	if err != nil {
		log.Fatalf("Error reading branch ref file: %v", err)
	}

	var parentCommit *models.Commit
	if len(parentCommitHash) > 0 {
		parentCommit, err = vcs_operations.GetCommitByHash(strings.TrimSpace(string(parentCommitHash)))
		if err != nil {
			log.Fatalf("Error retrieving parent commit: %v", err)
		}
	} else {
		// If there is no parent commit, create an initial commit for the main branch
		initialCommit := createInitialCommit()

		// Write the initial commit to the commits directory
		initialCommitData, err := json.Marshal(initialCommit)
		if err != nil {
			log.Fatalf("Error serializing initial commit data: %v", err)
		}
		initialCommitFilePath := filepath.Join(commitsDir, initialCommit.ID)
		if err := os.WriteFile(initialCommitFilePath, initialCommitData, 0644); err != nil {
			log.Fatalf("Error writing initial commit file: %v", err)
		}

		if err := vcs_operations.UpdateHEAD("refs/heads/main"); err != nil {
			log.Fatalf("Error updating HEAD with main branch reference: %v", err)
		}

		mainBranchRefPath := filepath.Join(".gitx", "refs", "heads", "main")
		if err := os.WriteFile(mainBranchRefPath, []byte(initialCommit.ID), 0644); err != nil {
			log.Fatalf("Error creating main branch ref file: %v", err)
		}

		parentCommit = &initialCommit
	}

	tree, err := vcs_operations.CreateTreeFromIndex(".gitx/index")
	if err != nil {
		log.Fatalf("Error creating tree from INDEX: %v", err)
	}

	newCommit := models.Commit{
		ID:        "",
		Parent:    []*models.Commit{},
		Tree:      tree,
		Message:   message,
		Author:    vcs_operations.GetCurrentUser(),
		Timestamp: time.Now(),
	}

	if parentCommit != nil {
		newCommit.Parent = append(newCommit.Parent, parentCommit)
	}

	newCommit.ID, err = vcs_operations.GenerateCommitID(newCommit.Tree, newCommit.Parent, newCommit.Message, newCommit.Author, newCommit.Timestamp)
	if err != nil {
		log.Fatalf("Error generating commit ID: %v", err)
	}

	// Serialize the commit object to JSON
	commitData, err := json.Marshal(newCommit)
	if err != nil {
		log.Fatalf("Error serializing commit data: %v", err)
	}

	// Write Commit Object to File
	commitFilePath := filepath.Join(commitsDir, newCommit.ID)
	if err := os.WriteFile(commitFilePath, commitData, 0644); err != nil {
		log.Fatalf("Error writing commit file: %v", err)
	}

	// Clear the INDEX file after committing
	if err := os.Truncate(".gitx/index", 0); err != nil {
		log.Fatalf("Error clearing INDEX file: %v", err)
	}

	// Update Metadata with the new commit
	metadataFile := ".gitx/metadata.json"
	if err := metadata_operations.UpdateMetadata(metadataFile, newCommit, "", ""); err != nil {
		log.Fatalf("Error updating metadata: %v", err)
	}

	if err := os.WriteFile(branchRefPath, []byte(newCommit.ID), 0644); err != nil {
		log.Fatalf("Error updating branch ref file: %v", err)
	}

	fmt.Printf("Commit created with ID: %s and message: %s\n", newCommit.ID, newCommit.Message)
}

func createInitialCommit() models.Commit {
	// Create an empty tree
	emptyTree := vcs_operations.CreateEmptyTree()

	// Set author and committer information
	author := vcs_operations.GetCurrentUser()
	committer := author

	// Set commit message and timestamp
	message := "Initial commit"
	timestamp := time.Now()

	// Generate commit ID
	commitID, err := vcs_operations.GenerateCommitID(emptyTree, nil, author, committer, timestamp)
	if err != nil {
		log.Fatalf("Error generating commit ID: %v", err)
	}

	// Create the initial commit object
	initialCommit := models.Commit{
		ID:        commitID,
		Parent:    nil, // No parent commit for the initial commit
		Tree:      emptyTree,
		Message:   message,
		Author:    author,
		Committer: committer,
		Timestamp: timestamp,
	}

	return initialCommit
}

// StatusHandler compares the files in the staging area with the tracked files in the metadata and the files in the working directory.
func StatusHandler() {
	// Helper function to convert to relative path
	relativePath := func(path string) string {
		basePath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %v", err)
		}
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return path // Fallback to absolute path if relative path computation fails
		}
		return relPath
	}

	// Step 1: Retrieve tracked files from metadata
	trackedFiles, err := metadata_operations.GetTrackedFiles("metadata.json")
	if err != nil {
		log.Fatalf("Error retrieving tracked files: %v", err)
	}

	// Step 2: Read the INDEX file to get the staging area
	indexFile := ".gitx/index"
	indexEntries, err := vcs_operations.ReadIndexFile(indexFile)
	if err != nil {
		log.Fatalf("Error reading INDEX file: %v", err)
	}

	stagingArea := make(map[string]string)
	for _, entry := range indexEntries {
		relPath := relativePath(entry.Path)
		stagingArea[relPath] = entry.Hash
	}

	// Step 3: Get list of files in the working directory
	workingDirFiles, err := getAllFilesInDir(".")
	if err != nil {
		log.Fatalf("Error retrieving files from working directory: %v", err)
	}

	// Helper function to check if a file path is within the .gitx directory
	isGitxFile := func(path string) bool {
		return strings.HasPrefix(path, ".gitx"+string(os.PathSeparator))
	}

	// Convert tracked files to relative paths for comparison
	relativeTrackedFiles := make(map[string]string)
	for path, hash := range trackedFiles {
		relPath := relativePath(path)
		relativeTrackedFiles[relPath] = hash
	}

	fmt.Println("Tracked Files:")
	for filePath, hashValue := range relativeTrackedFiles {
		fmt.Printf("\t%s: %s\n", filePath, hashValue)
	}

	// Step 4: Compare files
	fmt.Println("Changes to be committed:")
	for filePath, hashValue := range stagingArea {
		if isGitxFile(filePath) {
			continue // Skip .gitx files
		}
		if _, ok := relativeTrackedFiles[filePath]; !ok {
			fmt.Printf("\tnew file: %s\n", filePath)
		} else {
			if hashValue != relativeTrackedFiles[filePath] {
				fmt.Printf("\tmodified: %s\n", filePath)
			}
		}
	}

	fmt.Println("Changes not staged for commit:")
	for _, file := range workingDirFiles {
		if isGitxFile(file) {
			continue // Skip .gitx files
		}
		relativeFile := relativePath(file)
		if _, ok := stagingArea[relativeFile]; !ok {
			if hashValue, err := hash.SHA1Hash(file); err == nil {
				if trackedHash, ok := relativeTrackedFiles[relativeFile]; ok {
					if hashValue != trackedHash {
						fmt.Printf("\tmodified: %s\n", relativeFile)
					}
				} else {
					fmt.Printf("\tuntracked: %s\n", relativeFile)
				}
			} else {
				log.Printf("Error hashing file %s: %v", file, err)
			}
		}
	}
}

// getAllFilesInDir returns a list of all files in a directory.
func getAllFilesInDir(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
