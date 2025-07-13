package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"attendance-management/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestNewHandler(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)
	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.db)
}

// StudentHandler 테스트
func TestStudentHandler(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "GET 요청 - 모든 학생 조회",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
					AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321")
				mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)
			},
		},
		{
			name:           "POST 요청 - 학생 생성",
			method:         http.MethodPost,
			expectedStatus: http.StatusCreated,
			setupMock: func() {
				mock.ExpectExec("INSERT INTO students").WithArgs(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:           "PUT 요청 - Method Not Allowed",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var req *http.Request
			if tt.method == http.MethodPost {
				student := models.Student{StudentID: 1, Name: "김철수", Grade: "1학년", Phone: "010-1234-5678", ParentPhone: "010-8765-4321"}
				jsonData, _ := json.Marshal(student)
				req = httptest.NewRequest(tt.method, "/students", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, "/students", nil)
			}

			w := httptest.NewRecorder()
			handler.studentHandler.StudentHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// StudentByIDHandler 테스트
func TestStudentByIDHandler(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMock      func()
		body           interface{}
	}{
		{
			name:           "GET 요청 - 특정 학생 조회",
			method:         http.MethodGet,
			path:           "/students/1",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
					AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321")
				mock.ExpectQuery("SELECT \\* FROM students WHERE student_id = \\$1").WithArgs("1").WillReturnRows(rows)
			},
		},
		{
			name:           "PUT 요청 - 학생 정보 수정",
			method:         http.MethodPut,
			path:           "/students/1",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mock.ExpectExec("UPDATE students SET name = \\$1, grade = \\$2, phone = \\$3, parent_phone = \\$4 WHERE student_id = \\$5").
					WithArgs("김철수", "2학년", "010-1234-5678", "010-8765-4321", "1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			body: models.Student{Name: "김철수", Grade: "2학년", Phone: "010-1234-5678", ParentPhone: "010-8765-4321"},
		},
		{
			name:           "DELETE 요청 - 학생 삭제",
			method:         http.MethodDelete,
			path:           "/students/1",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mock.ExpectExec("DELETE FROM students WHERE student_id = \\$1").
					WithArgs("1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var req *http.Request
			if tt.body != nil {
				jsonData, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			w := httptest.NewRecorder()
			handler.studentHandler.StudentByIDHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// GetStudents 테스트
func TestStudentRepository_GetAll(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	expectedStudents := []models.Student{
		{StudentID: 1, Name: "김철수", Grade: "1학년", Phone: "010-1234-5678", ParentPhone: "010-8765-4321"},
		{StudentID: 2, Name: "이영희", Grade: "2학년", Phone: "010-2345-6789", ParentPhone: "010-9876-5432"},
	}

	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321").
		AddRow(2, "이영희", "2학년", "010-2345-6789", "010-9876-5432")

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	w := httptest.NewRecorder()

	handler.studentHandler.GetStudents(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseStudents []models.Student
	err := json.Unmarshal(w.Body.Bytes(), &responseStudents)
	require.NoError(t, err)
	assert.Equal(t, expectedStudents, responseStudents)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// CreateStudent 테스트
func TestCreateStudent(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	student := models.Student{
		StudentID:   3,
		Name:        "박민수",
		Grade:       "3학년",
		Phone:       "010-3456-7890",
		ParentPhone: "010-0987-6543",
	}

	mock.ExpectExec("INSERT INTO students").WithArgs(3, "박민수", "3학년", "010-3456-7890", "010-0987-6543").
		WillReturnResult(sqlmock.NewResult(1, 1))

	jsonData, _ := json.Marshal(student)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.studentHandler.CreateStudent(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Student created", response["message"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

// GetStudentByID 테스트
func TestGetStudentByID(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	expectedStudent := models.Student{
		StudentID:   1,
		Name:        "김철수",
		Grade:       "1학년",
		Phone:       "010-1234-5678",
		ParentPhone: "010-8765-4321",
	}

	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321")

	mock.ExpectQuery("SELECT \\* FROM students WHERE student_id = \\$1").WithArgs("1").WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/students/1", nil)
	w := httptest.NewRecorder()

	handler.studentHandler.GetStudentByID(w, req, "1")

	assert.Equal(t, http.StatusOK, w.Code)

	var responseStudent models.Student
	err := json.Unmarshal(w.Body.Bytes(), &responseStudent)
	require.NoError(t, err)
	assert.Equal(t, expectedStudent, responseStudent)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// UpdateStudentByID 테스트
func TestUpdateStudentByID(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	student := models.Student{
		Name:        "김철수",
		Grade:       "2학년",
		Phone:       "010-1234-5678",
		ParentPhone: "010-8765-4321",
	}

	mock.ExpectExec("UPDATE students SET name = \\$1, grade = \\$2, phone = \\$3, parent_phone = \\$4 WHERE student_id = \\$5").
		WithArgs("김철수", "2학년", "010-1234-5678", "010-8765-4321", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	jsonData, _ := json.Marshal(student)
	req := httptest.NewRequest(http.MethodPut, "/students/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.studentHandler.UpdateStudentByID(w, req, "1")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Student updated", response["message"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

// DeleteStudentByID 테스트
func TestDeleteStudentByID(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	mock.ExpectExec("DELETE FROM students WHERE student_id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/students/1", nil)
	w := httptest.NewRecorder()

	handler.studentHandler.DeleteStudentByID(w, req, "1")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Student deleted", response["message"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

// 에러 케이스 테스트들
func TestGetStudentsDBError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	w := httptest.NewRecorder()

	handler.studentHandler.GetStudents(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStudentInvalidJSON(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.studentHandler.CreateStudent(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetStudentByIDMissingID(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/students/", nil)
	w := httptest.NewRecorder()

	handler.studentHandler.GetStudentByID(w, req, "999")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetStudentByIDNotFound(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db, nil, nil, nil, nil, nil)

	mock.ExpectQuery("SELECT \\* FROM students WHERE student_id = \\$1").WithArgs("999").WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/students/999", nil)
	w := httptest.NewRecorder()

	handler.studentHandler.GetStudentByID(w, req, "999")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// getIDFromPath 함수 테스트
func TestGetIDFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"정상적인 ID", "/students/1", "1"},
		{"숫자가 아닌 ID", "/students/abc", "abc"},
		{"ID가 없는 경우", "/students/", ""},
		{"잘못된 경로", "/students", ""},
		{"루트 경로", "/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIDFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
