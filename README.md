# D-Day Backend API

Go Fiber로 구현된 D-Day 관리 백엔드 API

## 설정

### 1. 데이터베이스 설정
MariaDB/MySQL 서버를 설치하고 실행하세요.

```bash
# MariaDB 설치 (macOS)
brew install mariadb
brew services start mariadb

# 데이터베이스 생성
mysql -u root -p
```

### 2. 스키마 생성
```bash
mysql -u root -p < schema.sql
```

### 3. 환경변수 설정
```bash
cp .env.example .env
# .env 파일을 편집하여 데이터베이스 정보 입력
```

### 4. 의존성 설치 및 실행
```bash
go mod tidy
go run main.go database.go
```

## API 엔드포인트

- `GET /` - API 정보
- `GET /health` - 서버 상태
- `GET /api/v1/ddays` - 모든 D-Day 조회
- `POST /api/v1/ddays` - D-Day 생성
- `GET /api/v1/ddays/:id` - 특정 D-Day 조회
- `PUT /api/v1/ddays/:id` - D-Day 수정
- `DELETE /api/v1/ddays/:id` - D-Day 삭제

## 데이터 구조

```json
{
  "id": "uuid",
  "title": "제목",
  "target_date": "2024-12-31",
  "category": "개인",
  "memo": "메모",
  "is_important": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```# ddayback
