package todos

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	store *Store
	mux   *http.ServeMux
}

func NewHandler(store *Store) http.Handler {
	h := &Handler{
		store: store,
		mux:   http.NewServeMux(),
	}
	h.routes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("GET /health", h.health)
	h.mux.HandleFunc("GET /todos", h.listTodos)
	h.mux.HandleFunc("POST /todos", h.createTodo)
	h.mux.HandleFunc("GET /todos/{id}", h.getTodo)
	h.mux.HandleFunc("PUT /todos/{id}", h.updateTodo)
	h.mux.HandleFunc("DELETE /todos/{id}", h.deleteTodo)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listTodos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := todoID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, ok := h.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, ErrNotFound.Error())
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req todoRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.store.Create(req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, todo)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := todoID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req todoRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.store.Update(id, req.Title, req.Completed)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := todoID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !h.store.Delete(id) {
		writeError(w, http.StatusNotFound, ErrNotFound.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type todoRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func todoID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("todo id must be a positive number")
	}
	return id, nil
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return errors.New("request body must be valid JSON")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": strings.TrimSpace(message),
	})
}
