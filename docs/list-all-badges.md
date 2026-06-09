# `list-all-badges` Command

Queries the Jorbites Next.js API to fetch and display the complete list of available badges in the system.

## Usage
```bash
./jorbites-scripts list-all-badges [flags]
```

## Features & Logic
1. **Dynamic File Scanning**: Communicates with Jorbites Next.js route `/api/badges` which scans the `public/badges/` directory in real-time, meaning newly added badges are visible immediately without rebuilding the CLI.
2. **Tabular/Bulleted list**: Formats and prints the available badges as a sorted count list.

## Examples
### List all badges
```bash
$ ./jorbites-scripts list-all-badges
Fetching available badges from Jorbites API at http://localhost:3000...
Available Badges in Jorbites (total 53):
  - 10_min_event
  - 29_of_gnocchis
  ...
```
