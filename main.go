package main

import (
	"net/http"

	"attendance-management/internal/database"
	"attendance-management/internal/handlers"
	"attendance-management/internal/repositories"
	"attendance-management/internal/services"

	_ "github.com/lib/pq"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	studentRepo := repositories.NewStudentRepository(db)
	studentService := services.NewStudentService(studentRepo)
	studentHandler := handlers.NewStudentHandler(studentService)
	handler := handlers.NewHandler(db, studentHandler)

	http.ListenAndServe(":8080", handler)
}
