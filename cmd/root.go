package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Used for flags.
	cfgFile           string
	userLicense       string
	collection        *mongo.Collection
	sessionCollection *mongo.Collection
	ctx               context.Context

	rootCmd = &cobra.Command{
		Use:   "passmag",
		Short: "A password manager used to store securely passwords",
		Long:  `passmag`,
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("manager")
	collection = db.Collection("users")

	//defer client.Disconnect(ctx)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
