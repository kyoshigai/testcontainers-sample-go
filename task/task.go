package task

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type TaskHandler struct {
	TaskRepository *TaskRepository
}

type TaskRepository struct {
	DB *sql.DB
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.TaskRepository.InsertTask(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (r *TaskRepository) InsertTask(title string) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id
	`, title).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
