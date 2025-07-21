# 출결 관리 시스템 API 문서

## 개요

이 API는 학원 출결 관리 시스템을 위한 RESTful API입니다. 학생, 교사, 수업, 등록, 출결 정보를 관리할 수 있습니다.

## 기본 정보

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **인코딩**: UTF-8

## 데이터 모델

### Student (학생)
```json
{
  "student_id": 1001,
  "name": "홍길동",
  "grade": "초등 3학년",
  "phone": "010-1234-5678",
  "parent_phone": "010-9876-5432"
}
```

### Teacher (교사)
```json
{
  "teacher_id": "T001",
  "password": "hashed_password",
  "name": "김선생님",
  "phone": "010-1111-2222"
}
```

### Class (수업)
```json
{
  "class_id": 1,
  "class_name": "수학 기초반",
  "days": "월,수,금",
  "start_time": "14:00",
  "end_time": "15:30",
  "price": 150000,
  "teacher_id": "T001"
}
```

### Enrollment (수강신청)
```json
{
  "enrollment_id": 1,
  "student_id": 1001,
  "class_id": 1,
  "enrolled_date": "2024-01-15"
}
```

### Attendance (출결)
```json
{
  "student_id": 1001,
  "date": "2024-01-15",
  "check_in": "13:55",
  "check_out": "15:35",
  "status": "출석"
}
```

### Payment (결제)
```json
{
  "payment_id": 1,
  "student_id": 1001,
  "class_id": 1,
  "payment_date": "2024-01-15",
  "amount": 150000,
  "enrollment_id": 1
}
```

## API 엔드포인트

### 1. 학생 관리 API

#### 1.1 전체 학생 목록 조회
- **URL**: `GET /students`
- **설명**: 등록된 모든 학생의 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "student_id": 1001,
    "name": "홍길동",
    "grade": "초등 3학년",
    "phone": "010-1234-5678",
    "parent_phone": "010-9876-5432"
  }
]
```

#### 1.2 학생 등록
- **URL**: `POST /students`
- **설명**: 새로운 학생을 등록합니다.
- **요청 본문**:
```json
{
  "student_id": 1002,
  "name": "김철수",
  "grade": "초등 4학년",
  "phone": "010-2345-6789",
  "parent_phone": "010-8765-4321"
}
```
- **응답**: `201 Created`
```json
{
  "message": "Student created"
}
```

#### 1.3 특정 학생 조회
- **URL**: `GET /students/{student_id}`
- **설명**: 특정 학생의 정보를 조회합니다.
- **응답 예시**:
```json
{
  "student_id": 1001,
  "name": "홍길동",
  "grade": "초등 3학년",
  "phone": "010-1234-5678",
  "parent_phone": "010-9876-5432"
}
```

#### 1.4 학생 정보 수정
- **URL**: `PUT /students/{student_id}`
- **설명**: 특정 학생의 정보를 수정합니다.
- **요청 본문**:
```json
{
  "student_id": 1001,
  "name": "홍길동",
  "grade": "초등 4학년",
  "phone": "010-1234-5678",
  "parent_phone": "010-9876-5432"
}
```
- **응답**: `200 OK`
```json
{
  "message": "Student updated"
}
```

#### 1.5 학생 삭제
- **URL**: `DELETE /students/{student_id}`
- **설명**: 특정 학생을 삭제합니다.
- **응답**: `200 OK`
```json
{
  "message": "Student deleted"
}
```

### 2. 교사 관리 API

#### 2.1 전체 교사 목록 조회
- **URL**: `GET /teachers`
- **설명**: 등록된 모든 교사의 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "teacher_id": "T001",
    "name": "김선생님",
    "phone": "010-1111-2222"
  }
]
```

#### 2.2 교사 등록
- **URL**: `POST /teachers`
- **설명**: 새로운 교사를 등록합니다.
- **요청 본문**:
```json
{
  "teacher_id": "T002",
  "password": "hashed_password",
  "name": "이선생님",
  "phone": "010-3333-4444"
}
```
- **응답**: `201 Created`
```json
{
  "message": "Teacher created"
}
```

#### 2.3 특정 교사 조회
- **URL**: `GET /teachers/{teacher_id}`
- **설명**: 특정 교사의 정보를 조회합니다.
- **응답 예시**:
```json
{
  "teacher_id": "T001",
  "name": "김선생님",
  "phone": "010-1111-2222"
}
```

#### 2.4 교사 정보 수정
- **URL**: `PUT /teachers/{teacher_id}`
- **설명**: 특정 교사의 정보를 수정합니다.
- **요청 본문**:
```json
{
  "teacher_id": "T001",
  "password": "new_hashed_password",
  "name": "김선생님",
  "phone": "010-1111-3333"
}
```
- **응답**: `200 OK`
```json
{
  "message": "Teacher updated"
}
```

#### 2.5 교사 삭제
- **URL**: `DELETE /teachers/{teacher_id}`
- **설명**: 특정 교사를 삭제합니다.
- **응답**: `200 OK`
```json
{
  "message": "Teacher deleted"
}
```

### 3. 수업 관리 API

#### 3.1 전체 수업 목록 조회
- **URL**: `GET /classes`
- **설명**: 등록된 모든 수업의 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "class_id": 1,
    "class_name": "수학 기초반",
    "days": "월,수,금",
    "start_time": "14:00",
    "end_time": "15:30",
    "price": 150000,
    "teacher_id": "T001"
  }
]
```

#### 3.2 수업 등록
- **URL**: `POST /classes`
- **설명**: 새로운 수업을 등록합니다.
- **요청 본문**:
```json
{
  "class_id": 2,
  "class_name": "영어 중급반",
  "days": "화,목",
  "start_time": "16:00",
  "end_time": "17:30",
  "price": 180000,
  "teacher_id": "T002"
}
```
- **응답**: `201 Created`
```json
{
  "message": "Class created"
}
```

#### 3.3 특정 수업 조회
- **URL**: `GET /classes/{class_id}`
- **설명**: 특정 수업의 정보를 조회합니다.
- **응답 예시**:
```json
{
  "class_id": 1,
  "class_name": "수학 기초반",
  "days": "월,수,금",
  "start_time": "14:00",
  "end_time": "15:30",
  "price": 150000,
  "teacher_id": "T001"
}
```

#### 3.4 수업 정보 수정
- **URL**: `PUT /classes/{class_id}`
- **설명**: 특정 수업의 정보를 수정합니다.
- **요청 본문**:
```json
{
  "class_id": 1,
  "class_name": "수학 기초반",
  "days": "월,수,금",
  "start_time": "14:00",
  "end_time": "16:00",
  "price": 160000,
  "teacher_id": "T001"
}
```
- **응답**: `200 OK`
```json
{
  "message": "Class updated"
}
```

#### 3.5 수업 삭제
- **URL**: `DELETE /classes/{class_id}`
- **설명**: 특정 수업을 삭제합니다.
- **응답**: `200 OK`
```json
{
  "message": "Class deleted"
}
```

#### 3.6 교사별 수업 목록 조회
- **URL**: `GET /teachers/{teacher_id}/classes`
- **설명**: 특정 교사가 담당하는 수업 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "class_id": 1,
    "class_name": "수학 기초반",
    "days": "월,수,금",
    "start_time": "14:00",
    "end_time": "15:30",
    "price": 150000,
    "teacher_id": "T001"
  }
]
```

### 4. 수강신청 관리 API

#### 4.1 전체 수강신청 목록 조회
- **URL**: `GET /enrollments`
- **설명**: 모든 수강신청 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "enrollment_id": 1,
    "student_id": 1001,
    "class_id": 1,
    "enrolled_date": "2024-01-15"
  }
]
```

#### 4.2 수강신청 등록
- **URL**: `POST /enrollments`
- **설명**: 새로운 수강신청을 등록합니다.
- **요청 본문**:
```json
{
  "student_id": 1002,
  "class_id": 1,
  "enrolled_date": "2024-01-16"
}
```
- **응답**: `201 Created`
```json
{
  "message": "Enrollment created"
}
```

#### 4.3 특정 수강신청 조회
- **URL**: `GET /enrollments/{enrollment_id}`
- **설명**: 특정 수강신청의 정보를 조회합니다.
- **응답 예시**:
```json
{
  "enrollment_id": 1,
  "student_id": 1001,
  "class_id": 1,
  "enrolled_date": "2024-01-15"
}
```

#### 4.4 수강신청 정보 수정
- **URL**: `PUT /enrollments/{enrollment_id}`
- **설명**: 특정 수강신청의 정보를 수정합니다.
- **요청 본문**:
```json
{
  "student_id": 1001,
  "class_id": 2,
  "enrolled_date": "2024-01-15"
}
```
- **응답**: `200 OK`
```json
{
  "message": "Enrollment updated"
}
```

#### 4.5 수강신청 삭제
- **URL**: `DELETE /enrollments/{enrollment_id}`
- **설명**: 특정 수강신청을 삭제합니다.
- **응답**: `200 OK`
```json
{
  "message": "Enrollment deleted"
}
```

#### 4.6 학생별 수강신청 목록 조회
- **URL**: `GET /students/{student_id}/enrollments`
- **설명**: 특정 학생의 수강신청 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "enrollment_id": 1,
    "student_id": 1001,
    "class_id": 1,
    "enrolled_date": "2024-01-15"
  }
]
```

#### 4.7 수업별 수강신청 목록 조회
- **URL**: `GET /classes/{class_id}/enrollments`
- **설명**: 특정 수업의 수강신청 목록을 조회합니다.
- **응답 예시**:
```json
[
  {
    "enrollment_id": 1,
    "student_id": 1001,
    "class_id": 1,
    "enrolled_date": "2024-01-15"
  }
]
```

### 5. 출결 관리 API

#### 5.1 학생별 특정 날짜 출결 조회
- **URL**: `GET /students/{student_id}/attendance/{date}`
- **설명**: 특정 학생의 특정 날짜 출결 정보를 조회합니다.
- **응답 예시**:
```json
{
  "student_id": 1001,
  "date": "2024-01-15",
  "check_in": "13:55",
  "check_out": "15:35",
  "status": "출석"
}
```

#### 5.2 출결 등록
- **URL**: `POST /students/{student_id}/attendance/{date}`
- **설명**: 특정 학생의 특정 날짜 출결을 등록합니다.
- **요청 본문**:
```json
{
  "student_id": 1001,
  "date": "2024-01-15",
  "check_in": "13:55",
  "check_out": "15:35",
  "status": "출석"
}
```
- **응답**: `201 Created`
```json
{
  "message": "Attendance created"
}
```

#### 5.3 출결 수정
- **URL**: `PUT /students/{student_id}/attendance/{date}`
- **설명**: 특정 학생의 특정 날짜 출결을 수정합니다.
- **요청 본문**:
```json
{
  "student_id": 1001,
  "date": "2024-01-15",
  "check_in": "14:00",
  "check_out": "15:30",
  "status": "지각"
}
```
- **응답**: `200 OK`
```json
{
  "message": "Attendance updated"
}
```

#### 5.4 출결 삭제
- **URL**: `DELETE /students/{student_id}/attendance/{date}`
- **설명**: 특정 학생의 특정 날짜 출결을 삭제합니다.
- **응답**: `200 OK`
```json
{
  "message": "Attendance deleted"
}
```

#### 5.5 날짜별 전체 출결 조회
- **URL**: `GET /attendance/{date}`
- **설명**: 특정 날짜의 모든 학생 출결 정보를 조회합니다.
- **응답 예시**:
```json
[
  {
    "student_id": 1001,
    "date": "2024-01-15",
    "check_in": "13:55",
    "check_out": "15:35",
    "status": "출석"
  },
  {
    "student_id": 1002,
    "date": "2024-01-15",
    "check_in": "14:10",
    "check_out": "15:40",
    "status": "지각"
  }
]
```

## 상태 코드

- `200 OK`: 요청이 성공적으로 처리됨
- `201 Created`: 리소스가 성공적으로 생성됨
- `400 Bad Request`: 잘못된 요청
- `404 Not Found`: 요청한 리소스를 찾을 수 없음
- `405 Method Not Allowed`: 허용되지 않은 HTTP 메서드
- `500 Internal Server Error`: 서버 내부 오류

## 오류 응답

모든 오류 응답은 다음과 같은 형식을 따릅니다:

```json
{
  "error": "오류 메시지"
}
```

## 데이터베이스 스키마

### students 테이블
- `student_id` (int, PK): 학생 ID (1000-99999)
- `name` (varchar(20)): 학생 이름
- `grade` (varchar(10)): 학년
- `phone` (varchar(15)): 학생 전화번호
- `parent_phone` (varchar(15)): 부모님 전화번호

### teachers 테이블
- `teacher_id` (varchar(30), PK): 교사 ID
- `password` (varchar(100)): 비밀번호
- `name` (varchar(30)): 교사 이름
- `phone_number` (varchar(20)): 교사 전화번호

### classes 테이블
- `class_id` (int, PK): 수업 ID
- `class_name` (varchar(50)): 수업명
- `days` (varchar(20)): 수업 요일
- `start_time` (time): 시작 시간
- `end_time` (time): 종료 시간
- `price` (int): 수업료
- `teacher_id` (varchar(30), FK): 담당 교사 ID

### enrollments 테이블
- `enrollment_id` (serial, PK): 수강신청 ID
- `student_id` (int, FK): 학생 ID
- `class_id` (int, FK): 수업 ID
- `enrolled_date` (date): 등록 날짜

### attendance 테이블
- `student_id` (int, FK): 학생 ID
- `date` (date): 출결 날짜
- `check_in` (time): 등원 시간
- `check_out` (time): 하원 시간
- `status` (text): 출결 상태 ('출석', '결석', '지각')

### payments 테이블
- `payment_id` (serial, PK): 결제 ID
- `student_id` (int, FK): 학생 ID
- `class_id` (int, FK): 수업 ID
- `payment_date` (date): 결제 날짜
- `amount` (int): 결제 금액
- `enrollment_id` (int, FK): 수강신청 ID