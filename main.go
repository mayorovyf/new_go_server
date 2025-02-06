package main

import (
	"fmt"
	"log"
	"net/http"
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

	// Коллекция пользователей
	usersCollection := databaseInstance.Collection("users")

	// Установка маршрутов
	http.HandleFunc("/register", database.RegisterHandler(usersCollection))
	http.HandleFunc("/login", database.LoginHandler(usersCollection))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
