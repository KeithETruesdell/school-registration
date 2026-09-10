package database

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("student not found")

type Student struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Nickname   string `json:"nickname"`
	ParentName string `json:"parentName"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Email      string `json:"email"`
	Grade      string `json:"grade"`
}

type StudentStore struct {
	db *sql.DB
}

func NewStudentStore(db *sql.DB) *StudentStore {
	return &StudentStore{db: db}
}

func (s *StudentStore) List() ([]Student, error) {
	return ListStudents(s.db)
}

func (s *StudentStore) Get(id int64) (Student, error) {
	return GetStudent(s.db, id)
}

func (s *StudentStore) Create(student Student) (Student, error) {
	return CreateStudent(s.db, student)
}

func (s *StudentStore) Update(student Student) error {
	return UpdateStudent(s.db, student)
}

func (s *StudentStore) Delete(id int64) error {
	return DeleteStudent(s.db, id)
}

func ListStudents(db *sql.DB) ([]Student, error) {
	rows, err := db.Query(`
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
		return nil, err
	}
	defer rows.Close()

	students := []Student{}

	for rows.Next() {
		var student Student

		err := rows.Scan(
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
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func GetStudents(db *sql.DB) ([]Student, error) {
	return ListStudents(db)
}

func GetStudent(db *sql.DB, id int64) (Student, error) {
	var student Student

	err := db.QueryRow(`
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
			return Student{}, ErrNotFound
		}
		return Student{}, err
	}

	return student, nil
}

func GetStudentByID(db *sql.DB, id int64) (Student, error) {
	return GetStudent(db, id)
}

func CreateStudent(db *sql.DB, student Student) (Student, error) {
	res, err := db.Exec(`
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
		return Student{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Student{}, err
	}

	student.ID = id
	return student, nil
}

func UpdateStudent(db *sql.DB, student Student) error {
	res, err := db.Exec(`
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
		student.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func DeleteStudent(db *sql.DB, id int64) error {
	res, err := db.Exec(`DELETE FROM students WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
