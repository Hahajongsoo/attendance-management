package handlers

import (
	"attendance-management/internal/models"
	"attendance-management/internal/services"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type EnrollmentHandler struct {
	Service *services.EnrollmentService
}

func NewEnrollmentHandler(s *services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{Service: s}
}

func (h *EnrollmentHandler) EnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetEnrollments(w, r)
	case http.MethodPost:
		h.CreateEnrollment(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EnrollmentHandler) GetEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments, err := h.Service.GetAllEnrollments()
	if err != nil {
		log.Println("등록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, enrollments)
}

func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.CreateEnrollment(&enrollment)
	if err != nil {
		log.Println("등록 등록 실패:", err)
		http.Error(w, "Failed to create enrollment", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Enrollment created"})
}

func (h *EnrollmentHandler) EnrollmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetEnrollmentByID(w, r)
	case http.MethodPut:
		h.UpdateEnrollmentByID(w, r)
	case http.MethodDelete:
		h.DeleteEnrollmentByID(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EnrollmentHandler) GetEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	enrollmentID := getIDFromPath(r.URL.Path)
	if enrollmentID == "" {
		http.Error(w, "Missing enrollment ID", http.StatusBadRequest)
		return
	}

	enrollment, err := h.Service.GetEnrollmentByID(enrollmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Enrollment not found", http.StatusNotFound)
		} else {
			log.Println("등록 조회 실패:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, enrollment)
}

func (h *EnrollmentHandler) UpdateEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	enrollmentID := getIDFromPath(r.URL.Path)
	if enrollmentID == "" {
		http.Error(w, "Missing enrollment ID", http.StatusBadRequest)
		return
	}

	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateEnrollment(enrollmentID, &enrollment)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Enrollment not found", http.StatusNotFound)
		} else {
			log.Println("등록 수정 실패:", err)
			http.Error(w, "Failed to update enrollment", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Enrollment updated"})
}

func (h *EnrollmentHandler) DeleteEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	enrollmentID := getIDFromPath(r.URL.Path)
	if enrollmentID == "" {
		http.Error(w, "Missing enrollment ID", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteEnrollment(enrollmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Enrollment not found", http.StatusNotFound)
		} else {
			log.Println("등록 삭제 실패:", err)
			http.Error(w, "Failed to delete enrollment", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Enrollment deleted"})
}

func (h *EnrollmentHandler) EnrollmentByStudentIDHandler(w http.ResponseWriter, r *http.Request) {
	studentID := getIDFromPath(r.URL.Path)
	if studentID == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	enrollments, err := h.Service.GetEnrollmentsByStudentID(studentID)
	if err != nil {
		log.Println("등록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, enrollments)
}

func (h *EnrollmentHandler) EnrollmentByClassIDHandler(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	enrollments, err := h.Service.GetEnrollmentsByClassID(classID)
	if err != nil {
		log.Println("등록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, enrollments)
}
