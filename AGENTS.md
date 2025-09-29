# Agent Coding Guidelines for pbar (Go)

This document outlines the conventions and commands for working in this repository.

## 1. Core Commands

| Action | Command | Notes |
| :--- | :--- | :--- |
| **Build** | `go build -o pbar main.go` | Builds the main executable. |
| **Test All** | `go test -v ./...` | Runs all tests in the project. |
| **Test Single** | `go test -v ./pbar -run TestFunctionName` | Use the full function name. |
| **Lint** | `golangci-lint run ./...` | Requires `golangci-lint` to be installed. |
| **Format** | `go fmt ./...` | Standard Go formatting. |

## 2. Code Style & Conventions (Go)

*   **Formatting:** Strictly adhere to `go fmt` output.
*   **Naming:** Use `CamelCase` for exported functions, types, and struct fields. Use `camelCase` or `snake_case` for internal variables.
*   **Imports:** Use grouped imports (standard library, then third-party).
*   **Error Handling:** Return errors explicitly. Use `fmt.Errorf` for context wrapping.
*   **Structs:** Use JSON struct tags for serialization (e.g., `json:"field_name"`).
*   **Testing:** Use `testing.T` and `t.Run` for clear, isolated subtests.
