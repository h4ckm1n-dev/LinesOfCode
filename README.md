# Code Line Counter

This script is designed to recursively walk through a directory, count the lines of code in various file types, and display the results in a formatted table. It recognizes a variety of programming languages and configuration files by their extensions and provides a summary of the lines of code per file type.

## Features

- **Recursive Directory Walk**: The script traverses all subdirectories within the specified directory.
- **Line Counting**: Counts the lines of code for each file.
- **File Type Detection**: Identifies the file type based on its extension and provides an appropriate icon.
- **Concurrency**: Uses goroutines and a wait group to process files concurrently for better performance.
- **Formatted Output**: Displays the results in a nicely formatted ASCII table with colored output.

```md
+--------------------------+--------------------------+
| File Type                | Lines of Code            |
+--------------------------+--------------------------+
| 🐹 go                    | 1234                     |
| 🐍 py                    | 567                      |
| ✨ js                    | 890                      |
| 📝 md                    | 234                      |
| ❓ other                 | 45                       |
+--------------------------+--------------------------+
| Total                    | 3000                     |
+--------------------------+--------------------------+
```

## Script Breakdown
- Main Function: Sets up the environment, initializes variables, and starts the directory walk.
- countLines Function: Reads a file and counts the number of lines.
- getFileType Function: Determines the file type and returns an appropriate icon and label.
- processFile Function: Counts the lines of a given file and updates the total counts.
- File Type Icons: A map associating file extensions with icons for better visual representation.
- Concurrency Handling: Uses a wait group and mutex to safely count lines in files concurrently.

## Supported File Types

The script recognizes and counts lines for the following file types:

- Go (`.go`)
- Python (`.py`)
- JavaScript (`.js`)
- TypeScript (`.ts`)
- HTML (`.html`)
- CSS (`.css`)
- Java (`.java`)
- C (`.c`)
- C++ (`.cpp`)
- Ruby (`.rb`)
- PHP (`.php`)
- Rust (`.rs`)
- Shell (`.sh`)
- YAML (`.yaml`, `.yml`)
- JSON (`.json`)
- Markdown (`.md`)
- Lua (`.lua`)
- Terraform (`.tf`)
- Template (`.tpl`)
- Helm (`.helm`)
- Dockerfile (`Dockerfile`)
- Docker Compose (`docker-compose.yml`)
- Other (unrecognized file types)

## Usage

To run the script, simply execute the `main` function. It will use the current working directory as the starting point for counting lines of code.

### Prerequisites

Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/).

### Running the Script

1. Clone or download this repository.
2. Navigate to the directory containing the script.
3. Run the script using the following command:

```sh
go run main.go
```

## Dependencies
The script relies on the following packages:

- fmt: For formatted I/O.
- io: For basic I/O handling.
- io/fs: For filesystem-related interfaces.
- os: For operating system functionality.
- path/filepath: For manipulating filename paths.
- strings: For string manipulation.
- sync: For concurrency primitives.
- github.com/fatih/color: For colored output in the terminal.
