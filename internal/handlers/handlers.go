package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"attendance-management/internal/models"

	"golang.org/x/crypto/bcrypt"
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

		if len(segments) >= 3 {
			switch segments[2] {
			case "attendance":
				h.AttendanceHandler(w, r)
				return
			case "enrollments":
				h.EnrollmentByStudentIDHandler(w, r)
				return
			}
		}

		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case strings.HasPrefix(path, "/attendance/"):
		h.AttendanceByDateHandler(w, r)
		return
	case path == "/teachers":
		h.TeacherHandler(w, r)
		return
	case strings.HasPrefix(path, "/teachers/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.TeacherByIDHandler(w, r)
			return

		}
		if len(segments) >= 3 && segments[2] == "classes" {
			h.ClassByTeacherIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case path == "/classes":
		h.ClassHandler(w, r)
		return
	case strings.HasPrefix(path, "/classes/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.ClassByIDHandler(w, r)
			return
		}
		if len(segments) >= 3 && segments[2] == "enrollments" {
			h.EnrollmentByClassIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case path == "/enrollments":
		h.EnrollmentHandler(w, r)
		return
	case strings.HasPrefix(path, "/enrollments/"):
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) == 2 {
			h.EnrollmentByIDHandler(w, r)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
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
