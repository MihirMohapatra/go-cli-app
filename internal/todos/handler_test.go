package todos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTodoEndpoints(t *testing.T) {
	handler := NewHandler(NewStore())

	createBody := bytes.NewBufferString(`{"title":"Write CRUD API"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/todos", createBody)
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /todos status = %d, want %d", createRec.Code, http.StatusCreated)
	}

	var created Todo
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode created todo: %v", err)
	}
	if created.ID != 1 || created.Title != "Write CRUD API" || created.Completed {
		t.Fatalf("created todo = %+v", created)
	}

	updateBody := bytes.NewBufferString(`{"title":"Write CRUD API","completed":true}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/todos/1", updateBody)
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("PUT /todos/1 status = %d, want %d", updateRec.Code, http.StatusOK)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /todos/1 status = %d, want %d", getRec.Code, http.StatusOK)
	}

	var got Todo
	if err := json.NewDecoder(getRec.Body).Decode(&got); err != nil {
		t.Fatalf("decode fetched todo: %v", err)
	}
	if !got.Completed {
		t.Fatalf("GET /todos/1 completed = false, want true")
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /todos/1 status = %d, want %d", deleteRec.Code, http.StatusNoContent)
	}
}
