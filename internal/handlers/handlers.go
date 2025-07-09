package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"attendance-management/internal/models"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/students":
		h.StudentHandler(w, r)
	case strings.HasPrefix(path, "/students/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.StudentByIDHandler(w, r)
			return
		}

		if len(segments) >= 3 && segments[2] == "attendance" {
			h.AttendanceHandler(w, r)
			return
		}

		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case strings.HasPrefix(path, "/attendance/"):
		h.AttendanceByDateHandler(w, r)
	}
}

func (h *Handler) StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetStudents(w, r)
	case http.MethodPost:
		h.CreateStudent(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) StudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetStudentByID(w, r)
	case http.MethodPut:
		h.UpdateStudentByID(w, r)
	case http.MethodDelete:
		h.DeleteStudentByID(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT * FROM students")
	if err != nil {
		log.Println("DB 조회 오류:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.StudentID, &s.Name, &s.Grade, &s.Phone, &s.ParentPhone); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	writeJSON(w, http.StatusOK, students)
}

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("INSERT INTO students (student_id, name, grade, phone, parent_phone) VALUES ($1, $2, $3, $4, $5)",
		student.StudentID, student.Name, student.Grade, student.Phone, student.ParentPhone)
	if err != nil {
		log.Println("학생 등록 실패:", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("학생 등록 실패:", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Student created"})
}

func (h *Handler) UpdateStudentByID(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("UPDATE students SET name = $1, grade = $2, phone = $3, parent_phone = $4 WHERE student_id = $5",
		student.Name, student.Grade, student.Phone, student.ParentPhone, id)
	if err != nil {
		log.Println("학생 수정 실패:", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("학생 수정 실패:", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Student updated"})
}

func (h *Handler) DeleteStudentByID(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("DELETE FROM students WHERE student_id = $1", id)
	if err != nil {
		log.Println("학생 삭제 실패:", err)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Student deleted"})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	var student models.Student
	row := h.db.QueryRow("SELECT * FROM students WHERE student_id = $1", id)
	err := row.Scan(&student.StudentID, &student.Name, &student.Grade, &student.Phone, &student.ParentPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			log.Println("학생 조회 실패:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, student)
}

func getIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

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
	err := row.Scan(&attendance.StudentID, &attendance.Date, &attendance.CheckIn, &attendance.CheckOut, &attendance.Status)
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
	if len(parts) < 4 {
		return ""
	}
	return parts[3]
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
		var a models.Attendance
		if err := rows.Scan(&a.StudentID, &a.Date, &a.CheckIn, &a.CheckOut, &a.Status); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		attendances = append(attendances, a)
	}
	writeJSON(w, http.StatusOK, attendances)
}
