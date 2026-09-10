package database

import (
	"database/sql"
	"errors"
	"testing"

	"school-app/model"
)

func TestStudentStore_CRUD(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	store := NewStudentStore(db)

	// Create
	id, err := store.CreateStudent(model.Student{
		FirstName:  "John",
		LastName:   "Doe",
		Nickname:   "JD",
		ParentName: "Jane Doe",
		Address1:   "123 Main St",
		City:       "Anytown",
		State:      "NY",
		PostalCode: "12345",
		Email:      "john@example.com",
		Grade:      "10",
	})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero id")
	}

	// Get
	student, err := store.GetStudent(id)
	if err != nil {
		t.Fatalf("getting student: %v", err)
	}
	if student.FirstName != "John" {
		t.Fatalf("expected first name John, got %s", student.FirstName)
	}
	if student.LastName != "Doe" {
		t.Fatalf("expected last name Doe, got %s", student.LastName)
	}

	// List
	students, err := store.ListStudents()
	if err != nil {
		t.Fatalf("listing students: %v", err)
	}
	if len(students) != 1 {
		t.Fatalf("expected 1 student, got %d", len(students))
	}

	// Update
	student.FirstName = "Jane"
	student.LastName = "Smith"
	if err := store.UpdateStudent(id, *student); err != nil {
		t.Fatalf("updating student: %v", err)
	}

	updated, err := store.GetStudent(id)
	if err != nil {
		t.Fatalf("getting updated student: %v", err)
	}
	if updated.FirstName != "Jane" {
		t.Fatalf("expected first name Jane, got %s", updated.FirstName)
	}
	if updated.LastName != "Smith" {
		t.Fatalf("expected last name Smith, got %s", updated.LastName)
	}

	// Delete
	if err := store.DeleteStudent(id); err != nil {
		t.Fatalf("deleting student: %v", err)
	}

	_, err = store.GetStudent(id)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestStudentStore_GetStudent_NotFound(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	store := NewStudentStore(db)

	_, err = store.GetStudent(999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestStudentStore_UpdateStudent_NotFound(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	store := NewStudentStore(db)

	err = store.UpdateStudent(999, model.Student{FirstName: "Ghost", LastName: "User"})
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestStudentStore_DeleteStudent_NotFound(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	store := NewStudentStore(db)

	err = store.DeleteStudent(999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestStudentStore_ListStudents_Ordered(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	store := NewStudentStore(db)

	_, err = store.CreateStudent(model.Student{FirstName: "Charlie", LastName: "Brown"})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}

	_, err = store.CreateStudent(model.Student{FirstName: "Alice", LastName: "Adams"})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}

	_, err = store.CreateStudent(model.Student{FirstName: "Bob", LastName: "Brown"})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}

	students, err := store.ListStudents()
	if err != nil {
		t.Fatalf("listing students: %v", err)
	}

	if len(students) != 3 {
		t.Fatalf("expected 3 students, got %d", len(students))
	}

	// Expected order: Adams, Alice; Brown, Bob; Brown, Charlie
	if students[0].LastName != "Adams" {
		t.Fatalf("expected first student Adams, got %s", students[0].LastName)
	}
	if students[1].FirstName != "Bob" {
		t.Fatalf("expected second student Bob, got %s", students[1].FirstName)
	}
	if students[2].FirstName != "Charlie" {
		t.Fatalf("expected third student Charlie, got %s", students[2].FirstName)
	}
}
