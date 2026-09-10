package api

import (
	"database/sql"
	"net/http"
)

type API struct {
	db *sql.DB
}

func NewRouter(db *sql.DB) http.Handler {
	api := &API{
		db: db,
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
