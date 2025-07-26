package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"todo/model"

	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&model.Todo{})
}

func main() {
	initDB()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	nc.QueueSubscribe("todos.get", "workers", func(msg *nats.Msg) {
		var todos []model.Todo
		db.Find(&todos)
		data, _ := json.Marshal(todos)
		msg.Respond(data)
	})

	nc.QueueSubscribe("todos.create", "workers", func(msg *nats.Msg) {
		var todo model.Todo
		json.Unmarshal(msg.Data, &todo)
		db.Create(&todo)
		resp, _ := json.Marshal(todo)
		msg.Respond(resp)
	})

	nc.QueueSubscribe("todos.update", "workers", func(msg *nats.Msg) {
		var todo model.Todo
		json.Unmarshal(msg.Data, &todo)
		db.Save(&todo)
		resp, _ := json.Marshal(todo)
		msg.Respond(resp)
	})

	nc.QueueSubscribe("todos.delete", "workers", func(msg *nats.Msg) {
		var todo model.Todo
		json.Unmarshal(msg.Data, &todo)
		db.Delete(&todo)
		msg.Respond([]byte(`{"status":"deleted"}`))
	})

	log.Println("DB Service is listening on NATS...")
	select {}
}
