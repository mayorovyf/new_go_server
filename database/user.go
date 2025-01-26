package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// User представляет пользователя
type User struct {
	Username string `bson:"username"` // Поля BSON соответствуют структуре MongoDB
	Password string `bson:"password"`
}

// Save сохраняет объект пользователя в коллекцию MongoDB
func (u *User) Save(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("ошибка добавления пользователя в базу данных: %v", err)
	}
	return nil
}
