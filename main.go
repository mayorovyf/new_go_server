package main

import (
	"fmt"
	"log"
	"os"

	"new_go_server/config"
	"new_go_server/database"
)

func main() {
	// Загрузка .env файла
	config.LoadEnv()

	// Подключение к MongoDB
	mongoURI := config.BuildMongoURI()
	client := database.ConnectToMongoDB(mongoURI)
	defer database.DisconnectFromMongoDB(client)

	// Получение ссылки на базу данных
	dbName := os.Getenv("DB_NAME")
	databaseInstance := client.Database(dbName)
	usersCollection := databaseInstance.Collection("users")

	// Создание нового пользователя
	newUser := &database.User{
		Username: "example_user",
		Password: "example_password",
	}

	// Сохранение пользователя в базе данных
	err := newUser.Save(usersCollection)
	if err != nil {
		log.Fatalf("Ошибка сохранения пользователя: %v", err)
	}

	fmt.Println("Пользователь успешно добавлен в базу данных!")
}
