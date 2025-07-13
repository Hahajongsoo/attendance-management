package handlers

import (
	"attendance-management/internal/models"
	"attendance-management/internal/services"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type ClassHandler struct {
	Service *services.ClassService
}

func NewClassHandler(s *services.ClassService) *ClassHandler {
	return &ClassHandler{Service: s}
}
func (h *ClassHandler) ClassHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetClasses(w, r)
	case http.MethodPost:
		h.CreateClass(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ClassHandler) GetClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.Service.GetAllClasses()
	if err != nil {
		log.Println("클래스 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, classes)
}

func (h *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.CreateClass(&class)
	if err != nil {
		log.Println("클래스 등록 실패:", err)
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Class created"})
}

func (h *ClassHandler) ClassByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetClassByID(w, r)
	case http.MethodPut:
		h.UpdateClassByID(w, r)
	case http.MethodDelete:
		h.DeleteClassByID(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ClassHandler) GetClassByID(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	class, err := h.Service.GetClassByID(classID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			log.Println("클래스 조회 실패:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, class)
}

func (h *ClassHandler) UpdateClassByID(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateClass(classID, &class)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			log.Println("클래스 수정 실패:", err)
			http.Error(w, "Failed to update class", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Class updated"})
}

func (h *ClassHandler) DeleteClassByID(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteClass(classID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			log.Println("클래스 삭제 실패:", err)
			http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Class deleted"})
}

func (h *ClassHandler) ClassByTeacherIDHandler(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}

	classes, err := h.Service.GetClassesByTeacherID(teacherID)
	if err != nil {
		log.Println("클래스 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, classes)
}
