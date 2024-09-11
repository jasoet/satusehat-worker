# Go Tools for Database and FHIR Healthcare API Interaction

This repository contains a collection of tools written in Go, for interacting with the Database and the FHIR healthcare
API provided by the SatuSehat service.

## Getting Started

- Clone this repository
- Make sure you have [Go SDK 1.22.1](https://golang.org/dl/) installed.
- Install Task Runner
- Execute `task install-dev` to install the development tools
- Execute `task run`
- Execute `task run -- -d -f .config.yaml start` to run the application
- Execute `task test` to run the tests

## Task Runner

This project uses [Taskfile](https://taskfile.dev/) as a task runner.

- To install Taskfile, check the installation guide
  at [https://taskfile.dev/#/installation](https://taskfile.dev/#/installation)
- Use `task <command>` to run a specific task.
- Use `task --list` to list all available tasks.
- Use `task --watch <command>` to run a specific task and watch for changes.

## Tasks

| Task Name        | Description                                                                                                                                                                                                                                          |
|------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `mkdir`          | This task creates a required directory. It also tests if the directory was correctly created.                                                                                                                                                        |
| `vendor`         | This task runs the Go vendor command. execute 'go mod tidy', and then makes all dependencies as vendored so the Go project can be built reliably with 'go mod vendor'.                                                                               |
| `install-dev`    | This task installs the development tools used in the project like goreleaser, and wire.                                                                                                                                                              |
| `build-snapshot` | This task builds a snapshot of the application after making sure the directories are created and dependencies are vendored. It uses goreleaser to build the snapshot.                                                                                |
| `run`            | This task runs the main.go file after ensuring the directories and dependencies are ready.                                                                                                                                                           |
| `test`           | This task runs the Go tests after making sure directories are created and dependencies are vendored. It provides verbosity with -v option, checks for race conditions with -race and profiles code coverage with -coverprofile and -covermode flags. |

It is advised to use these tasks for promoting consistent development practices and reduce the time spent on setup
processes.

## Contributing

We welcome as many contributors as possible. Feel free to fork this repository and submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.