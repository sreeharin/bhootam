# Bhootam

Bhootam is an experimental, in-memory task queue library written in Go. It's designed as a toy project for learning and experimentation, not intended for production use. It provides basic task queuing functionality with support for timeouts, retries, and panic recovery.

## Features

- **In-Memory Queue**: Simple channel-based job queue for asynchronous task execution.
- **Timeout Support**: Configurable timeouts for tasks to prevent hanging operations.
- **Retry Mechanism**: Automatic retry for failed jobs with configurable retry counts.
- **Panic Recovery**: Workers recover from panics in task functions.

## License

See LICENSE file for details.
