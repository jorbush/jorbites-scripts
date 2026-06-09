# Configuration Options

`jorbites-scripts` reads database connection settings and application URLs from environment variables, CLI flags, or a `.env` file.

## strict `.env` Loading
The CLI will attempt to load a `.env` file containing configuration variables:
1. From the path explicitly specified via the `--env-file <path>` flag.
2. From the current working directory (`./.env`) if no flag is specified.

---

## Required Environment Variables
The following variables are **strictly required**. If they are missing from both the `.env` file and the system environment, the CLI will panic and exit immediately:

* `DATABASE_URL` (or `MONGO_URI`): The MongoDB connection string (e.g., `mongodb+srv://user:pass@cluster.mongodb.net/example_db`).
* `JORBITES_URL`: The base URL of the Jorbites Next.js application (used to validate badge names dynamically).

### Optional Environment Variables
* `MONGO_DB`: The target database name. If not set, it defaults to the database specified in the connection string path, or `example` as a fallback.

---

## Global Flags
You can override any environment variable on the fly using CLI flags:

* `--env-file <path>`: Path to a custom environment file (panics if the specified file does not exist).
* `--db-url <string>`: MongoDB connection URI (overrides `DATABASE_URL` / `MONGO_URI`).
* `--db-name <string>`: MongoDB database name (overrides `MONGO_DB`).
* `--app-url <string>`: Jorbites Next.js app URL (overrides `JORBITES_URL`).

---

## Example Usage with Flags
```bash
# Run using custom DB URL and DB Name
./jorbites-scripts list-badges 669b7be8f163ac944bc8a16e --db-url "mongodb://localhost:27017" --db-name "prod" --app-url "http://localhost:3000"
```
