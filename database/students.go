package database

import (
	"database/sql"
	"errors"
	"fmt"

	"school-app/model"
)

// Student is an alias for model.Student.
type Student = model.Student

// StudentStore handles all database operations for students.
type StudentStore struct {
	db *sql.DB
}

// NewStudentStore creates a new StudentStore.
func NewStudentStore(db *sql.DB) *StudentStore {
	return &StudentStore{db: db}
}

// ListStudents returns all students ordered by last name, then first name.
func (s *StudentStore) ListStudents() ([]model.Student, error) {
	rows, err := s.db.Query(`
		SELECT
			id,
			first_name,
			last_name,
			nickname,
			parent_name,
			address1,
			address2,
			city,
			state,
			postal_code,
			email,
			grade
		FROM students
		ORDER BY last_name, first_name
	`)
	if err != nil {
		return nil, fmt.Errorf("querying students: %w", err)
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var student model.Student
		if err := rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Nickname,
			&student.ParentName,
			&student.Address1,
			&student.Address2,
			&student.City,
			&student.State,
			&student.PostalCode,
			&student.Email,
			&student.Grade,
		); err != nil {
			return nil, fmt.Errorf("scanning student: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating students: %w", err)
	}

	return students, nil
}

// GetStudent returns a single student by ID.
func (s *StudentStore) GetStudent(id int64) (*model.Student, error) {
	var student model.Student

	err := s.db.QueryRow(`
		SELECT
			id,
			first_name,
			last_name,
			nickname,
			parent_name,
			address1,
			address2,
			city,
			state,
			postal_code,
			email,
			grade
		FROM students
		WHERE id = ?
	`, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Nickname,
		&student.ParentName,
		&student.Address1,
		&student.Address2,
		&student.City,
		&student.State,
		&student.PostalCode,
		&student.Email,
		&student.Grade,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("getting student %d: %w", id, err)
	}

	return &student, nil
}

// CreateStudent inserts a new student and returns the generated ID.
func (s *StudentStore) CreateStudent(student model.Student) (int64, error) {
	result, err := s.db.Exec(`
		INSERT INTO students (
			first_name,
			last_name,
			nickname,
			parent_name,
			address1,
			address2,
			city,
			state,
			postal_code,
			email,
			grade
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		student.FirstName,
		student.LastName,
		student.Nickname,
		student.ParentName,
		student.Address1,
		student.Address2,
		student.City,
		student.State,
		student.PostalCode,
		student.Email,
		student.Grade,
	)
	if err != nil {
		return 0, fmt.Errorf("inserting student: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert id: %w", err)
	}

	return id, nil
}

// UpdateStudent updates an existing student by ID.
func (s *StudentStore) UpdateStudent(id int64, student model.Student) error {
	result, err := s.db.Exec(`
		UPDATE students SET
			first_name = ?,
			last_name = ?,
			nickname = ?,
			parent_name = ?,
			address1 = ?,
			address2 = ?,
			city = ?,
			state = ?,
			postal_code = ?,
			email = ?,
			grade = ?
		WHERE id = ?
	`,
		student.FirstName,
		student.LastName,
		student.Nickname,
		student.ParentName,
		student.Address1,
		student.Address2,
		student.City,
		student.State,
		student.PostalCode,
		student.Email,
		student.Grade,
		id,
	)
	if err != nil {
		return fmt.Errorf("updating student %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteStudent removes a student by ID.
func (s *StudentStore) DeleteStudent(id int64) error {
	result, err := s.db.Exec(`
		DELETE FROM students
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("deleting student %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
