package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"school-app/database"
)

type Student = database.Student

func (api *API) listStudents(
	w http.ResponseWriter,
	r *http.Request,
) {
	students, err := database.ListStudents(api.db)
	if err != nil {
		http.Error(
			w,
			"unable to retrieve students",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, students)
}

func (api *API) getStudent(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid student id",
			http.StatusBadRequest,
		)
		return
	}

	student, err := database.GetStudent(api.db, id)
	if err != nil {
		http.Error(
			w,
			"student not found",
			http.StatusNotFound,
		)
		return
	}

	writeJSON(w, http.StatusOK, student)
}

func (api *API) createStudent(
	w http.ResponseWriter,
	r *http.Request,
) {
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	created, err := database.CreateStudent(api.db, student)
	if err != nil {
		http.Error(
			w,
			"unable to create student",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (api *API) updateStudent(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid student id",
			http.StatusBadRequest,
		)
		return
	}

	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	student.ID = id
	err = database.UpdateStudent(api.db, student)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(
				w,
				"student not found",
				http.StatusNotFound,
			)
			return
		}
		http.Error(
			w,
			"unable to update student",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, student)
}

func (api *API) deleteStudent(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid student id",
			http.StatusBadRequest,
		)
		return
	}

	err = database.DeleteStudent(api.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(
				w,
				"student not found",
				http.StatusNotFound,
			)
			return
		}
		http.Error(
			w,
			"unable to delete student",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(value)
}
