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

	// Коллекция командиров
	commandersCollection := databaseInstance.Collection("commanders")

	// Установка маршрутов для пользователей
	http.HandleFunc("/register", database.RegisterUserHandler(usersCollection, commandersCollection))
	http.HandleFunc("/login", database.LoginHandler(usersCollection))

	// Установка маршрутов для командиров
	http.HandleFunc("/register_commander", database.RegisterCommanderHandler(commandersCollection))
	http.HandleFunc("/login_commander", database.LoginCommanderHandler(commandersCollection))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
