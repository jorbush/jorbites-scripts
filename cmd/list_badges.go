package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	listBadgesCmd.Flags().StringVar(&userIDFlag, "user-id", "", "ID of the user")
	RootCmd.AddCommand(listBadgesCmd)
}

var listBadgesCmd = &cobra.Command{
	Use:   "list-badges [userID]",
	Short: "List a user's badges",
	Long:  `Queries and prints a formatted list of badges currently assigned to the target user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var userID string
		if len(args) > 0 {
			userID = args[0]
		} else {
			userID = userIDFlag
		}

		if userID == "" {
			return fmt.Errorf("missing user ID; specify it as a positional argument or use --user-id")
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

		cmd.Printf("User: %s (%s)\n", userName, userID)
		if len(badges) == 0 {
			cmd.Println("Badges: [No badges assigned]")
			return nil
		}

		cmd.Printf("Badges (total %d):\n", len(badges))
		for _, b := range badges {
			cmd.Printf("  - %s\n", b)
		}

		return nil
	},
}
