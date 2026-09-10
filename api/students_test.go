package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"school-app/api"
	"school-app/database"
)

func setupTestServer(t *testing.T) (http.Handler, *database.Student) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	initial := database.Student{
		FirstName:  "John",
		LastName:   "Doe",
		Nickname:   "Johnny",
		ParentName: "Jane Doe",
		Address1:   "100 Elm St",
		City:       "Springfield",
		State:      "IL",
		PostalCode: "62701",
		Email:      "john.doe@example.com",
		Grade:      "9th",
	}
	created, err := database.CreateStudent(db, initial)
	if err != nil {
		t.Fatalf("failed to seed student: %v", err)
	}

	handler := api.NewRouter(db)
	return handler, &created
}

func TestListStudents(t *testing.T) {
	handler, seeded := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/students", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var students []api.Student
	if err := json.NewDecoder(rec.Body).Decode(&students); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(students) != 1 {
		t.Fatalf("expected 1 student, got %d", len(students))
	}
	if students[0].ID != seeded.ID {
		t.Errorf("expected student ID %d, got %d", seeded.ID, students[0].ID)
	}
}

func TestGetStudentSuccess(t *testing.T) {
	handler, seeded := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/students/1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var student api.Student
	if err := json.NewDecoder(rec.Body).Decode(&student); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if student.ID != seeded.ID {
		t.Errorf("expected student ID %d, got %d", seeded.ID, student.ID)
	}
	if student.FirstName != seeded.FirstName {
		t.Errorf("expected student first name %s, got %s", seeded.FirstName, student.FirstName)
	}
}

func TestGetStudentNotFound(t *testing.T) {
	handler, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/students/9999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestGetStudentInvalidID(t *testing.T) {
	handler, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/students/abc", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateStudent(t *testing.T) {
	handler, _ := setupTestServer(t)

	payload := api.Student{
		FirstName:  "Alice",
		LastName:   "Wonderland",
		Nickname:   "Ali",
		ParentName: "Queen of Hearts",
		Address1:   "1 Rabbit Hole",
		City:       "Fantasy",
		State:      "NY",
		PostalCode: "10001",
		Email:      "alice@wonderland.com",
		Grade:      "8th",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var created api.Student
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID <= 0 {
		t.Errorf("expected valid positive ID, got %d", created.ID)
	}
	if created.FirstName != payload.FirstName || created.LastName != payload.LastName {
		t.Errorf("expected name %s %s, got %s %s", payload.FirstName, payload.LastName, created.FirstName, created.LastName)
	}
}

func TestCreateStudentInvalidBody(t *testing.T) {
	handler, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader([]byte("not valid json")))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestUpdateStudent(t *testing.T) {
	handler, seeded := setupTestServer(t)

	updated := *seeded
	updated.FirstName = "Jonathan"
	updated.Grade = "10th"
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/students/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp api.Student
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.FirstName != "Jonathan" || resp.Grade != "10th" {
		t.Errorf("expected updated fields, got %+v", resp)
	}
}

func TestUpdateStudentNotFound(t *testing.T) {
	handler, _ := setupTestServer(t)

	payload := api.Student{FirstName: "Ghost"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/students/9999", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestDeleteStudent(t *testing.T) {
	handler, seeded := setupTestServer(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/students/1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	// Verify it's gone
	getReq := httptest.NewRequest(http.MethodGet, "/api/students/1", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after deletion, got %d", getRec.Code)
	}

	_ = seeded
}

func TestDeleteStudentNotFound(t *testing.T) {
	handler, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/students/9999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
