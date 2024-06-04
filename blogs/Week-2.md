# Week 2: Exploring Git Internals

In the second week of developing GitX, we delved into the internals of Git and gained a deeper understanding of its underlying structure and components. This knowledge will be crucial in building a lightweight version control system that aligns with Git's principles.

## Git Internals: Blob, Tree, and Commit

We started by learning about the three fundamental objects in Git: blob, tree, and commit.

- **Blob**: A blob represents the content of a file at a specific point in time. It is essentially a binary representation of the file's contents.

- **Tree**: A tree object represents a directory in Git. It maintains a list of blobs and other tree objects, effectively representing the directory structure.

- **Commit**: A commit object captures a specific state of the repository at a given point in time. It contains references to the tree object representing the project's directory structure, along with metadata such as the author, timestamp, and commit message.

## First Commit in GitX

During this week, we made our first commit in developing the GitX project. This initial commit marks the starting point of our version control history and sets the foundation for future development.

 
## Merkle Tree and Git

In addition to the blob, tree, and commit objects, we also learned about the Merkle tree in Git. The Merkle tree is a data structure that represents the entire state of the repository by recursively hashing the contents of the tree and its sub-trees.

Each commit in Git corresponds to a root node of a Merkle tree. This allows Git to efficiently track changes in the repository by comparing the hash of the current root node with the hash of the previous root node.

## Implementing the Second Commit

To implement the second commit in GitX, we need to make changes to the repository and create a new commit object. This involves modifying files, adding new files, or deleting existing files.

Once the changes are made, we can use the `git add` command to stage the changes and then use the `git commit` command to create a new commit object. The commit object will reference the updated tree object that represents the new state of the repository.

## Subsequent Commits

After the second commit, we can continue making changes to the repository and creating subsequent commits. Each commit will have a unique identifier and will reference the updated tree object.

By creating a series of commits, we can track the history of changes in the repository and easily revert to previous states if needed.

## Conclusion

The second week of developing GitX was focused on gaining a deeper understanding of Git's internals, including the blob, tree, and commit objects. We also made our first commit and explored other key aspects of Git. Armed with this knowledge, we are now well-equipped to continue building our lightweight version control system.

