package main

import (
	"log"
	"net/http"

	"school-app/api"
	"school-app/database"
)

func main() {
	db, err := database.Open("school.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	studentStore := database.NewStudentStore(db)
	handler := api.NewRouter(studentStore)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("School app running on http://localhost:8080")

	log.Fatal(server.ListenAndServe())
}
