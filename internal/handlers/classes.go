package handlers

import (
	"attendance-management/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) ClassHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetClasses(w, r)
	case http.MethodPost:
		h.CreateClass(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetClasses(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT * FROM classes")
	if err != nil {
		log.Println("출결 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ClassID, &class.ClassName, &class.Days, &class.StartTime.Time, &class.EndTime.Time, &class.Price, &class.TeacherID); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		classes = append(classes, class)
	}
	writeJSON(w, http.StatusOK, classes)
}

func (h *Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	_, err := h.db.Exec("INSERT INTO classes (class_id, class_name, days, start_time, end_time, price, teacher_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		class.ClassID, class.ClassName, class.Days, class.StartTime.Time, class.EndTime.Time, class.Price, class.TeacherID)
	if err != nil {
		log.Println("클래스 등록 실패:", err)
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Class created"})
}

func (h *Handler) ClassByIDHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetClassByID(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}
	row := h.db.QueryRow("SELECT * FROM classes WHERE class_id = $1", classID)
	var class models.Class
	err := row.Scan(&class.ClassID, &class.ClassName, &class.Days, &class.StartTime.Time, &class.EndTime.Time, &class.Price, &class.TeacherID)
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

func (h *Handler) UpdateClassByID(w http.ResponseWriter, r *http.Request) {
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
	result, err := h.db.Exec("UPDATE classes SET class_name = $1, days = $2, start_time = $3, end_time = $4, price = $5, teacher_id = $6 WHERE class_id = $7",
		class.ClassName, class.Days, class.StartTime.Time, class.EndTime.Time, class.Price, class.TeacherID, classID)
	if err != nil {
		log.Println("클래스 수정 실패:", err)
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("클래스 수정 실패:", err)
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Class updated"})
}

func (h *Handler) DeleteClassByID(w http.ResponseWriter, r *http.Request) {
	classID := getIDFromPath(r.URL.Path)
	if classID == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}
	result, err := h.db.Exec("DELETE FROM classes WHERE class_id = $1", classID)
	if err != nil {
		log.Println("클래스 삭제 실패:", err)
		http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("클래스 삭제 실패:", err)
		http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Class deleted"})
}

func (h *Handler) ClassByTeacherIDHandler(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}
	rows, err := h.db.Query("SELECT * FROM classes WHERE teacher_id = $1", teacherID)
	if err != nil {
		log.Println("클래스 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ClassID, &class.ClassName, &class.Days, &class.StartTime.Time, &class.EndTime.Time, &class.Price, &class.TeacherID); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		classes = append(classes, class)
	}
	writeJSON(w, http.StatusOK, classes)
}
