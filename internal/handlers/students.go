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

// GetStudents 전체 학생 목록을 조회합니다
// @Summary 전체 학생 목록 조회
// @Description 등록된 모든 학생의 목록을 조회합니다
// @Tags students
// @Accept json
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {string} string "Internal Server Error"
// @Router /students [get]
func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.Service.GetAllStudents()
	if err != nil {
		log.Println("학생 목록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, students)
}

// GetStudentByID 특정 학생의 정보를 조회합니다
// @Summary 특정 학생 조회
// @Description 특정 학생의 정보를 조회합니다
// @Tags students
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Student not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id} [get]
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

// CreateStudent 새로운 학생을 등록합니다
// @Summary 학생 등록
// @Description 새로운 학생을 등록합니다
// @Tags students
// @Accept json
// @Produce json
// @Param student body models.Student true "학생 정보"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Router /students [post]
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

// UpdateStudentByID 특정 학생의 정보를 수정합니다
// @Summary 학생 정보 수정
// @Description 특정 학생의 정보를 수정합니다
// @Tags students
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Param student body models.Student true "수정할 학생 정보"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Student not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id} [put]
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

// DeleteStudentByID 특정 학생을 삭제합니다
// @Summary 학생 삭제
// @Description 특정 학생을 삭제합니다
// @Tags students
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id} [delete]
func (h *StudentHandler) DeleteStudentByID(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Service.DeleteStudent(id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Student deleted"})
}
