package task

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu       sync.RWMutex
	seq      int64
	items    map[int64]*Task
	filePath string
}

// NewRepo теперь принимает путь к файлу для сохранения
func NewRepo(filePath string) *Repo {
	r := &Repo{
		items:    make(map[int64]*Task),
		filePath: filePath,
	}
	_ = r.load() // Загружаем данные при старте, если файл существует
	return r
}

// Внутренний метод сохранения в JSON-файл
func (r *Repo) save() error {
	tasks := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		tasks = append(tasks, t)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

// Внутренний метод загрузки из JSON-файла
func (r *Repo) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файла еще нет — это нормально для первого запуска
		}
		return err
	}

	var tasks []*Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	var maxID int64
	for _, t := range tasks {
		r.items[t.ID] = t
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	r.seq = maxID
	return nil
}

// List с поддержкой фильтрации по done и пагинации (старый метод без аргументов удален)
func (r *Repo) List(doneFilter *bool, page, limit int) []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*Task
	for _, t := range r.items {
		if doneFilter != nil && t.Done != *doneFilter {
			continue
		}
		filtered = append(filtered, t)
	}

	// Пагинация
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	if start > len(filtered) {
		return []*Task{}
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end]
}

func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	now := time.Now()
	t := &Task{ID: r.seq, Title: title, CreatedAt: now, UpdatedAt: now, Done: false}
	r.items[t.ID] = t
	_ = r.save() // Сохраняем на диск при создании
	return t
}

func (r *Repo) Update(id int64, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	_ = r.save() // Сохраняем на диск при обновлении
	return t, nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	_ = r.save() // Сохраняем на диск при удалении
	return nil
}