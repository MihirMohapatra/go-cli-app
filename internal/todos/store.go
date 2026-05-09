package todos

import (
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

type Store struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]Todo
}

func NewStore() *Store {
	return &Store{
		nextID: 1,
		items:  make(map[int64]Todo),
	}
}

func (s *Store) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]Todo, 0, len(s.items))
	for _, todo := range s.items {
		todos = append(todos, todo)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos
}

func (s *Store) Get(id int64) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.items[id]
	return todo, ok
}

func (s *Store) Create(title string) (Todo, error) {
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

func (s *Store) Update(id int64, title string, completed bool) (Todo, error) {
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

func (s *Store) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return false
	}

	delete(s.items, id)
	return true
}
