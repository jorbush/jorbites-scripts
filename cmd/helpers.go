package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	userIDFlag    string
	badgeFlag     string
	forceFlag     bool
	readerScanner *bufio.Scanner // helper for testing mock stdin
)

func parseUserAndBadgeArgs(args []string) (string, string, error) {
	var userID, badgeName string

	if len(args) >= 2 {
		userID = args[0]
		badgeName = args[1]
	} else {
		userID = userIDFlag
		badgeName = badgeFlag
	}

	if userID == "" {
		return "", "", fmt.Errorf("missing user ID; specify it as a positional argument or use --user-id")
	}

	if badgeName == "" {
		return "", "", fmt.Errorf("missing badge name; specify it as a positional argument or use --badge")
	}

	return userID, badgeName, nil
}

func fetchAvailableBadges(baseAppURL string) ([]string, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	url := fmt.Sprintf("%s/api/badges", strings.TrimSuffix(baseAppURL, "/"))

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	var badges []string
	if err := json.NewDecoder(resp.Body).Decode(&badges); err != nil {
		return nil, err
	}

	return badges, nil
}

func askForConfirmation(prompt string) (bool, error) {
	_, _ = fmt.Fprint(RootCmd.OutOrStdout(), prompt)
	var text string

	if readerScanner != nil {
		// Used in unit tests
		if readerScanner.Scan() {
			text = readerScanner.Text()
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			text = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return false, err
		}
	}

	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes", nil
}
