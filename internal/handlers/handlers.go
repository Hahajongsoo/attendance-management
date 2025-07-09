package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"attendance-management/internal/models"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetStudents(w, r)
	case "POST":
		h.CreateStudent(w, r)
	case "PUT":
		h.UpdateStudent(w, r)
	case "DELETE":
		h.DeleteStudent(w, r)
	}
}

func (h *Handler) GetStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT * FROM students")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	defer rows.Close()
	students := []models.Student{}
	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.StudentID, &student.Name, &student.Grade, &student.Phone, &student.ParentPhone)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		students = append(students, student)
	}
	json.NewEncoder(w).Encode(students)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	h.db.Exec("INSERT INTO students (student_id, name, grade, phone, parent_phone) VALUES ($1, $2, $3, $4, $5)", student.StudentID, student.Name, student.Grade, student.Phone, student.ParentPhone)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	h.db.Exec("UPDATE students SET name = $1, grade = $2, phone = $3, parent_phone = $4 WHERE student_id = $5", student.Name, student.Grade, student.Phone, student.ParentPhone, student.StudentID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	h.db.Exec("DELETE FROM students WHERE student_id = $1", student.StudentID)
	w.WriteHeader(http.StatusOK)
}
