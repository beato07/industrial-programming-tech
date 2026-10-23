package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен; если файла нет — ошибка игнорируется
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback — прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:root@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}
	
	// 4) Проверка ListDone (выполненные задачи)
	ctxDone, cancelDone := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelDone()

	doneTasks, err := repo.ListDone(ctxDone, true)
	if err != nil {
		log.Fatalf("ListDone(true) error: %v", err)
	}
	fmt.Println("\n=== Выполненные задачи (done=true) ===")
	if len(doneTasks) == 0 {
		fmt.Println("Нет выполненных задач.")
	} else {
		for _, t := range doneTasks {
			// ИЗМЕНЕНИЕ ЗДЕСЬ: добавили %-24s и %-5v для ровных столбцов
			fmt.Printf("#%d | %-24s | done=%-5v | %s\n", t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
		}
	}

	// 5) Проверка ListDone (невыполненные задачи)
	notDoneTasks, err := repo.ListDone(ctxDone, false)
	if err != nil {
		log.Fatalf("ListDone(false) error: %v", err)
	}
	fmt.Println("\n=== Невыполненные задачи (done=false) ===")
	if len(notDoneTasks) == 0 {
		fmt.Println("Нет невыполненных задач.")
	} else {
		for _, t := range notDoneTasks {
			// ИЗМЕНЕНИЕ ЗДЕСЬ: добавили %-24s и %-5v для ровных столбцов
			fmt.Printf("#%d | %-24s | done=%-5v | %s\n", t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
		}
	}
	
	// 6) Проверка FindByID (используем ctxDone, чтобы не вышел таймаут)
	fmt.Println("\n=== Поиск задачи по ID ===")
	searchID := 1
	task, err := repo.FindByID(ctxDone, searchID) // <-- Лучше ctxDone, а не ctx
	if err != nil {
		log.Fatalf("FindByID error: %v", err)
	}
	
	if task != nil {
		fmt.Printf("Найдена задача #%d:\n", task.ID)
		fmt.Printf("  Название: %s\n", task.Title)
		fmt.Printf("  Выполнено: %v\n", task.Done)
		fmt.Printf("  Создано: %s\n", task.CreatedAt.Format(time.RFC3339))
	} else {
		fmt.Printf("Задача с ID %d не найдена.\n", searchID)
	}
	
	// 7) Проверка CreateMany (Задание 3 - Массовая вставка)
	fmt.Println("\n=== Массовая вставка (CreateMany) ===")
	newTasks := []string{"Массовая задача 1", "Массовая задача 2", "Массовая задача 3"}
	err = repo.CreateMany(ctxDone, newTasks) // Передаем массив строк
	if err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}
	fmt.Println("Успешно вставлено задач:", len(newTasks))
	
	// 8) Измерение статистики пула (Задание 4)
	fmt.Println("\n=== Статистика пула соединений ===")
	stats := db.Stats()
	fmt.Printf("Открыто соединений (OpenConnections): %d\n", stats.OpenConnections)
	fmt.Printf("Используется сейчас (InUse): %d\n", stats.InUse)
	fmt.Printf("Простаивает (Idle): %d\n", stats.Idle)
	fmt.Printf("Максимально открытых (MaxOpenConnections): %d\n", stats.MaxOpenConnections)
}