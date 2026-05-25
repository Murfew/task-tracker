# task-tracker

A simple command-line task tracker written in Go. Tasks are persisted to a local `tasks.json` file.

## Installation

```bash
git clone https://github.com/murfew/task-tracker.git
cd task-tracker
go build -o task-tracker
```

## Usage

```
task-tracker <command> [arguments]
```

### Commands

| Command | Arguments | Description |
|---|---|---|
| `add` | `<description>` | Add a new task |
| `update` | `<id> <description>` | Update a task's description |
| `delete` | `<id>` | Delete a task |
| `mark-in-progress` | `<id>` | Mark a task as in progress |
| `mark-done` | `<id>` | Mark a task as done |
| `list` | `[status]` | List all tasks, optionally filtered by status |

### Status values

- `todo` (default)
- `in-progress`
- `done`

### Examples

```bash
# Add tasks
task-tracker add "Buy groceries"
task-tracker add "Write unit tests"

# List all tasks
task-tracker list

# List only in-progress tasks
task-tracker list in-progress

# Update a task's description
task-tracker update 1 "Buy groceries and cook dinner"

# Change a task's status
task-tracker mark-in-progress 2
task-tracker mark-done 1

# Delete a task
task-tracker delete 2
```

## Data storage

Tasks are stored in `tasks.json` in the current working directory. The file is created automatically on first use.

## License

[MIT](LICENSE)
