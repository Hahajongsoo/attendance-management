package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"attendance-management/internal/models"
	"attendance-management/internal/services"
)

type StudentHandler struct {
	Service *services.StudentService
}

func NewStudentHandler(s *services.StudentService) *StudentHandler {
	return &StudentHandler{Service: s}
}

func (h *StudentHandler) StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetStudents(w, r)
	case http.MethodPost:
		h.CreateStudent(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *StudentHandler) StudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetStudentByID(w, r, id)
	case http.MethodPut:
		h.UpdateStudentByID(w, r, id)
	case http.MethodDelete:
		h.DeleteStudentByID(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.Service.GetAllStudents()
	if err != nil {
		log.Println("학생 목록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request, id string) {
	student, err := h.Service.GetStudentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, student)
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.Service.CreateStudent(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Student created"})
}

func (h *StudentHandler) UpdateStudentByID(w http.ResponseWriter, r *http.Request, id string) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.Service.UpdateStudent(id, &student); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Update failed", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Student updated"})
}

func (h *StudentHandler) DeleteStudentByID(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Service.DeleteStudent(id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Student deleted"})
}
