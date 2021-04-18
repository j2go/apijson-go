package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

var db *sqlx.DB

const dataSourceName = "apijson:1234qqqq@tcp(y.tadev.cn:53306)/sys"

func init() {
	var err error
	db, err = sqlx.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("db connect error", err)
	}
}

func QueryOne(sql string, args ...interface{}) (map[string]interface{}, error) {
	if !strings.Contains(strings.ToLower(sql), "limit") {
		sql += " limit 1"
	}
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colSize := len(columns)
	record := make(map[string]interface{})
	values := make([]interface{}, colSize)
	valuePointers := make([]interface{}, colSize)
	if rows.Next() {
		for i := 0; i < colSize; i++ {
			valuePointers[i] = &values[i]
		}
		rows.Scan(valuePointers...)
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			record[col] = v
		}
	}
	return record, nil
}

func QueryAll(sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0, 8)
	values := make([]interface{}, count)
	valuePointers := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePointers[i] = &values[i]
		}
		rows.Scan(valuePointers...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}
