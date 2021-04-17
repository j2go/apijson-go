package db

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Open("mysql", "apijson:1234qqqq@tcp(y.tadev.cn:53306)/sys")
	if err != nil {
		log.Fatal("db connect error", err)
	}
}

type TableShow struct {
	name string `db:"Table_in_sys"`
}

func getJSON(sqlString string) (string, error) {
	rows, err := db.Query(sqlString)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
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
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	log.Println(string(jsonData))
	return string(jsonData), nil
}
