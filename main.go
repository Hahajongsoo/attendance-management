package main

import (
	"log"
	"net/http"
	"os"

	"attendance-management/internal/database"
	"attendance-management/internal/handlers"
	"attendance-management/internal/repositories"
	"attendance-management/internal/services"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	studentRepo := repositories.NewStudentRepository(db)
	teacherRepo := repositories.NewTeacherRepository(db)
	enrollmentRepo := repositories.NewEnrollmentRepository(db)
	classRepo := repositories.NewClassRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)

	studentService := services.NewStudentService(studentRepo)
	teacherService := services.NewTeacherService(teacherRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo)
	classService := services.NewClassService(classRepo)
	messageSender := services.NewTwilioSender(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), os.Getenv("TWILIO_FROM_NUMBER"))
	notificationService := services.NewNotificationService(studentRepo, messageSender)
	attendanceService := services.NewAttendanceService(attendanceRepo, classRepo, notificationService)

	studentHandler := handlers.NewStudentHandler(studentService)
	teacherHandler := handlers.NewTeacherHandler(teacherService)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	classHandler := handlers.NewClassHandler(classService)

	handler := handlers.NewHandler(db, enrollmentHandler, studentHandler, attendanceHandler, classHandler, teacherHandler)

	http.ListenAndServe(":8080", handler)
}
