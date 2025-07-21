# 출결 관리 시스템 (Attendance Management System)

학원 출결 관리 프로그램입니다. 학생, 교사, 수업, 등록, 출결 정보를 관리할 수 있는 RESTful API를 제공합니다.

## 주요 기능

- **학생 관리**: 학생 등록, 조회, 수정, 삭제
- **교사 관리**: 교사 등록, 조회, 수정, 삭제
- **수업 관리**: 수업 등록, 조회, 수정, 삭제
- **수강신청 관리**: 수강신청 등록 및 관리
- **출결 관리**: 출결 등록, 조회, 수정, 삭제
- **SMS 알림**: Twilio를 통한 출결 상태 알림

## 기술 스택

- **언어**: Go 1.24.2
- **데이터베이스**: PostgreSQL
- **SMS 서비스**: Twilio
- **테스트**: testify, go-sqlmock

## API 문서

자세한 API 문서는 [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)를 참조하세요.

## 설치 및 실행

### 1. 의존성 설치
```bash
go mod download
```

### 2. 환경 변수 설정
`.env` 파일을 생성하고 다음 변수들을 설정하세요:
```
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIO_FROM_NUMBER=your_twilio_from_number
```

### 3. 데이터베이스 설정
PostgreSQL 데이터베이스를 설정하고 `migrations/` 폴더의 SQL 파일들을 순서대로 실행하세요.

### 4. 서버 실행
```bash
go run main.go
```

서버는 `http://localhost:8080`에서 실행됩니다.

## 테스트

```bash
# 전체 테스트 실행
go test ./...

# 커버리지 확인
go test -cover ./...
```

## 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다.
