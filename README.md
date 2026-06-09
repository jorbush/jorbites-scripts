# Jorbites Admin Scripts

A Go-based command-line interface (CLI) tool designed to manage Jorbites database operations remotely (ideal for deployment on a Raspberry Pi 3B).

## Getting Started

### 1. Prerequisites
- Go version 1.22+
- Access to the Jorbites MongoDB database

### 2. Configuration
The CLI reads database connection settings and application URLs from environment variables, flags, or `.env` files.

See [docs/configuration.md](docs/configuration.md) for detailed configuration options and strict env requirements.

### 3. Build & Installation
To build the CLI tool locally:
```bash
make build
```
For deployment on a Raspberry Pi 3B, see [docs/deployment.md](docs/deployment.md) for cross-compilation and setup instructions.

---

## Available Commands

The CLI supports the following subcommands. For detailed usage manuals, flags, and examples, click on the command name:

* [**`assign-badge`**](docs/assign-badge.md): Assigns a badge to a specific user (includes Next.js API typo validation).
* [**`delete-badge`**](docs/delete-badge.md): Removes an assigned badge from a user.
* [**`list-badges`**](docs/list-badges.md): Lists all badges currently assigned to a user.
* [**`list-all-badges`**](docs/list-all-badges.md): Queries the Next.js API to show all badges available in the system.

---

## Development & Testing

To run the unit tests:
```bash
make test
```

To format and lint the code:
```bash
make lint
```
