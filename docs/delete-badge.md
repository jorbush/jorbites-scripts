# `delete-badge` Command

Removes an assigned badge from a user.

## Usage
```bash
./jorbites-scripts delete-badge [userID] [badgeName] [flags]
# OR
./jorbites-scripts delete-badge --user-id [userID] --badge [badgeName] [flags]
```

## Arguments & Flags
* `userID`: The target user ID.
* `badgeName`: The name of the badge to remove.
* `--user-id`: Flag alternative to the positional `userID` argument.
* `--badge`: Flag alternative to the positional `badgeName` argument.

## Features & Logic
1. **Ownership verification**: The CLI fetches the user's current badges first. If the user does not possess the badge, it prints an info log and exits without updating the database.
2. **MongoDB pull update**: Deletes the badge from the user's `badges` string array in MongoDB via a `$pull` atomic operation.

## Examples
### Remove badge
```bash
./jorbites-scripts delete-badge 669b7be8f163ac944bc8a16e level_100
```
