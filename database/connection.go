package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToMongoDB подключается к MongoDB и возвращает клиент
func ConnectToMongoDB(uri string) *mongo.Client {
	fmt.Println("Подключение к MongoDB по адресу:", uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	fmt.Println("Успешно подключено к MongoDB!")
	return client
}

// DisconnectFromMongoDB отключает клиент от MongoDB
func DisconnectFromMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("Ошибка при отключении от MongoDB: %v", err)
	}

	fmt.Println("Отключение от MongoDB завершено.")
}
