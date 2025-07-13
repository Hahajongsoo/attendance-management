package handlers

import (
	"attendance-management/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"attendance-management/internal/services"
)

type TeacherHandler struct {
	Service *services.TeacherService
}

func NewTeacherHandler(s *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{Service: s}
}

func (h *TeacherHandler) TeacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTeachers(w, r)
	case http.MethodPost:
		h.CreateTeacher(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TeacherHandler) TeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTeacherByID(w, r)
	case http.MethodPut:
		h.UpdateTeacherByID(w, r)
	case http.MethodDelete:
		h.DeleteTeacherByID(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.Service.GetAllTeachers()
	if err != nil {
		log.Println("교사 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var teacherResponses []models.TeacherResponse
	for _, teacher := range teachers {
		teacherResponses = append(teacherResponses, teacher.ToResponse())
	}
	writeJSON(w, http.StatusOK, teacherResponses)
}

func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.CreateTeacher(&teacher)
	if err != nil {
		log.Println("교사 등록 실패:", err)
		http.Error(w, "Failed to create teacher", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Teacher created"})
}

func (h *TeacherHandler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}

	teacher, err := h.Service.GetTeacherByID(teacherID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
		} else {
			log.Println("교사 조회 실패:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, teacher.ToResponse())
}

func (h *TeacherHandler) UpdateTeacherByID(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateTeacher(teacherID, &teacher)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
		} else {
			log.Println("교사 수정 실패:", err)
			http.Error(w, "Failed to update teacher", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher updated"})
}

func (h *TeacherHandler) DeleteTeacherByID(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteTeacher(teacherID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
		} else {
			log.Println("교사 삭제 실패:", err)
			http.Error(w, "Failed to delete teacher", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher deleted"})
}
