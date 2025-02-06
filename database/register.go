package database

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register регистрирует нового пользователя в коллекции
func (u *User) Register(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверка, существует ли пользователь с таким же уникальным ID
	var existingUser User
	err := collection.FindOne(ctx, bson.M{"unique_id": u.UniqueID}).Decode(&existingUser)
	if err == nil {
		return errors.New("пользователь с таким уникальным ID уже существует")
	} else if err != mongo.ErrNoDocuments {
		return err // Вернуть ошибку, если она связана не с отсутствием документов
	}

	// Сохранение нового пользователя
	_, err = collection.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

// RegisterHandler обрабатывает HTTP-запрос для регистрации пользователя
func RegisterHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var newUser User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "не удалось разобрать запрос: "+err.Error(), http.StatusBadRequest)
			return
		}

		if newUser.Username == "" || newUser.Password == "" || newUser.UniqueID == "" {
			http.Error(w, "имя пользователя, уникальный ID и пароль обязательны", http.StatusBadRequest)
			return
		}

		err = newUser.Register(collection)
		if err != nil {
			http.Error(w, "ошибка регистрации: "+err.Error(), http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"пользователь успешно зарегистрирован"}`))
	}
}
