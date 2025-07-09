package main

import (
	"net/http"

	"attendance-management/internal/database"
	"attendance-management/internal/handlers"

	_ "github.com/lib/pq"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	handler := handlers.NewHandler(db)
	http.HandleFunc("/students", handler.StudentHandler)
	http.HandleFunc("/students/", handler.StudentByIDHandler)
	http.ListenAndServe(":8080", nil)
}
