package userstore

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

type FirestoreImpl struct {
	projectID  string
	collection string
}

func NewFirestoreImpl(projectID, collection string) UserStore {
	return &FirestoreImpl{projectID, collection}
}

type Document struct {
	UserID  string `firestore:"userid"`
	Visited int    `firestore:"visited"`
}

func (f *FirestoreImpl) Increment(ctx context.Context, userID string) (int, error) {
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return 0, fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	doc := Document{UserID: userID, Visited: 1}
	snap, err := client.Collection(f.collection).Doc(userID).Get(ctx)
	if err == nil {
		snap.DataTo(&doc)
		doc.Visited += 1
	}
	result, err := client.Collection(f.collection).Doc(userID).Set(ctx, doc)
	if err != nil {
		return 0, fmt.Errorf("firestore Set failed; %w", err)
	}
	log.Printf("firestore Set result: %v", result)
	return doc.Visited, nil
}
