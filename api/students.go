package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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

func (api *API) listStudents(
	w http.ResponseWriter,
	r *http.Request,
) {
	rows, err := api.db.Query(`
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
		http.Error(
			w,
			"unable to retrieve students",
			http.StatusInternalServerError,
		)
		return
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
			http.Error(
				w,
				"unable to read student",
				http.StatusInternalServerError,
			)
			return
		}

		students = append(students, student)
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

	var student Student

	err = api.db.QueryRow(`
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
		http.Error(
			w,
			"student not found",
			http.StatusNotFound,
		)
		return
	}

	writeJSON(w, http.StatusOK, student)
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
