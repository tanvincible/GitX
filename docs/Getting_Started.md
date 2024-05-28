---
title: Getting Started
---
Welcome to GitX! This guide will help you get started with GitX and set up your first repository.

## Installation Guide

### Prerequisites:
Go (version 1.22.1 or higher)

### Steps

1. **Clone the repository:**
   ```bash
   git clone https://github.com/TanviPooranmal/GitX
   cd GitX
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the project:**
   ```bash
   go build -o gitx cmd/main.go
   ```

4. **Verify the installation:**
   ```bash
   ./gitx --version
   ```

---

## Configuration Guide

### Configuration File

1. **Create a configuration file:**
   Create a `config.yaml` file in the project root.

   ```yaml
   server:
     port: 8080
   database:
     host: localhost
     port: 5432
     user: gitxuser
     password: gitxpassword
     dbname: gitx_db
   ```

2. **Update the configuration file:**
   Edit the `config.yaml` file with your specific settings.

### Environment Variables

1. **Set environment variables (optional):**

   ```bash
   export GITX_SERVER_PORT=8080
   export GITX_DB_HOST=localhost
   export GITX_DB_PORT=5432
   export GITX_DB_USER=gitxuser
   export GITX_DB_PASSWORD=gitxpassword
   export GITX_DB_NAME=gitx_db
   ```

---

## Quick Start Guide

### Running GitX

1. **Start the server:**

   ```bash
   ./gitx --config config.yaml
   ```

2. **Access the application:**
   Open your browser and navigate to `http://localhost:8080`.

### Using GitX

1. **Initialize a new repository:**

   ```bash
   ./gitx init my-repo
   ```

2. **Add files to the repository:**

   ```bash
   ./gitx add file1.txt file2.txt
   ```

3. **Commit changes:**

   ```bash
   ./gitx commit -m "Initial commit"
   ```

4. **View the repository status:**

   ```bash
   ./gitx status
   ```
   