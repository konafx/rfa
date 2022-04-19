package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Participant struct{
	ID string
	Active bool `firestore:"active"`
}

func GetParticipants(ctx context.Context, projectID string) (participants []Participant, err error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return
	}

	iter := client.Collection("rfa-participants").Documents(ctx)
	for {
		var participant Participant

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return participants, err
		}
		
		if err = doc.DataTo(&participant); err != nil {
			return participants, err
		}
		participant.ID = doc.Ref.ID
		if participant.Active {
			participants = append(participants, participant)
		}
	}

	return
}

type Replacer struct{
	Before string `firestore:"before"`
	After string `firestore:"after"`
}

func GetReplacers(ctx context.Context, projectID string) (replacers []Replacer, err error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return
	}

	iter := client.Collection("rfa-replacers").Documents(ctx)
	for {
		var replacer Replacer

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return replacers, err
		}

		if err = doc.DataTo(&replacer); err != nil {
			return replacers, err
		}
		replacers = append(replacers, replacer)
	}

	return
}
