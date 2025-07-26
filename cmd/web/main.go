package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"todo/model"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func initNATS() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable is not set")
	}

	var err error
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	msg, _ := nc.Request("todos.get", nil, nats.DefaultTimeout)
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(msg.Data))
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	data, _ := json.Marshal(todo)
	msg, _ := nc.Request("todos.create", data, nats.DefaultTimeout)
	w.WriteHeader(http.StatusCreated)
	io.Copy(w, bytes.NewReader(msg.Data))
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	data, _ := json.Marshal(todo)
	msg, _ := nc.Request("todos.update", data, nats.DefaultTimeout)
	io.Copy(w, bytes.NewReader(msg.Data))
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	data, _ := json.Marshal(todo)
	msg, _ := nc.Request("todos.delete", data, nats.DefaultTimeout)
	io.Copy(w, bytes.NewReader(msg.Data))
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodos(w, r)
	case http.MethodPost:
		createTodo(w, r)
	case http.MethodPut:
		updateTodo(w, r)
	case http.MethodDelete:
		deleteTodo(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	initNATS()

	http.HandleFunc("/todos", todosHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("API Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
