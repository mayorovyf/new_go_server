package database

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoginCommanderHandler обрабатывает HTTP-запрос для логина командира
func LoginCommanderHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Получаем логин и пароль из запроса
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "не удалось разобрать запрос: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Проверка обязательных полей
		if credentials.Username == "" || credentials.Password == "" {
			http.Error(w, "логин и пароль обязательны", http.StatusBadRequest)
			return
		}

		// Проверка, существует ли командир с такими логином и паролем
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var commander Commander
		err = collection.FindOne(ctx, bson.M{"username": credentials.Username, "password": credentials.Password}).Decode(&commander)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "неверный логин или пароль", http.StatusUnauthorized)
			} else {
				http.Error(w, "ошибка при проверке данных: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Если командир найден, отправляем его данные
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(commander)
		if err != nil {
			http.Error(w, "не удалось отправить данные: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
