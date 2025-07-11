package handlers

import (
	"attendance-management/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) EnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetEnrollments(w, r)
	case http.MethodPost:
		h.CreateEnrollment(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetEnrollments(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments")
	if err != nil {
		log.Println("등록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		enrollments = append(enrollments, enrollment)
	}
	writeJSON(w, http.StatusOK, enrollments)
}

func (h *Handler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	_, err := h.db.Exec("INSERT INTO enrollments (student_id, class_id, enrolled_date) VALUES ($1, $2, $3)",
		enrollment.StudentID, enrollment.ClassID, enrollment.EnrolledDate.Time)
	if err != nil {
		log.Println("등록 등록 실패:", err)
		http.Error(w, "Failed to create enrollment", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Enrollment created"})
}

func (h *Handler) EnrollmentByIDHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	enrollmentID := getIDFromPath(r.URL.Path)
	if enrollmentID == "" {
		http.Error(w, "Missing enrollment ID", http.StatusBadRequest)
		return
	}
	row := h.db.QueryRow("SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments WHERE enrollment_id = $1", enrollmentID)
	var enrollment models.Enrollment
	err := row.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time)
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

func (h *Handler) UpdateEnrollmentByID(w http.ResponseWriter, r *http.Request) {
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
	result, err := h.db.Exec("UPDATE enrollments SET student_id = $1, class_id = $2, enrolled_date = $3 WHERE enrollment_id = $4",
		enrollment.StudentID, enrollment.ClassID, enrollment.EnrolledDate.Time, enrollmentID)
	if err != nil {
		log.Println("등록 수정 실패:", err)
		http.Error(w, "Failed to update enrollment", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("등록 수정 실패:", err)
		http.Error(w, "Failed to update enrollment", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Enrollment updated"})
}

func (h *Handler) DeleteEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	enrollmentID := getIDFromPath(r.URL.Path)
	if enrollmentID == "" {
		http.Error(w, "Missing enrollment ID", http.StatusBadRequest)
		return
	}
	result, err := h.db.Exec("DELETE FROM enrollments WHERE enrollment_id = $1", enrollmentID)
	if err != nil {
		log.Println("등록 삭제 실패:", err)
		http.Error(w, "Failed to delete enrollment", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("등록 삭제 실패:", err)
		http.Error(w, "Failed to delete enrollment", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Enrollment deleted"})
}

func (h *Handler) EnrollmentByStudentIDHandler(w http.ResponseWriter, r *http.Request) {
	studentID := getIDFromPath(r.URL.Path)
	if studentID == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}
	rows, err := h.db.Query("SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments WHERE student_id = $1", studentID)
	if err != nil {
		log.Println("등록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		enrollments = append(enrollments, enrollment)
	}
	writeJSON(w, http.StatusOK, enrollments)
}

func (h *Handler) EnrollmentByClassIDHandler(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}
	rows, err := h.db.Query("SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments WHERE class_id = $1", classID)
	if err != nil {
		log.Println("등록 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		enrollments = append(enrollments, enrollment)
	}
	writeJSON(w, http.StatusOK, enrollments)
}
