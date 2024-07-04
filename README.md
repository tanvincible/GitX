# GitX

GitX is a powerful version control system inspired by Git, designed to provide a simple and efficient way to manage your source code. With GitX, you can enjoy nearly all of Git's features, including support for merkle trees, compression, and metadata operations. Whether you are a solo developer or part of a team, GitX empowers you to track changes, collaborate seamlessly, and ensure the integrity of your codebase. Experience the ease and flexibility of GitX for all your version control needs.

> [!IMPORTANT]
> GitX is actively under development, and hence, it is highly unstable.

## Table of Contents

- [Installation](#installation)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [License](#license)

## Installation

You have two options to install GitX depending on your needs: `git clone` and `go get`.

### Using `git clone`

Use `git clone` if you plan to:
- Contribute to the project.
- Modify the source code.
- Inspect the full repository history.

To clone the repository and build the project, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/TanviPooranmal/GitX
   cd gitx
   ```
2. Download the necessary dependencies:
    ```
    go mod tidy
    ```
3. Build the project:
    ```
    go build -o gitx.exe cmd/main.go
    ```
4. Run the executable:
    ```
    ./gitx <command>
    ```

### Using `go get`

Use `go get` if you want to use GitX as a library in another Go project.

1. To fetch and install GitX as a dependency, run:

    ```
    go get github.com/TanviPooranmal/gitx
    ```
2. Import and use GitX in your Go project:
    ```go
    package main

    import (
        "github.com/TanviPooranmal/gitx"
    )

    func main() {
        // Example usage of GitX library
        gitx.SomeFunction()
    }
    ```

## Project Structure

The project structure is outlined at a high level, as follows:

```
GitX
│   go.mod                   # Go module file
│   LICENSE                  # License file
│   README.md                # Readme file
│
├───cmd/                     # Main application entry point
│       main.go
│
├───docs/                    # Documentation files
│          
├───internal/                # Internal packages
│   │   merkletree.go
│   │
│   ├───compression/         # Compression logic
│   │       compression.go
│   │
│   ├───hash/                # Hashing logic
│   │       hash.go
│   │
│   ├───metadata/            # Metadata handling
│   │       metadata.go
│   │
│   └───storage/             # Storage handling
│           storage.go
│
├───models/                  # Data models
│       blob.go
│       commit.go
│       structs.go
│       tree.go
│
└───utils/                   # Utility packages
    ├───file_operations/     # File operations utility
    │       file_operations.go
    │
    ├───metadata_operations/ # Metadata operations utility
    │       metadata_operations.go
    │
    └───vcs_operations/      # Version control operations utility
            vcs_operations.go
```

## Documentation

You can find the detailed documentation [here](https://tanvipooranmal.github.io/GitX/) and, also in the [docs\ ](https://github.com/TanviPooranmal/GitX/tree/main/docs) directory. It includes guides, references, and troubleshooting information to help you get started with `GitX`.

## License

GitX is licensed under the GNU GENERAL PUBLIC LICENSE. See the [LICENSE](LICENSE.MD) file for more details.
