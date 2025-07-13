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

func getDateFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return parts[len(parts)-1]
}

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
