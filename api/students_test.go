package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"school-app/database"
	"school-app/model"
)

func setupTestRouter(t *testing.T) http.Handler {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	store := database.NewStudentStore(db)
	return NewRouter(store)
}

func TestAPI_StudentLifecycle(t *testing.T) {
	router := setupTestRouter(t)

	// Initially empty list
	req := httptest.NewRequest(http.MethodGet, "/api/students", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var students []model.Student
	if err := json.NewDecoder(rec.Body).Decode(&students); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(students) != 0 {
		t.Fatalf("expected 0 students, got %d", len(students))
	}

	// Create a student
	newStudent := model.Student{
		FirstName:  "Alice",
		LastName:   "Wonderland",
		Nickname:   "Ali",
		ParentName: "Queen of Hearts",
		Address1:   "100 Rabbit Hole",
		City:       "Fantasy",
		State:      "CA",
		PostalCode: "90210",
		Email:      "alice@wonderland.com",
		Grade:      "9",
	}
	body, _ := json.Marshal(newStudent)
	req = httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var created model.Student
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected non-zero ID for created student")
	}
	if created.FirstName != "Alice" || created.LastName != "Wonderland" {
		t.Fatalf("created student fields mismatch: %+v", created)
	}

	// Get student by ID
	req = httptest.NewRequest(http.MethodGet, "/api/students/1", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var fetched model.Student
	if err := json.NewDecoder(rec.Body).Decode(&fetched); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if fetched.ID != created.ID || fetched.Email != "alice@wonderland.com" {
		t.Fatalf("fetched student mismatch: %+v", fetched)
	}

	// Update student
	fetched.Grade = "10"
	fetched.Nickname = "Alicia"
	updateBody, _ := json.Marshal(fetched)
	req = httptest.NewRequest(http.MethodPut, "/api/students/1", bytes.NewReader(updateBody))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var updated model.Student
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Grade != "10" || updated.Nickname != "Alicia" {
		t.Fatalf("updated student mismatch: %+v", updated)
	}

	// List students now has 1
	req = httptest.NewRequest(http.MethodGet, "/api/students", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	var listAfterCreate []model.Student
	if err := json.NewDecoder(rec.Body).Decode(&listAfterCreate); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(listAfterCreate) != 1 {
		t.Fatalf("expected 1 student, got %d", len(listAfterCreate))
	}

	// Delete student
	req = httptest.NewRequest(http.MethodDelete, "/api/students/1", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rec.Code)
	}

	// Get student after delete -> 404
	req = httptest.NewRequest(http.MethodGet, "/api/students/1", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found after delete, got %d", rec.Code)
	}
}

func TestAPI_Errors(t *testing.T) {
	router := setupTestRouter(t)

	// Get not found
	req := httptest.NewRequest(http.MethodGet, "/api/students/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	// Get invalid ID
	req = httptest.NewRequest(http.MethodGet, "/api/students/invalid", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	// Create invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader([]byte("{invalid json")))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	// Update not found
	validStudent := model.Student{FirstName: "Nobody"}
	body, _ := json.Marshal(validStudent)
	req = httptest.NewRequest(http.MethodPut, "/api/students/999", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	// Update invalid ID
	req = httptest.NewRequest(http.MethodPut, "/api/students/invalid", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	// Update invalid JSON
	req = httptest.NewRequest(http.MethodPut, "/api/students/1", bytes.NewReader([]byte("{bad")))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	// Delete not found
	req = httptest.NewRequest(http.MethodDelete, "/api/students/999", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	// Delete invalid ID
	req = httptest.NewRequest(http.MethodDelete, "/api/students/invalid", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
