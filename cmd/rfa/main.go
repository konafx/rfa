package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tosh223/rfa/search"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use: "rfa",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		projectID, _ := cmd.PersistentFlags().GetString("project-id")
		location, _ := cmd.PersistentFlags().GetString("location")
		twitterID, _ := cmd.PersistentFlags().GetString("twitter-id")
		sizeStr, _ := cmd.PersistentFlags().GetString("search-size")

		var rfa search.Rfa
		rfa.ProjectID = projectID
		rfa.Location = location
		rfa.TwitterID = twitterID
		rfa.Size = sizeStr
		err := rfa.Search(ctx)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed", zap.Error(err))
	}
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zapdriver failed %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("project-id", "p", "", "GCP Project ID")
	rootCmd.PersistentFlags().StringP("location", "l", "us", "BigQuery location")
	rootCmd.PersistentFlags().StringP("twitter-id", "u", "", "Twitter ID")
	rootCmd.PersistentFlags().StringP("search-size", "s", "1", "search size")
}
