package todos

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
)

var ErrNotFound = errors.New("todo not found")

type Todo struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Store interface {
	List(ctx context.Context) ([]Todo, error)
	Get(ctx context.Context, id int64) (Todo, error)
	Create(ctx context.Context, title string) (Todo, error)
	Update(ctx context.Context, id int64, title string, completed bool) (Todo, error)
	Delete(ctx context.Context, id int64) error
}

type MemoryStore struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]Todo
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID: 1,
		items:  make(map[int64]Todo),
	}
}

func (s *MemoryStore) List(ctx context.Context) ([]Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]Todo, 0, len(s.items))
	for _, todo := range s.items {
		todos = append(todos, todo)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

func (s *MemoryStore) Get(ctx context.Context, id int64) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}

	return todo, nil
}

func (s *MemoryStore) Create(ctx context.Context, title string) (Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Todo{}, errors.New("title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo := Todo{
		ID:    s.nextID,
		Title: title,
	}
	s.items[todo.ID] = todo
	s.nextID++

	return todo, nil
}

func (s *MemoryStore) Update(ctx context.Context, id int64, title string, completed bool) (Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Todo{}, errors.New("title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return Todo{}, ErrNotFound
	}

	todo := Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
	}
	s.items[id] = todo

	return todo, nil
}

func (s *MemoryStore) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}

	delete(s.items, id)
	return nil
}
