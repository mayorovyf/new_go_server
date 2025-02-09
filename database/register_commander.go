package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// GenerateAPIKey создает случайный API-ключ
func GenerateAPIKey() string {
	bytes := make([]byte, 16) // 16 байт = 32 символа hex
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// RegisterCommander регистрирует нового командира в базе данных
func RegisterCommander(collection *mongo.Collection, c *Commander) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверка существующего пользователя с таким username
	var existingCommander Commander
	err := collection.FindOne(ctx, bson.M{"username": c.Username}).Decode(&existingCommander)
	if err == nil {
		return errors.New("пользователь с таким именем уже существует")
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	// Генерация API-ключа
	c.APIKey = GenerateAPIKey()
	if c.APIKey == "" {
		return errors.New("ошибка генерации API-ключа")
	}

	// Вставка нового командира в базу данных
	_, err = collection.InsertOne(ctx, c)
	if err != nil {
		return err
	}

	return nil
}

// RegisterCommanderHandler обрабатывает HTTP-запрос для регистрации командира
func RegisterCommanderHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var newCommander Commander
		err := json.NewDecoder(r.Body).Decode(&newCommander)
		if err != nil {
			http.Error(w, "не удалось разобрать запрос: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Генерация уникального ID
		newCommander.UniqueID = uuid.New().String()

		// Проверка обязательных полей
		if newCommander.Username == "" || newCommander.Password == "" || newCommander.FirstName == "" || newCommander.LastName == "" {
			http.Error(w, "имя пользователя, пароль, имя и фамилия обязательны", http.StatusBadRequest)
			return
		}

		// Регистрация командира через функцию
		err = RegisterCommander(collection, &newCommander)
		if err != nil {
			http.Error(w, "ошибка регистрации: "+err.Error(), http.StatusConflict)
			return
		}

		// Успешный ответ с API-ключом и уникальным ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "командир успешно зарегистрирован",
			"unique_id": newCommander.UniqueID,
			"api_key":   newCommander.APIKey,
		})
	}
}
