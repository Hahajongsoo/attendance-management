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

// GetClasses 전체 수업 목록을 조회합니다
// @Summary 전체 수업 목록 조회
// @Description 등록된 모든 수업의 목록을 조회합니다
// @Tags classes
// @Accept json
// @Produce json
// @Success 200 {array} models.Class
// @Failure 500 {string} string "Internal Server Error"
// @Router /classes [get]
func (h *ClassHandler) GetClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.Service.GetAllClasses()
	if err != nil {
		log.Println("클래스 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, classes)
}

// CreateClass 새로운 수업을 등록합니다
// @Summary 수업 등록
// @Description 새로운 수업을 등록합니다
// @Tags classes
// @Accept json
// @Produce json
// @Param class body models.Class true "수업 정보"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /classes [post]
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

// GetClassByID 특정 수업의 정보를 조회합니다
// @Summary 특정 수업 조회
// @Description 특정 수업의 정보를 조회합니다
// @Tags classes
// @Accept json
// @Produce json
// @Param class_id path string true "수업 ID"
// @Success 200 {object} models.Class
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Class not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /classes/{class_id} [get]
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

// UpdateClassByID 특정 수업의 정보를 수정합니다
// @Summary 수업 정보 수정
// @Description 특정 수업의 정보를 수정합니다
// @Tags classes
// @Accept json
// @Produce json
// @Param class_id path string true "수업 ID"
// @Param class body models.Class true "수정할 수업 정보"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Class not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /classes/{class_id} [put]
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

// DeleteClassByID 특정 수업을 삭제합니다
// @Summary 수업 삭제
// @Description 특정 수업을 삭제합니다
// @Tags classes
// @Accept json
// @Produce json
// @Param class_id path string true "수업 ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Class not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /classes/{class_id} [delete]
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

// ClassByTeacherIDHandler 교사별 수업 목록을 조회합니다
// @Summary 교사별 수업 목록 조회
// @Description 특정 교사가 담당하는 수업 목록을 조회합니다
// @Tags classes
// @Accept json
// @Produce json
// @Param teacher_id path string true "교사 ID"
// @Success 200 {array} models.Class
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /teachers/{teacher_id}/classes [get]
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
