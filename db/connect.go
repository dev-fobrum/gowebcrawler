package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getConnection() (client *mongo.Client, ctx context.Context, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err == nil {

			err = client.Ping(ctx, nil)
			if err == nil {
				return client, ctx, nil
			}
		}
		log.Printf("Tentativa %d: Falha ao conectar ao MongoDB: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, nil, fmt.Errorf("Não foi possível conectar ao MongoDB após %d tentativas: %v", maxRetries, err)
}
