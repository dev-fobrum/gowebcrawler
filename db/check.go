package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func VisitedLink(link string) bool {
	client, ctx, err := getConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v\n", err)
	}
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")

	opts := options.Count().SetLimit(1)

	n, err := c.CountDocuments(
		context.TODO(),
		bson.D{{Key: "link", Value: link}},
		opts)
	if err != nil {
		panic(err)
	}

	return n > 0
}
