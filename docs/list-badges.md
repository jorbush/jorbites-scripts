# `list-badges` Command

Displays all badges currently assigned to a user.

## Usage
```bash
./jorbites-scripts list-badges [userID] [flags]
# OR
./jorbites-scripts list-badges --user-id [userID] [flags]
```

## Arguments & Flags
* `userID`: The target user ID.
* `--user-id`: Flag alternative to the positional `userID` argument.

## Features & Logic
1. **Projection Query**: Performs a MongoDB projection query on the `User` collection, retrieving only the `name` and `badges` fields to optimize database throughput and bandwidth.
2. **Formatted list**: Outputs the user's name and email, followed by a bulleted count list of their badges. If the user has no badges assigned, it outputs `Badges: [No badges assigned]`.

## Examples
### List badges
```bash
$ ./jorbites-scripts list-badges 669b7be8f163ac944bc8a16e
User: test (669b7be8f163ac944bc8a16e)
Badges (total 3):
  - badge_1st_jorbites_contest
  - badge_2nd_jorbites_contest
  - recipe_of_the_week
```
