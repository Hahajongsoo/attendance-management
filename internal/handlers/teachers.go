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

// GetTeachers 전체 교사 목록을 조회합니다
// @Summary 전체 교사 목록 조회
// @Description 등록된 모든 교사의 목록을 조회합니다
// @Tags teachers
// @Accept json
// @Produce json
// @Success 200 {array} models.TeacherResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /teachers [get]
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

// CreateTeacher 새로운 교사를 등록합니다
// @Summary 교사 등록
// @Description 새로운 교사를 등록합니다
// @Tags teachers
// @Accept json
// @Produce json
// @Param teacher body models.Teacher true "교사 정보"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /teachers [post]
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

// GetTeacherByID 특정 교사의 정보를 조회합니다
// @Summary 특정 교사 조회
// @Description 특정 교사의 정보를 조회합니다
// @Tags teachers
// @Accept json
// @Produce json
// @Param teacher_id path string true "교사 ID"
// @Success 200 {object} models.TeacherResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Teacher not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /teachers/{teacher_id} [get]
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

// UpdateTeacherByID 특정 교사의 정보를 수정합니다
// @Summary 교사 정보 수정
// @Description 특정 교사의 정보를 수정합니다
// @Tags teachers
// @Accept json
// @Produce json
// @Param teacher_id path string true "교사 ID"
// @Param teacher body models.Teacher true "수정할 교사 정보"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Teacher not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /teachers/{teacher_id} [put]
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

// DeleteTeacherByID 특정 교사를 삭제합니다
// @Summary 교사 삭제
// @Description 특정 교사를 삭제합니다
// @Tags teachers
// @Accept json
// @Produce json
// @Param teacher_id path string true "교사 ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Teacher not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /teachers/{teacher_id} [delete]
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
