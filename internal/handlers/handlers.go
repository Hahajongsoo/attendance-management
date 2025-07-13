package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	db                *sql.DB
	enrollmentHandler *EnrollmentHandler
	studentHandler    *StudentHandler
	attendanceHandler *AttendanceHandler
	classHandler      *ClassHandler
	teacherHandler    *TeacherHandler
}

func NewHandler(db *sql.DB, enrollmentHandler *EnrollmentHandler, studentHandler *StudentHandler, attendanceHandler *AttendanceHandler, classHandler *ClassHandler, teacherHandler *TeacherHandler) *Handler {
	return &Handler{
		db:                db,
		enrollmentHandler: enrollmentHandler,
		studentHandler:    studentHandler,
		attendanceHandler: attendanceHandler,
		classHandler:      classHandler,
		teacherHandler:    teacherHandler,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/students":
		h.studentHandler.StudentHandler(w, r)
	case strings.HasPrefix(path, "/students/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.studentHandler.StudentByIDHandler(w, r)
			return
		}

		if len(segments) >= 3 {
			switch segments[2] {
			case "attendance":
				h.attendanceHandler.AttendanceHandler(w, r)
				return
			case "enrollments":
				h.enrollmentHandler.EnrollmentByStudentIDHandler(w, r)
				return
			}
		}

		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case strings.HasPrefix(path, "/attendance/"):
		h.attendanceHandler.AttendanceByDateHandler(w, r)
		return
	case path == "/teachers":
		h.teacherHandler.TeacherHandler(w, r)
		return
	case strings.HasPrefix(path, "/teachers/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.teacherHandler.TeacherByIDHandler(w, r)
			return

		}
		if len(segments) >= 3 && segments[2] == "classes" {
			h.classHandler.ClassByTeacherIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case path == "/classes":
		h.classHandler.ClassHandler(w, r)
		return
	case strings.HasPrefix(path, "/classes/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.classHandler.ClassByIDHandler(w, r)
			return
		}
		if len(segments) >= 3 && segments[2] == "enrollments" {
			h.enrollmentHandler.EnrollmentByClassIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case path == "/enrollments":
		h.enrollmentHandler.EnrollmentHandler(w, r)
		return
	case strings.HasPrefix(path, "/enrollments/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.enrollmentHandler.EnrollmentByIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
}

func getIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
