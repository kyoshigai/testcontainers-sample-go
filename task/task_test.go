package task_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testcontainers-sample/infra"
	"testcontainers-sample/task"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func TestCreateTask(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(3*time.Minute),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err = infra.InitDB(dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = infra.CreateTables(db)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title": "test task"}`))
	w := httptest.NewRecorder()

	taskHandler := &task.TaskHandler{
		TaskRepository: &task.TaskRepository{DB: db},
	}
	taskHandler.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Body.String() != "{\"id\":1}\n" {
		t.Errorf("expected body %q, got %q", `{"id":1}`, w.Body.String())
	}

	var id int
	var title string
	err = db.QueryRow("SELECT id, title FROM tasks").Scan(&id, &title)
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Errorf("expected id %d, got %d", 1, id)
	}

	if title != "test task" {
		t.Errorf("expected title %q, got %q", "test task", title)
	}
}
