package main

import (
	"log"
	"net/http"

	"testcontainers-sample/infra"
	"testcontainers-sample/task"

	_ "github.com/lib/pq"
)

func main() {
	db, _ := infra.InitDB("postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	_ = infra.CreateTables(db)

	taskRepo := &task.TaskRepository{DB: db}
	taskHandler := &task.TaskHandler{TaskRepository: taskRepo}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)

	log.Println("Listening on :9999")
	log.Fatal(http.ListenAndServe(":9999", mux))
}
