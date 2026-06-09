package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listAllCmd)
}

var listAllCmd = &cobra.Command{
	Use:   "list-all-badges",
	Short: "List all available badges",
	Long:  `Queries the Jorbites API to retrieve and display the complete list of available badges in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Printf("Fetching available badges from Jorbites API at %s...\n", appURL)
		badges, err := fetchAvailableBadges(appURL)
		if err != nil {
			return fmt.Errorf("failed to fetch available badges from Jorbites API: %w", err)
		}

		if len(badges) == 0 {
			cmd.Println("No badges found in the Jorbites app assets.")
			return nil
		}

		cmd.Printf("Available Badges in Jorbites (total %d):\n", len(badges))
		for _, b := range badges {
			cmd.Printf("  - %s\n", b)
		}

		return nil
	},
}
