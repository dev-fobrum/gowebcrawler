package db

import (
	"fmt"
	"log"
)

func Insert(collection string, data interface{}) error {
	client, ctx, err := getConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v\n", err)
	}
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection(collection)

	_, err = c.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("Erro ao inserir documento: %v", err)
	}

	return nil
}
