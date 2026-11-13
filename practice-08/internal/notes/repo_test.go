package notes_test

import (
	"context"
	"os"
	"testing"

	"example.com/practice-08/internal/db"
	"example.com/practice-08/internal/notes"
)

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()

	// Подключаемся к тестовой базе данных (можно использовать отдельное имя для тестов)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	deps, err := db.ConnectMongo(ctx, mongoURI, "pz8_test")
	if err != nil {
		t.Fatal(err)
	}

	// Очищаем и удаляем тестовую базу после завершения теста
	t.Cleanup(func() {
		_ = deps.Database.Drop(ctx)
		_ = deps.Client.Disconnect(ctx)
	})

	r, err := notes.NewRepo(deps.Database)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем создание заметки
	created, err := r.Create(ctx, "T1", "C1")
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем получение заметки по ID
	got, err := r.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}

	if got.Title != "T1" {
		t.Fatalf("want T1 got %s", got.Title)
	}
	if got.Content != "C1" {
		t.Fatalf("want C1 got %s", got.Content)
	}
}