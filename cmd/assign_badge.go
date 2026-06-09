package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	assignCmd.Flags().StringVar(&userIDFlag, "user-id", "", "ID of the user")
	assignCmd.Flags().StringVar(&badgeFlag, "badge", "", "Name of the badge to assign")
	assignCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Skip remote validation against Jorbites API")
	RootCmd.AddCommand(assignCmd)
}

var assignCmd = &cobra.Command{
	Use:   "assign-badge [userID] [badgeName]",
	Short: "Assign a badge to a user",
	Long:  `Assigns the specified badge to the target user. If the badge name is not recognized by the Jorbites API, it prompts for confirmation unless the --force flag is passed.`,
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

		for _, b := range badges {
			if b == badgeName {
				cmd.Printf("Info: User '%s' (%s) already has the badge '%s'. No changes made.\n", userName, userID, badgeName)
				return nil
			}
		}

		if !forceFlag {
			availableBadges, err := fetchAvailableBadges(appURL)
			if err != nil {
				// API failed (offline/network issue). Warn but proceed since we shouldn't block DB modifications.
				cmd.Printf("Warning: Could not contact Jorbites API at %s for validation: %v. Proceeding without validation.\n", appURL, err)
			} else {
				found := false
				for _, b := range availableBadges {
					if b == badgeName {
						found = true
						break
					}
				}

				if !found {
					cmd.Printf("Warning: Badge '%s' is not in the list of recognized badges from Jorbites API.\n", badgeName)
					confirm, err := askForConfirmation("Do you want to assign it anyway? [y/N]: ")
					if err != nil {
						return err
					}
					if !confirm {
						return fmt.Errorf("assignment aborted by user")
					}
				}
			}
		}

		err = dbClient.AssignBadgeToUser(ctx, userID, badgeName)
		if err != nil {
			return err
		}

		cmd.Printf("Success: Badge '%s' assigned to user '%s' (%s).\n", badgeName, userName, userID)
		return nil
	},
}
