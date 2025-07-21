package handlers

import (
	"attendance-management/internal/models"
	"attendance-management/internal/services"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AttendanceHandler struct {
	Service *services.AttendanceService
}

func NewAttendanceHandler(s *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{Service: s}
}
func (h *AttendanceHandler) AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAttendance(w, r)
	case http.MethodPost:
		h.CreateAttendance(w, r)
	case http.MethodPut:
		h.UpdateAttendance(w, r)
	case http.MethodDelete:
		h.DeleteAttendance(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GetAttendance 학생별 특정 날짜 출결을 조회합니다
// @Summary 학생별 특정 날짜 출결 조회
// @Description 특정 학생의 특정 날짜 출결 정보를 조회합니다
// @Tags attendance
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Param date path string true "날짜 (YYYY-MM-DD)"
// @Success 200 {object} models.Attendance
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Attendance not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id}/attendance/{date} [get]
func (h *AttendanceHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	studentID := getIDFromPath(r.URL.Path)
	if studentID == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	date := getDateFromPath(r.URL.Path)
	if date == "" {
		http.Error(w, "Missing date", http.StatusBadRequest)
		return
	}

	attendance, err := h.Service.GetAttendanceByStudentIDAndDate(studentID, date)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Attendance not found", http.StatusNotFound)
		} else {
			log.Println("출결 조회 실패:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, attendance)
}

// CreateAttendance 학생별 특정 날짜 출결을 등록합니다
// @Summary 출결 등록
// @Description 특정 학생의 특정 날짜 출결을 등록합니다
// @Tags attendance
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Param date path string true "날짜 (YYYY-MM-DD)"
// @Param attendance body models.Attendance true "출결 정보"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id}/attendance/{date} [post]
func (h *AttendanceHandler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
	studentID := getIDFromPath(r.URL.Path)
	if studentID == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	date := getDateFromPath(r.URL.Path)
	if date == "" {
		http.Error(w, "Missing date", http.StatusBadRequest)
		return
	}

	var attendance models.Attendance
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if studentID != strconv.Itoa(attendance.StudentID) {
		http.Error(w, "Student ID mismatch", http.StatusBadRequest)
		return
	}
	if date != attendance.Date.Time.Format("2006-01-02") {
		http.Error(w, "Date mismatch", http.StatusBadRequest)
		return
	}

	err := h.Service.CreateAttendance(&attendance)
	if err != nil {
		log.Println("출결 등록 실패:", err)
		http.Error(w, "Failed to create attendance", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Attendance created"})
}

// UpdateAttendance 학생별 특정 날짜 출결을 수정합니다
// @Summary 출결 수정
// @Description 특정 학생의 특정 날짜 출결을 수정합니다
// @Tags attendance
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Param date path string true "날짜 (YYYY-MM-DD)"
// @Param attendance body models.Attendance true "수정할 출결 정보"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Attendance not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id}/attendance/{date} [put]
func (h *AttendanceHandler) UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	studentID := getIDFromPath(r.URL.Path)
	if studentID == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	date := getDateFromPath(r.URL.Path)
	if date == "" {
		http.Error(w, "Missing date", http.StatusBadRequest)
		return
	}

	var attendance models.Attendance
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if studentID != strconv.Itoa(attendance.StudentID) {
		http.Error(w, "Student ID mismatch", http.StatusBadRequest)
		return
	}
	if date != attendance.Date.Time.Format("2006-01-02") {
		http.Error(w, "Date mismatch", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateAttendance(studentID, date, &attendance)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Attendance not found", http.StatusNotFound)
		} else {
			log.Println("출결 수정 실패:", err)
			http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Attendance updated"})
}

// DeleteAttendance 학생별 특정 날짜 출결을 삭제합니다
// @Summary 출결 삭제
// @Description 특정 학생의 특정 날짜 출결을 삭제합니다
// @Tags attendance
// @Accept json
// @Produce json
// @Param student_id path string true "학생 ID"
// @Param date path string true "날짜 (YYYY-MM-DD)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Attendance not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{student_id}/attendance/{date} [delete]
func (h *AttendanceHandler) DeleteAttendance(w http.ResponseWriter, r *http.Request) {
	studentID := getIDFromPath(r.URL.Path)
	if studentID == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	date := getDateFromPath(r.URL.Path)
	if date == "" {
		http.Error(w, "Missing date", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteAttendance(studentID, date)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Attendance not found", http.StatusNotFound)
		} else {
			log.Println("출결 삭제 실패:", err)
			http.Error(w, "Failed to delete attendance", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Attendance deleted"})
}

// getDateFromPath URL 경로에서 날짜를 추출합니다
// @Summary 날짜 추출
// @Description URL 경로에서 날짜를 추출합니다
func getDateFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return parts[len(parts)-1]
}

// AttendanceByDateHandler 날짜별 전체 출결을 조회합니다
// @Summary 날짜별 전체 출결 조회
// @Description 특정 날짜의 모든 학생 출결 정보를 조회합니다
// @Tags attendance
// @Accept json
// @Produce json
// @Param date path string true "날짜 (YYYY-MM-DD)"
// @Success 200 {array} models.Attendance
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /attendance/{date} [get]
func (h *AttendanceHandler) AttendanceByDateHandler(w http.ResponseWriter, r *http.Request) {
	date := getDateFromPath(r.URL.Path)
	if date == "" {
		http.Error(w, "Missing date", http.StatusBadRequest)
		return
	}

	attendances, err := h.Service.GetAttendanceByDate(date)
	if err != nil {
		log.Println("출결 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, attendances)
}
