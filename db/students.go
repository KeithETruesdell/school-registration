package db

import (
	"database/sql"

	"school-app/database"
)

type Student = database.Student
type StudentStore = database.StudentStore

var (
	ErrNotFound     = database.ErrNotFound
	NewStudentStore = database.NewStudentStore
)

func ListStudents(db *sql.DB) ([]Student, error) {
	return database.ListStudents(db)
}

func GetStudents(db *sql.DB) ([]Student, error) {
	return database.GetStudents(db)
}

func GetStudent(db *sql.DB, id int64) (Student, error) {
	return database.GetStudent(db, id)
}

func GetStudentByID(db *sql.DB, id int64) (Student, error) {
	return database.GetStudentByID(db, id)
}

func CreateStudent(db *sql.DB, student Student) (Student, error) {
	return database.CreateStudent(db, student)
}

func UpdateStudent(db *sql.DB, student Student) error {
	return database.UpdateStudent(db, student)
}

func DeleteStudent(db *sql.DB, id int64) error {
	return database.DeleteStudent(db, id)
}
