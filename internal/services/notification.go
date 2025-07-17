package services

import (
	"fmt"
	"log"
	"time"

	"attendance-management/internal/repositories"
)

type NotificationService interface {
	SendAttendanceMessage(studentID string, status string, t time.Time)
}

type NotificationServiceImpl struct {
	studentRepo   repositories.StudentRepository
	messageSender MessageSender
}

func NewNotificationService(studentRepo repositories.StudentRepository, messageSender MessageSender) NotificationService {
	return &NotificationServiceImpl{
		studentRepo:   studentRepo,
		messageSender: messageSender,
	}
}

func (n *NotificationServiceImpl) SendAttendanceMessage(studentID string, status string, t time.Time) {
	student, err := n.studentRepo.GetByID(studentID)
	if err != nil {
		log.Printf("학생 조회 오류: %v", err)
		return
	}
	message := fmt.Sprintf("[출결 알림]: %s 학생이 %s 하였습니다. %s 시간: %s",
		student.Name, status, status, t.Format("15:04"))
	go n.messageSender.Send(student.ParentPhone, message)
}
