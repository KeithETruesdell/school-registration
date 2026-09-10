package api

import (
	"net/http"

	"school-app/database"
)

type API struct {
	studentStore *database.StudentStore
}

func NewRouter(studentStore *database.StudentStore) http.Handler {
	api := &API{
		studentStore: studentStore,
	}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /api/students",
		api.listStudents,
	)

	mux.HandleFunc(
		"POST /api/students",
		api.createStudent,
	)

	mux.HandleFunc(
		"GET /api/students/{id}",
		api.getStudent,
	)

	mux.HandleFunc(
		"PUT /api/students/{id}",
		api.updateStudent,
	)

	mux.HandleFunc(
		"DELETE /api/students/{id}",
		api.deleteStudent,
	)

	return mux
}
