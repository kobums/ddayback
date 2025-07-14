-- D-Day 애플리케이션 데이터베이스 스키마
-- MariaDB/MySQL용

-- 데이터베이스 생성 (필요시)
CREATE DATABASE IF NOT EXISTS dday 
    CHARACTER SET utf8mb4 
    COLLATE utf8mb4_unicode_ci;

USE dday;

-- D-Day 테이블
CREATE TABLE IF NOT EXISTS ddays_tb (
    d_id VARCHAR(36) PRIMARY KEY,
    d_title VARCHAR(255) NOT NULL,
    d_target_date DATE NOT NULL,
    d_category VARCHAR(50) NOT NULL DEFAULT '개인',
    d_memo TEXT,
    d_is_important BOOLEAN DEFAULT FALSE,
    d_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    d_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 인덱스
    INDEX idx_d_target_date (d_target_date),
    INDEX idx_d_category (d_category),
    INDEX idx_d_is_important (d_is_important),
    INDEX idx_d_created_at (d_created_at)
);

-- 카테고리 테이블 (선택사항 - 향후 확장용)
CREATE TABLE IF NOT EXISTS categories_tb (
    cat_id INT AUTO_INCREMENT PRIMARY KEY,
    cat_name VARCHAR(50) NOT NULL UNIQUE,
    cat_color VARCHAR(7), -- hex color code
    cat_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 기본 카테고리 데이터 삽입
INSERT INTO categories_tb (cat_name, cat_color) VALUES 
    ('개인', '#007bff'),
    ('학업', '#28a745'),
    ('업무', '#ffc107'),
    ('기타', '#6c757d')
ON DUPLICATE KEY UPDATE cat_name = VALUES(cat_name);

-- 사용자 테이블 (향후 확장용)
CREATE TABLE IF NOT EXISTS users_tb (
    u_id VARCHAR(36) PRIMARY KEY,
    u_email VARCHAR(255) UNIQUE NOT NULL,
    u_name VARCHAR(100) NOT NULL,
    u_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    u_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- D-Day에 사용자 연결 (향후 확장용)
-- ALTER TABLE ddays_tb ADD COLUMN d_user_id VARCHAR(36);
-- ALTER TABLE ddays_tb ADD FOREIGN KEY (d_user_id) REFERENCES users_tb(u_id) ON DELETE CASCADE;

-- 알림 설정 테이블 (향후 확장용)
CREATE TABLE IF NOT EXISTS notifications_tb (
    n_id INT AUTO_INCREMENT PRIMARY KEY,
    n_dday_id VARCHAR(36) NOT NULL,
    n_days_before INT NOT NULL, -- 며칠 전에 알림
    n_is_active BOOLEAN DEFAULT TRUE,
    n_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (n_dday_id) REFERENCES ddays_tb(d_id) ON DELETE CASCADE,
    INDEX idx_n_dday_id (n_dday_id)
);

-- 샘플 데이터 (개발용)
INSERT INTO ddays_tb (d_id, d_title, d_target_date, d_category, d_memo, d_is_important) VALUES 
    (UUID(), '대학교 졸업', '2024-08-15', '학업', '졸업논문 제출 마감', TRUE),
    (UUID(), '여행 계획', '2024-07-20', '개인', '유럽 여행 준비', FALSE),
    (UUID(), '프로젝트 마감', '2024-07-30', '업무', '새로운 앱 출시', TRUE)
ON DUPLICATE KEY UPDATE d_title = VALUES(d_title);