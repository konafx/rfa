package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tosh223/rfa/firestore"
	"github.com/tosh223/rfa/search"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Page struct {
	Title string
	Count int
}

func main() {
	var (
		onGcp  = os.Getenv("ON_GCP")
		logger *zap.Logger
		err    error
	)

	if onGcp == "true" {
		logger, err = zapdriver.NewProduction()
	} else {
		logger, err = zapdriver.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("zapdriver failed %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	http.HandleFunc("/", handler)
	http.HandleFunc("/for/participants", participantHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	zap.L().Info(fmt.Sprintf("listening on port %s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		zap.L().Fatal("fatal", zap.Error(err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	query := r.URL.Query()

	var projectID string = os.Getenv("GCP_PROJECT_ID")
	var twitterID string
	var location string
	var size string

	if len(query["projectId"]) > 0 {
		projectID = query["projectId"][0]
	} else if projectID == "" {
		msg := "Parameter[projectId] not found."
		fmt.Fprintf(w, msg)
		zap.L().Warn(msg)
		return
	}

	if len(query["twitterId"]) > 0 {
		twitterID = query["twitterId"][0]
	} else {
		msg := "Parameter[twitterId] not found."
		fmt.Fprintf(w, msg)
		zap.L().Warn(msg)
		return
	}

	if len(query["location"]) > 0 {
		location = query["location"][0]
	} else {
		location = "us"
	}

	if len(query["size"]) > 0 {
		size = query["size"][0]
	} else {
		size = "15"
	}

	var rfa search.Rfa
	rfa.ProjectID = projectID
	rfa.Location = location
	rfa.TwitterID = twitterID
	rfa.Size = size
	err := rfa.Search(ctx)

	if err != nil {
		fmt.Fprintf(w, "Failed %v", err)
		zap.L().Fatal("fatal", zap.Error(err))
	} else {
		fmt.Fprintf(w, "Success")
		zap.L().Info("Success")
	}

	return
}

func participantHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var projectID string = os.Getenv("GCP_PROJECT_ID")

	var rfa search.Rfa
	rfa.ProjectID = projectID
	rfa.Location = "us"
	rfa.Size = "15"

	participants, err := firestore.GetParticipants(ctx, projectID)
	if err != nil {
		msg := fmt.Sprintf("Failed %v", err)
		fmt.Fprintf(w, msg)
		zap.L().Fatal(msg, zap.Error(err))
		return
	}
	zap.S().Debug("participants: %v", participants)

	var eg errgroup.Group
	for _, v := range participants {
		rfa.TwitterID = v.ID
		if rfa.TwitterID == "" {
			msg := "Failed getting TwitterID"
			fmt.Fprintf(w, msg)
			zap.L().Fatal(msg)
			return
		}
		rfa := rfa
		eg.Go(func() error {
			return rfa.Search(ctx)
		})
	}

	if err := eg.Wait(); err != nil {
		msg := fmt.Sprintf("Failed %v", err)
		fmt.Fprintf(w, msg)
		zap.L().Error(msg, zap.Error(err))
	} else {
		fmt.Fprintf(w, "Success")
		zap.L().Info("Success")
	}

	return
}
