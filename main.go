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
	teacherRepo := repositories.NewTeacherRepository(db)
	enrollmentRepo := repositories.NewEnrollmentRepository(db)
	classRepo := repositories.NewClassRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)

	studentService := services.NewStudentService(studentRepo)
	teacherService := services.NewTeacherService(teacherRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo)
	classService := services.NewClassService(classRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo)

	studentHandler := handlers.NewStudentHandler(studentService)
	teacherHandler := handlers.NewTeacherHandler(teacherService)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	classHandler := handlers.NewClassHandler(classService)

	handler := handlers.NewHandler(db, enrollmentHandler, studentHandler, attendanceHandler, classHandler, teacherHandler)

	http.ListenAndServe(":8080", handler)
}
