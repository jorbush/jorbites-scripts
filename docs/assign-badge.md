# `assign-badge` Command

Assigns a badge to a specific user.

## Usage
```bash
./jorbites-scripts assign-badge [userID] [badgeName] [flags]
# OR
./jorbites-scripts assign-badge --user-id [userID] --badge [badgeName] [flags]
```

## Arguments & Flags
* `userID`: The target user ID (must be a valid MongoDB ObjectID).
* `badgeName`: The name of the badge to assign.
* `--user-id`: Flag alternative to the positional `userID` argument.
* `--badge`: Flag alternative to the positional `badgeName` argument.
* `-f`, `--force`: Skip remote validation warning and confirmation prompt for unrecognized badge names.

## Features & Logic
1. **Already Owned check**: The command checks the user's current badges first. If the user already possesses the badge, the CLI prints an info log and exits without redundant database writes.
2. **Next.js API Validation**: If the Next.js app is reachable, the command checks the badge name against `/api/badges`. 
3. **Typo Protection**: If the badge name is not recognized, the CLI prints a warning and requests manual confirmation:
   ```
   Warning: Badge 'custom_badge' is not in the list of recognized badges from Jorbites API.
   Do you want to assign it anyway? [y/N]:
   ```
   If `--force` is used, the confirmation prompt is bypassed.
4. **Offline Resilience**: If the Jorbites API is unreachable, a warning is printed, but the assignment proceeds to prevent blocking DB operations.

## Examples
### Standard assignment
```bash
./jorbites-scripts assign-badge 669b7be8f163ac944bc8a16e level_100
```
### Force custom badge
```bash
./jorbites-scripts assign-badge 669b7be8f163ac944bc8a16e custom_badge --force
```
