package handlers

import (
	"attendance-management/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) AttendanceHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetAttendance(w http.ResponseWriter, r *http.Request) {
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

	row := h.db.QueryRow("SELECT * FROM attendance WHERE student_id = $1 AND date = $2", studentID, date)
	var attendance models.Attendance
	err := row.Scan(&attendance.StudentID, &attendance.Date.Time, &attendance.CheckIn.Time, &attendance.CheckOut.Time, &attendance.Status)
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

func (h *Handler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.db.Exec("INSERT INTO attendance (student_id, date, check_in, check_out, status) VALUES ($1, $2, $3, $4, $5)",
		studentID, date, attendance.CheckIn.Time, attendance.CheckOut.Time, attendance.Status)
	if err != nil {
		log.Println("출결 등록 실패:", err)
		http.Error(w, "Failed to create attendance", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Attendance created"})
}

func (h *Handler) UpdateAttendance(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.db.Exec("UPDATE attendance SET check_in = $1, check_out = $2, status = $3 WHERE student_id = $4 AND date = $5",
		attendance.CheckIn.Time, attendance.CheckOut.Time, attendance.Status, studentID, date)
	if err != nil {
		log.Println("출결 수정 실패:", err)
		http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Attendance updated"})
}

func (h *Handler) DeleteAttendance(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.db.Exec("DELETE FROM attendance WHERE student_id = $1 AND date = $2", studentID, date)
	if err != nil {
		log.Println("출결 삭제 실패:", err)
		http.Error(w, "Failed to delete attendance", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Attendance deleted"})
}

func getDateFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return parts[len(parts)-1]
}

func (h *Handler) AttendanceByDateHandler(w http.ResponseWriter, r *http.Request) {
	date := getDateFromPath(r.URL.Path)
	if date == "" {
		http.Error(w, "Missing date", http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query("SELECT * FROM attendance WHERE date = $1 ORDER BY check_in ASC", date)
	if err != nil {
		log.Println("출결 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var attendances []models.Attendance
	for rows.Next() {
		var attendance models.Attendance
		if err := rows.Scan(&attendance.StudentID, &attendance.Date.Time, &attendance.CheckIn.Time, &attendance.CheckOut.Time, &attendance.Status); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		attendances = append(attendances, attendance)
	}
	writeJSON(w, http.StatusOK, attendances)
}
