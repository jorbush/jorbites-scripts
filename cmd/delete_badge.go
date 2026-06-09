package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringVar(&userIDFlag, "user-id", "", "ID of the user")
	deleteCmd.Flags().StringVar(&badgeFlag, "badge", "", "Name of the badge to remove")
	RootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete-badge [userID] [badgeName]",
	Short: "Remove a badge from a user",
	Long:  `Deletes the specified badge from the target user's badges list if they currently possess it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, badgeName, err := parseUserAndBadgeArgs(args)
		if err != nil {
			return err
		}

		if dbClient == nil {
			return fmt.Errorf("database URL is not configured; set DATABASE_URL or MONGO_URI, or pass --db-url")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		badges, userName, err := dbClient.GetUserBadges(ctx, userID)
		if err != nil {
			return err
		}

		hasBadge := false
		for _, b := range badges {
			if b == badgeName {
				hasBadge = true
				break
			}
		}

		if !hasBadge {
			cmd.Printf("Info: User '%s' (%s) does not have the badge '%s'. No changes made.\n", userName, userID, badgeName)
			return nil
		}

		err = dbClient.DeleteBadgeFromUser(ctx, userID, badgeName)
		if err != nil {
			return err
		}

		cmd.Printf("Success: Badge '%s' removed from user '%s' (%s).\n", badgeName, userName, userID)
		return nil
	},
}
