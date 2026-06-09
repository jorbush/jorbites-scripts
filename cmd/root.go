package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jorbush/jorbites-scripts/pkg/db"
	"github.com/spf13/cobra"
)

var (
	envFile  string
	dbURL    string
	dbName   string
	appURL   string
	dbClient db.DBClient
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "jorbites-scripts",
	Short: "Admin scripts to manage database operations for Jorbites",
	Long:  `A command-line administration tool for performing Jorbites database updates, user badge assignments, and auditing.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		loadEnv()

		if dbURL == "" {
			dbURL = os.Getenv("DATABASE_URL")
			if dbURL == "" {
				dbURL = os.Getenv("MONGO_URI")
			}
		}

		if dbName == "" {
			dbName = os.Getenv("MONGO_DB")
		}

		if appURL == "" {
			appURL = os.Getenv("JORBITES_URL")
		}

		if appURL == "" {
			panic("missing required environment variable: JORBITES_URL")
		}

		if cmd.Name() != "list-all-badges" {
			if dbURL == "" {
				panic("missing required environment variable: DATABASE_URL or MONGO_URI")
			}
			if dbClient == nil {
				client, err := db.NewMongoClient(context.Background(), dbURL, dbName)
				if err != nil {
					panic(fmt.Sprintf("failed to initialize database client: %v", err))
				}
				dbClient = client
			}
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if dbClient != nil {
			return dbClient.Close(context.Background())
		}
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&envFile, "env-file", "", "Path to custom .env file")
	RootCmd.PersistentFlags().StringVar(&dbURL, "db-url", "", "MongoDB connection URI (overrides DATABASE_URL/MONGO_URI)")
	RootCmd.PersistentFlags().StringVar(&dbName, "db-name", "", "MongoDB database name (overrides MONGO_DB)")
	RootCmd.PersistentFlags().StringVar(&appURL, "app-url", "", "Jorbites Next.js application URL (overrides JORBITES_URL)")
}

func loadEnv() {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			panic(fmt.Sprintf("failed to load specified env file '%s': %v", envFile, err))
		}
		return
	}
	_ = godotenv.Load()
}

func GetDBClient() db.DBClient {
	return dbClient
}

func SetDBClient(client db.DBClient) {
	dbClient = client
}

func GetAppURL() string {
	return appURL
}

func SetAppURL(url string) {
	appURL = url
}
