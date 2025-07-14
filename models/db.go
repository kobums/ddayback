package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	*sql.DB
}

type Where struct {
	Column  string
	Value   interface{}
	Compare string
}

type Custom struct {
	Query string
	Args  []interface{}
}

type Paging struct {
	Page     int
	PageSize int
}

type Ordering struct {
	OrderBy string
}

var DB *Connection

func InitDatabase() error {
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "dday")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	maxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "300"))

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))

	DB = &Connection{db}
	log.Println("Database connected successfully")

	return createTables()
}

func (c *Connection) Begin() (*sql.Tx, error) {
	return c.DB.Begin()
}

func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS ddays_tb (
		d_id VARCHAR(36) PRIMARY KEY,
		d_title VARCHAR(255) NOT NULL,
		d_target_date DATE NOT NULL,
		d_category VARCHAR(50) NOT NULL DEFAULT '개인',
		d_memo TEXT,
		d_is_important BOOLEAN DEFAULT FALSE,
		d_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		d_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
		INDEX idx_d_target_date (d_target_date),
		INDEX idx_d_category (d_category),
		INDEX idx_d_is_important (d_is_important),
		INDEX idx_d_created_at (d_created_at)
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Tables created/verified successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewPaging(page, pageSize int) Paging {
	return Paging{Page: page, PageSize: pageSize}
}

func NewOrdering(orderBy string) Ordering {
	return Ordering{OrderBy: orderBy}
}

func NewWhere(column string, value interface{}, compare string) Where {
	return Where{Column: column, Value: value, Compare: compare}
}

func NewCustom(query string, args ...interface{}) Custom {
	return Custom{Query: query, Args: args}
}