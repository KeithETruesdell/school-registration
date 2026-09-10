package database_test

import (
	"errors"
	"path/filepath"
	"testing"

	"school-app/database"
)

func setupTestDB(t *testing.T) *database.StudentStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return database.NewStudentStore(db)
}

func sampleStudent() database.Student {
	return database.Student{
		FirstName:  "Jane",
		LastName:   "Doe",
		Nickname:   "Janie",
		ParentName: "John Doe",
		Address1:   "123 Main St",
		Address2:   "Apt 4B",
		City:       "Anytown",
		State:      "CA",
		PostalCode: "90210",
		Email:      "jane.doe@example.com",
		Grade:      "10th",
	}
}

func TestListStudentsEmpty(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	students, err := database.ListStudents(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(students) != 0 {
		t.Fatalf("expected 0 students, got %d", len(students))
	}
}

func TestCreateAndGetStudent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	student := sampleStudent()
	created, err := database.CreateStudent(db, student)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}
	if created.ID <= 0 {
		t.Fatalf("expected positive student ID, got %d", created.ID)
	}
	if created.FirstName != student.FirstName || created.LastName != student.LastName {
		t.Fatalf("created student data mismatch: got %+v", created)
	}

	retrieved, err := database.GetStudent(db, created.ID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}
	if retrieved.ID != created.ID {
		t.Fatalf("expected ID %d, got %d", created.ID, retrieved.ID)
	}
	if retrieved.Email != student.Email || retrieved.Grade != student.Grade {
		t.Fatalf("retrieved student data mismatch: got %+v", retrieved)
	}
}

func TestGetStudentNotFound(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	_, err = database.GetStudent(db, 99999)
	if err == nil {
		t.Fatal("expected error for non-existent student, got nil")
	}
	if !errors.Is(err, database.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListStudentsOrdering(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	s1 := sampleStudent()
	s1.LastName = "Smith"
	s1.FirstName = "Alice"

	s2 := sampleStudent()
	s2.LastName = "Adams"
	s2.FirstName = "Bob"

	s3 := sampleStudent()
	s3.LastName = "Adams"
	s3.FirstName = "Anna"

	if _, err := database.CreateStudent(db, s1); err != nil {
		t.Fatalf("failed to create s1: %v", err)
	}
	if _, err := database.CreateStudent(db, s2); err != nil {
		t.Fatalf("failed to create s2: %v", err)
	}
	if _, err := database.CreateStudent(db, s3); err != nil {
		t.Fatalf("failed to create s3: %v", err)
	}

	students, err := database.ListStudents(db)
	if err != nil {
		t.Fatalf("failed to list students: %v", err)
	}
	if len(students) != 3 {
		t.Fatalf("expected 3 students, got %d", len(students))
	}

	// Should be ordered by last_name, first_name: Adams Anna, Adams Bob, Smith Alice
	if students[0].FirstName != "Anna" || students[0].LastName != "Adams" {
		t.Errorf("expected students[0] to be Anna Adams, got %s %s", students[0].FirstName, students[0].LastName)
	}
	if students[1].FirstName != "Bob" || students[1].LastName != "Adams" {
		t.Errorf("expected students[1] to be Bob Adams, got %s %s", students[1].FirstName, students[1].LastName)
	}
	if students[2].FirstName != "Alice" || students[2].LastName != "Smith" {
		t.Errorf("expected students[2] to be Alice Smith, got %s %s", students[2].FirstName, students[2].LastName)
	}
}

func TestUpdateStudent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	created, err := database.CreateStudent(db, sampleStudent())
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	created.FirstName = "Janet"
	created.Grade = "11th"
	err = database.UpdateStudent(db, created)
	if err != nil {
		t.Fatalf("failed to update student: %v", err)
	}

	updated, err := database.GetStudent(db, created.ID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}
	if updated.FirstName != "Janet" || updated.Grade != "11th" {
		t.Fatalf("updated student fields mismatch: got %+v", updated)
	}
}

func TestUpdateStudentNotFound(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	s := sampleStudent()
	s.ID = 99999
	err = database.UpdateStudent(db, s)
	if err == nil {
		t.Fatal("expected error updating non-existent student, got nil")
	}
	if !errors.Is(err, database.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteStudent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	created, err := database.CreateStudent(db, sampleStudent())
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	err = database.DeleteStudent(db, created.ID)
	if err != nil {
		t.Fatalf("failed to delete student: %v", err)
	}

	_, err = database.GetStudent(db, created.ID)
	if !errors.Is(err, database.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestDeleteStudentNotFound(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	err = database.DeleteStudent(db, 99999)
	if err == nil {
		t.Fatal("expected error deleting non-existent student, got nil")
	}
	if !errors.Is(err, database.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestStudentStore(t *testing.T) {
	store := setupTestDB(t)

	s, err := store.Create(sampleStudent())
	if err != nil {
		t.Fatalf("store.Create failed: %v", err)
	}

	got, err := store.Get(s.ID)
	if err != nil {
		t.Fatalf("store.Get failed: %v", err)
	}
	if got.ID != s.ID {
		t.Fatalf("expected ID %d, got %d", s.ID, got.ID)
	}

	list, err := store.List()
	if err != nil {
		t.Fatalf("store.List failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 student in store, got %d", len(list))
	}

	s.Nickname = "NewNick"
	if err := store.Update(s); err != nil {
		t.Fatalf("store.Update failed: %v", err)
	}

	if err := store.Delete(s.ID); err != nil {
		t.Fatalf("store.Delete failed: %v", err)
	}
}
