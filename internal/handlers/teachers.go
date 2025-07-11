package handlers

import (
	"attendance-management/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) TeacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTeachers(w, r)
	case http.MethodPost:
		h.CreateTeacher(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT teacher_id, name, phone_number FROM teachers")
	if err != nil {
		log.Println("출결 조회 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teachers []models.TeacherResponse
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(&teacher.TeacherID, &teacher.Name, &teacher.Phone); err != nil {
			log.Println("rows.Scan 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		teachers = append(teachers, teacher.ToResponse())
	}
	writeJSON(w, http.StatusOK, teachers)
}

func (h *Handler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("비밀번호 암호화 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	teacher.Password = string(hashedPassword)
	_, err = h.db.Exec("INSERT INTO teachers (teacher_id, password, name, phone_number) VALUES ($1, $2, $3, $4)",
		teacher.TeacherID, teacher.Password, teacher.Name, teacher.Phone)
	if err != nil {
		log.Println("교사 등록 실패:", err)
		http.Error(w, "Failed to create teacher", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Teacher created"})
}

func (h *Handler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}
	row := h.db.QueryRow("SELECT teacher_id, name, phone_number FROM teachers WHERE teacher_id = $1", teacherID)
	var teacher models.Teacher
	err := row.Scan(&teacher.TeacherID, &teacher.Name, &teacher.Phone)
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

func (h *Handler) UpdateTeacherByID(w http.ResponseWriter, r *http.Request) {
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("비밀번호 암호화 실패:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	teacher.Password = string(hashedPassword)
	result, err := h.db.Exec("UPDATE teachers SET name = $1, phone_number = $2, password = $3 WHERE teacher_id = $4",
		teacher.Name, teacher.Phone, teacher.Password, teacherID)
	if err != nil {
		log.Println("교사 수정 실패:", err)
		http.Error(w, "Failed to update teacher", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("교사 수정 실패:", err)
		http.Error(w, "Failed to update teacher", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher updated"})
}

func (h *Handler) DeleteTeacherByID(w http.ResponseWriter, r *http.Request) {
	teacherID := getIDFromPath(r.URL.Path)
	if teacherID == "" {
		http.Error(w, "Missing teacher ID", http.StatusBadRequest)
		return
	}
	result, err := h.db.Exec("DELETE FROM teachers WHERE teacher_id = $1", teacherID)
	if err != nil {
		log.Println("교사 삭제 실패:", err)
		http.Error(w, "Failed to delete teacher", http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("교사 삭제 실패:", err)
		http.Error(w, "Failed to delete teacher", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher deleted"})
}
