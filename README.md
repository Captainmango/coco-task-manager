# coco-task-manager 
*A lightweight cron expression parser written in Go*

`coco-task-manager` is a simple, fast, and self-contained **cron expression parser** that expands a standard cron string into its individual components.  
It follows the common *five-field* cron syntax and outputs the parsed schedule in a clean, human-readable table.

## Features

- Parses standard cron expressions with 5 fields:
  - minute
  - hour
  - day of month
  - month
  - day of week
- Supports:
  - Wildcards (`*`)
  - Lists (e.g. `1,15,30`)
  - Ranges (e.g. `1-5`)
  - Step values (e.g. `*/15`)
  - Singular values (e.g. `4`)
- Nicely formatted output table
- Pure Go — no external dependencies

## Example

### Input:
```bash
go run ./cmd/cli.go "*/15 0 1,15 * 1-5"
```
