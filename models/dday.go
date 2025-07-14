package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type DDay struct {
	ID          string    `json:"id" db:"d_id"`
	Title       string    `json:"title" db:"d_title"`
	TargetDate  string    `json:"target_date" db:"d_target_date"`
	Category    string    `json:"category" db:"d_category"`
	Memo        string    `json:"memo" db:"d_memo"`
	IsImportant bool      `json:"is_important" db:"d_is_important"`
	CreatedAt   time.Time `json:"created_at" db:"d_created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"d_updated_at"`
}

type DdayManager struct {
	Conn *Connection
}

func NewDdayManager() *DdayManager {
	return &DdayManager{Conn: DB}
}

func (m *DdayManager) GetAll(args ...interface{}) ([]DDay, error) {
	query := "SELECT d_id, d_title, d_target_date, d_category, d_memo, d_is_important, d_created_at, d_updated_at FROM ddays_tb"
	whereClause, orderClause, limitClause, queryArgs := m.buildQuery(args...)

	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	if orderClause != "" {
		query += " ORDER BY " + orderClause
	} else {
		query += " ORDER BY d_target_date ASC"
	}
	if limitClause != "" {
		query += " " + limitClause
	}

	rows, err := m.Conn.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ddays []DDay
	for rows.Next() {
		var dday DDay
		err := rows.Scan(&dday.ID, &dday.Title, &dday.TargetDate,
			&dday.Category, &dday.Memo, &dday.IsImportant,
			&dday.CreatedAt, &dday.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ddays = append(ddays, dday)
	}

	return ddays, nil
}

func (m *DdayManager) GetByID(id string) (*DDay, error) {
	var dday DDay
	query := "SELECT d_id, d_title, d_target_date, d_category, d_memo, d_is_important, d_created_at, d_updated_at FROM ddays_tb WHERE d_id = ?"

	err := m.Conn.QueryRow(query, id).Scan(&dday.ID, &dday.Title, &dday.TargetDate,
		&dday.Category, &dday.Memo, &dday.IsImportant,
		&dday.CreatedAt, &dday.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &dday, nil
}

func (m *DdayManager) Create(dday *DDay) error {
	query := `INSERT INTO ddays_tb (d_id, d_title, d_target_date, d_category, d_memo, d_is_important, d_created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.Conn.Exec(query, dday.ID, dday.Title, dday.TargetDate,
		dday.Category, dday.Memo, dday.IsImportant, dday.CreatedAt)
	return err
}

func (m *DdayManager) Update(id string, dday *DDay) error {
	query := `UPDATE ddays_tb SET d_title = ?, d_target_date = ?, d_category = ?, d_memo = ?, d_is_important = ? 
			  WHERE d_id = ?`

	_, err := m.Conn.Exec(query, dday.Title, dday.TargetDate, dday.Category,
		dday.Memo, dday.IsImportant, id)
	return err
}

func (m *DdayManager) Delete(id string) error {
	query := "DELETE FROM ddays_tb WHERE d_id = ?"
	_, err := m.Conn.Exec(query, id)
	return err
}

func (m *DdayManager) Count(args ...interface{}) (int, error) {
	query := "SELECT COUNT(*) FROM ddays_tb"
	whereClause, _, _, queryArgs := m.buildQuery(args...)

	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	var count int
	err := m.Conn.QueryRow(query, queryArgs...).Scan(&count)
	return count, err
}

func (m *DdayManager) CreateWithTx(tx *sql.Tx, dday *DDay) error {
	query := `INSERT INTO ddays_tb (d_id, d_title, d_target_date, d_category, d_memo, d_is_important, d_created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(query, dday.ID, dday.Title, dday.TargetDate,
		dday.Category, dday.Memo, dday.IsImportant, dday.CreatedAt)
	return err
}

func (m *DdayManager) UpdateWithTx(tx *sql.Tx, id string, dday *DDay) error {
	query := `UPDATE ddays_tb SET d_title = ?, d_target_date = ?, d_category = ?, d_memo = ?, d_is_important = ? 
			  WHERE d_id = ?`

	_, err := tx.Exec(query, dday.Title, dday.TargetDate, dday.Category,
		dday.Memo, dday.IsImportant, id)
	return err
}

func (m *DdayManager) DeleteWithTx(tx *sql.Tx, id string) error {
	query := "DELETE FROM ddays_tb WHERE d_id = ?"
	_, err := tx.Exec(query, id)
	return err
}

func (m *DdayManager) buildQuery(args ...interface{}) (string, string, string, []interface{}) {
	var whereConditions []string
	var queryArgs []interface{}
	var orderClause string
	var limitClause string

	for _, arg := range args {
		switch v := arg.(type) {
		case Where:
			whereConditions = append(whereConditions, fmt.Sprintf("%s %s ?", v.Column, v.Compare))
			queryArgs = append(queryArgs, v.Value)
		case Custom:
			whereConditions = append(whereConditions, v.Query)
			queryArgs = append(queryArgs, v.Args...)
		case Ordering:
			orderClause = v.OrderBy
		case Paging:
			if v.Page > 0 && v.PageSize > 0 {
				offset := (v.Page - 1) * v.PageSize
				limitClause = fmt.Sprintf("LIMIT %d, %d", offset, v.PageSize)
			}
		}
	}

	whereClause := strings.Join(whereConditions, " AND ")
	return whereClause, orderClause, limitClause, queryArgs
}