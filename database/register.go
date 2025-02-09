package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
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

// RegisterUserHandler обрабатывает HTTP-запрос для регистрации пользователя
func RegisterUserHandler(usersCollection *mongo.Collection, commandersCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Декодируем тело запроса один раз в объединённую структуру
		var requestData struct {
			Username string `json:"username"`
			Password string `json:"password"`
			APIKey   string `json:"api_key"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "не удалось разобрать запрос: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Проверка наличия API-ключа в запросе
		if requestData.APIKey == "" {
			http.Error(w, "API ключ обязателен", http.StatusBadRequest)
			return
		}

		// Поиск командира по API-ключу
		commander, err := findCommanderByAPIKey(commandersCollection, requestData.APIKey)
		if err != nil {
			http.Error(w, "ошибка проверки API-ключа: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if commander == nil {
			http.Error(w, "неверный API-ключ", http.StatusUnauthorized)
			return
		}

		// Проверка обязательных полей
		if requestData.Username == "" || requestData.Password == "" {
			http.Error(w, "имя пользователя и пароль обязательны", http.StatusBadRequest)
			return
		}

		// Генерация уникального ID пользователя и создание объекта пользователя
		newUser := User{
			Username: requestData.Username,
			Password: requestData.Password,
			UniqueID: uuid.New().String(),
			Class:    commander.UniqueID, // Устанавливаем class как UniqueID командира
		}

		// Регистрация пользователя в базе
		err = newUser.Register(usersCollection)
		if err != nil {
			http.Error(w, "ошибка регистрации: "+err.Error(), http.StatusConflict)
			return
		}

		// Успешная регистрация
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "пользователь успешно зарегистрирован",
			"class":  newUser.Class,
		})
	}
}

// findCommanderByAPIKey ищет командира по API-ключу
func findCommanderByAPIKey(collection *mongo.Collection, apiKey string) (*Commander, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var commander Commander
	err := collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&commander)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Командир не найден
		}
		return nil, err
	}
	return &commander, nil
}
