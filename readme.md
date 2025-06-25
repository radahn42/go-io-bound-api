# Go I/O Bound Task API

This project implements a simple HTTP API in Go for managing long-running I/O bound tasks. The API allows creating, retrieving the status of, and deleting tasks. All task data is stored in memory. It leverages the `slog` package for structured logging and `go-chi/chi` for routing, and `google/uuid` for unique ID generation.

## Features

* **Create Task**: Start a new simulated I/O bound task that runs for 3-5 minutes.
* **Get Task Status**: Retrieve the current status, creation date, and processing duration of a task.
* **Delete Task**: Remove a task from memory.

## Prerequisites

* Go (version 1.21 or higher recommended for `slog`)

## How to Run

1.  **Navigate to the project directory:**

    ```bash
    cd go-io-bound-api
    ```

2.  **Download Go modules:**

    This command will download the necessary dependencies (`github.com/google/uuid` and `github.com/go-chi/chi`).

    ```bash
    go mod tidy
    ```

3.  **Build the application:**

    This will create an executable file named `io-bound-api` (or `io-bound-api.exe` on Windows) in your project directory.

    ```bash
    go build -o io-bound-api ./cmd/main.go
    ```

4.  **Run the application:**

    ```bash
    ./io-bound-api
    ```

    The server will start on `http://localhost:8080`. You will see structured log messages from `slog` indicating server status and task processing events.

## API Endpoints

### 1. Create a New Task

Starts a new long-running task. The server will respond immediately, indicating that the task has been accepted for processing.

* **URL:** `/tasks`
* **Method:** `POST`
* **Request Body:** None
* **Response (202 Accepted):**

    ```json
    {
        "id": "c40a0dc8-56ae-4433-acdd-9ec5254783aa",
        "message": "Task accepted and processing",
        "status": "pending"
    }
    ```

  _Note: The ID format is a standard UUID (Universally Unique Identifier). The status will internally change to "running" very quickly as the background goroutine starts._

**Example using `curl`:**

```bash
curl -X POST http://localhost:8080/tasks
```

### 2. Get Task Status

Retrieves the status and timestamps of a task by its ID.

* **URL:** `/tasks/{id}`
* **Method:** `GET`
* **Request Body:** None
* **Response (200 OK):**

    ```json
    {
        "id": "c40a0dc8-56ae-4433-acdd-9ec5254783aa",
        "status": "running",
        "created_at": "2025-06-25T15:21:54.9807805+03:00",
        "started_at": "2025-06-25T15:21:54.9807805+03:00",
        "completed_at": "0001-01-01T00:00:00Z",
        "duration": "10s"
    }
    ```

**Example using `curl`:**

```bash
curl -X GET http://localhost:8080/tasks/c40a0dc8-56ae-4433-acdd-9ec5254783aa
```

### 3. Delete Task

Deletes a task by its ID from memory. This does not cancel a running task, only removes its record.
* **URL:** `/tasks/{id}`
* **Method:** `DELETE`
* **Request Body:** None
* **Response (204 No Content):**

**Example using `curl`:**

```bash
curl -X DELETE http://localhost:8080/tasks/c40a0dc8-56ae-4433-acdd-9ec5254783aa
```

## Tests

```bash
go test -v ./...
```