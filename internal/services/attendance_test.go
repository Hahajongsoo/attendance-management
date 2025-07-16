package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptr(t time.Time) *time.Time {
	return &t
}
func TestAttendanceService_CreateAttendance(t *testing.T) {
	type fields struct {
		classRepo      *mocks.ClassRepository
		attendanceRepo *mocks.AttendanceRepository
	}
	type args struct {
		studentID string
		now       time.Time
		classTime *time.Time
		existing  *models.Attendance
		updated   *models.Attendance
	}

	tests := []struct {
		name       string
		args       args
		wantError  bool
		setupMocks func(f fields, a args)
	}{
		{
			name: "등원 기록 x 등록한 수업 x",
			args: args{
				studentID: "1",
				now:       time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local),
				classTime: nil,
				existing:  nil,
				updated:   nil,
			},
			wantError: false,
			setupMocks: func(f fields, a args) {
				f.attendanceRepo.On("GetByStudentIDAndDate", a.studentID, "2025-07-14").
					Return(nil, nil)
				f.classRepo.On("GetClassesForStudentByWeekday", a.studentID, "월").
					Return([]models.Class{}, nil)
				f.attendanceRepo.On("Create", mock.MatchedBy(func(att *models.Attendance) bool {
					return att.Status == "출석"
				})).
					Return(nil)
			},
		},
		{
			name: "등원 기록 x 등록한 수업 o 지각",
			args: args{
				studentID: "1",
				now:       time.Date(2025, 7, 14, 10, 5, 0, 0, time.Local),
				classTime: ptr(time.Date(0, 1, 1, 10, 0, 0, 0, time.Local)),
				existing:  nil,
				updated:   nil,
			},
			wantError: false,
			setupMocks: func(f fields, a args) {
				f.attendanceRepo.On("GetByStudentIDAndDate", a.studentID, "2025-07-14").
					Return(nil, nil)
				f.classRepo.On("GetClassesForStudentByWeekday", a.studentID, "월").
					Return([]models.Class{
						{ClassID: 1, StartTime: models.TimeOnly{Time: *a.classTime}},
					}, nil)
				f.attendanceRepo.On("Create", mock.MatchedBy(func(att *models.Attendance) bool {
					return att.Status == "지각"
				})).
					Return(nil)
			},
		},
		{
			name: "등원 기록 x 등록한 수업 o 출석",
			args: args{
				studentID: "1",
				now:       time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local),
				classTime: ptr(time.Date(0, 1, 1, 10, 0, 0, 0, time.Local)),
				existing:  nil,
				updated:   nil,
			},
			wantError: false,
			setupMocks: func(f fields, a args) {
				f.attendanceRepo.On("GetByStudentIDAndDate", a.studentID, "2025-07-14").
					Return(nil, nil)
				f.classRepo.On("GetClassesForStudentByWeekday", a.studentID, "월").
					Return([]models.Class{
						{ClassID: 1, StartTime: models.TimeOnly{Time: *a.classTime}},
					}, nil)
				f.attendanceRepo.On("Create", mock.MatchedBy(func(att *models.Attendance) bool {
					return att.Status == "출석"
				})).
					Return(nil)
			},
		},
		{
			name: "등원 기록 o",
			args: args{
				studentID: "1",
				now:       time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local),
				classTime: nil,
				existing: &models.Attendance{
					StudentID: 1,
					Date:      models.DateOnly{Time: time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)},
					CheckIn:   models.TimeOnly{Time: time.Date(2025, 7, 14, 9, 0, 0, 0, time.Local)},
				},
				updated: &models.Attendance{
					StudentID: 1,
					Date:      models.DateOnly{Time: time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)},
					CheckIn:   models.TimeOnly{Time: time.Date(2025, 7, 14, 9, 0, 0, 0, time.Local)},
					CheckOut:  models.TimeOnly{Time: time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)},
				},
			},
			wantError: false,
			setupMocks: func(f fields, a args) {
				f.attendanceRepo.On("GetByStudentIDAndDate", a.studentID, "2025-07-14").
					Return(a.existing, nil)
				f.attendanceRepo.On("Update", a.studentID, "2025-07-14", a.updated).
					Return(int64(1), nil)
			},
		},
		{
			name: "등원 기록 o 하원 시간이 더 빠른 경우",
			args: args{
				studentID: "1",
				now:       time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local),
				classTime: nil,
				existing: &models.Attendance{
					StudentID: 1,
					Date:      models.DateOnly{Time: time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)},
					CheckIn:   models.TimeOnly{Time: time.Date(2025, 7, 14, 11, 0, 0, 0, time.Local)},
				},
			},
			wantError: true,
			setupMocks: func(f fields, a args) {
				f.attendanceRepo.On("GetByStudentIDAndDate", a.studentID, "2025-07-14").
					Return(a.existing, nil)
			},
		},
		{
			name: "하원기록 o",
			args: args{
				studentID: "1",
				now:       time.Date(2025, 7, 14, 10, 5, 0, 0, time.Local),
				classTime: nil,
				existing: &models.Attendance{
					StudentID: 1,
					Date:      models.DateOnly{Time: time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)},
					CheckIn:   models.TimeOnly{Time: time.Date(2025, 7, 14, 9, 0, 0, 0, time.Local)},
					CheckOut:  models.TimeOnly{Time: time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)},
				},
			},
			wantError: true,
			setupMocks: func(f fields, a args) {
				f.attendanceRepo.On("GetByStudentIDAndDate", a.studentID, "2025-07-14").
					Return(a.existing, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClassRepo := new(mocks.ClassRepository)
			mockAttendanceRepo := new(mocks.AttendanceRepository)
			tt.setupMocks(fields{mockClassRepo, mockAttendanceRepo}, tt.args)

			service := &AttendanceService{
				classRepo:      mockClassRepo,
				attendanceRepo: mockAttendanceRepo,
			}

			att := &models.Attendance{
				StudentID: 1,
				Date:      models.DateOnly{Time: tt.args.now},
				CheckIn:   models.TimeOnly{Time: tt.args.now},
			}

			err := service.CreateAttendance(att)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockClassRepo.AssertExpectations(t)
			mockAttendanceRepo.AssertExpectations(t)
		})
	}
}
