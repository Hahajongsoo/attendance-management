package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	db             *sql.DB
	studentHandler *StudentHandler
}

func NewHandler(db *sql.DB, studentHandler *StudentHandler) *Handler {
	return &Handler{db: db, studentHandler: studentHandler}
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
				h.AttendanceHandler(w, r)
				return
			case "enrollments":
				h.EnrollmentByStudentIDHandler(w, r)
				return
			}
		}

		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case strings.HasPrefix(path, "/attendance/"):
		h.AttendanceByDateHandler(w, r)
		return
	case path == "/teachers":
		h.TeacherHandler(w, r)
		return
	case strings.HasPrefix(path, "/teachers/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.TeacherByIDHandler(w, r)
			return

		}
		if len(segments) >= 3 && segments[2] == "classes" {
			h.ClassByTeacherIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case path == "/classes":
		h.ClassHandler(w, r)
		return
	case strings.HasPrefix(path, "/classes/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.ClassByIDHandler(w, r)
			return
		}
		if len(segments) >= 3 && segments[2] == "enrollments" {
			h.EnrollmentByClassIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case path == "/enrollments":
		h.EnrollmentHandler(w, r)
		return
	case strings.HasPrefix(path, "/enrollments/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.EnrollmentByIDHandler(w, r)
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
