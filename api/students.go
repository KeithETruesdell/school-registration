package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"school-app/model"
)

func (api *API) listStudents(
	w http.ResponseWriter,
	r *http.Request,
) {
	students, err := api.studentStore.ListStudents()
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

	student, err := api.studentStore.GetStudent(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(
				w,
				"student not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"unable to retrieve student",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, student)
}

func (api *API) createStudent(
	w http.ResponseWriter,
	r *http.Request,
) {
	var student model.Student

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	id, err := api.studentStore.CreateStudent(student)
	if err != nil {
		http.Error(
			w,
			"unable to create student",
			http.StatusInternalServerError,
		)
		return
	}

	student.ID = id

	writeJSON(w, http.StatusCreated, student)
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

	var student model.Student

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if err := api.studentStore.UpdateStudent(id, student); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	if err := api.studentStore.DeleteStudent(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
