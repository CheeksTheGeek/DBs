# Go-DB-Lite

## Table of Contents
1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Project Structure](#project-structure)
6. [Core Components](#core-components)
7. [Database Operations](#database-operations)
8. [Data Types](#data-types)
9. [Command Parsing](#command-parsing)
10. [ANSI Color Support](#ansi-color-support)
11. [B+ Tree Implementation](#b-tree-implementation)
12. [Testing](#testing)
13. [Build and Run](#build-and-run)
14. [Contributing](#contributing)
15. [License](#license)

## Introduction

Go-DB-Lite is a lightweight, file-based database management system implemented in Go. It provides a simple SQL-like interface for basic database operations, including creating databases and tables, inserting data, and querying data. This project is designed as a learning tool and a foundation for understanding database internals.

## Features

- Create and manage multiple databases
- Create tables with various column types
- Insert data into tables
- Basic query support (SELECT)
- In-memory and file-based storage options
- ANSI color-coded console output
- Command history support
- B+ Tree index implementation for efficient data retrieval

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/chaitanyasharma/DBs.git
   ```

2. Navigate to the project directory:
   ```
   cd DBs/go-db-lite
   ```

3. Build the project:
   ```
   make build
   ```

### Usage

To run the database:
```
make run
```

This will start the CLI interface where you can enter commands.

## Supported Commands

- `CREATE TABLE`: Create a new table
- `INSERT`: Insert data into a table
- `SELECT`: Retrieve data from a table
- `UPDATE`: Modify existing data in a table
- `DELETE`: Remove data from a table
- `DROP TABLE`: Delete a table
- `SHOW DATABASES`: List all databases
- `USE`: Switch to a specific database
- `EXIT`: Quit the program

For a full list of commands and their syntax, type `HELP` in the CLI.

## Project Structure

- `cmd/`: Contains the main application entry point
- `internal/`: Internal packages
  - `ansi/`: ANSI color codes for CLI output
  - `parser/`: SQL command parser
  - `tree/`: B+ Tree implementation
  - `types/`: Common type definitions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- This project is inspired by various database systems and is intended for educational purposes.

