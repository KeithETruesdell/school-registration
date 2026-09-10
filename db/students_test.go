package db_test

import (
	"path/filepath"
	"testing"

	"school-app/db"
)

func TestDbPackageQueries(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer database.Close()

	list, err := db.ListStudents(database)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d", len(list))
	}

	created, err := db.CreateStudent(database, db.Student{
		FirstName: "Sam",
		LastName:  "Smith",
		Email:     "sam.smith@example.com",
		Grade:     "12th",
	})
	if err != nil {
		t.Fatalf("CreateStudent failed: %v", err)
	}

	got, err := db.GetStudent(database, created.ID)
	if err != nil {
		t.Fatalf("GetStudent failed: %v", err)
	}
	if got.FirstName != "Sam" {
		t.Errorf("expected first name Sam, got %s", got.FirstName)
	}

	gotByID, err := db.GetStudentByID(database, created.ID)
	if err != nil {
		t.Fatalf("GetStudentByID failed: %v", err)
	}
	if gotByID.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, gotByID.ID)
	}

	all, err := db.GetStudents(database)
	if err != nil {
		t.Fatalf("GetStudents failed: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 student, got %d", len(all))
	}

	got.LastName = "Jones"
	if err := db.UpdateStudent(database, got); err != nil {
		t.Fatalf("UpdateStudent failed: %v", err)
	}

	if err := db.DeleteStudent(database, got.ID); err != nil {
		t.Fatalf("DeleteStudent failed: %v", err)
	}

	store := db.NewStudentStore(database)
	if _, err := store.List(); err != nil {
		t.Fatalf("store.List failed: %v", err)
	}
}
