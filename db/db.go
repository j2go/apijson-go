package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/j2go/apijson/logger"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

var db *sqlx.DB

const accessAliasSQL = "select name,alias from Access where alias is not null"

type TableMeta struct {
	Name    string
	Columns map[string]ColumnMeta
}

type ColumnMeta struct {
	Field   string         `db:"Field"`
	Type    string         `db:"Type"`
	Null    string         `db:"Null"`
	Key     string         `db:"Key"`
	Default sql.NullString `db:"Default"`
	Extra   sql.NullString `db:"Extra"`
}

type Access struct {
	Name  string `db:"name"`
	Alias string `db:"alias"`
}

var AllTable = make(map[string]TableMeta)

func Init(database string, dataSource string) {
	dataSourceName := dataSource + "/" + database
	showTableSQL := "select TABLE_NAME from information_schema.tables where table_schema='" + database + "' and table_type='BASE TABLE'"

	var err error
	db, err = sqlx.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("db connect error", err)
	}
	logger.Info("LoadTableMeta START")
	if rows, err := db.Query(showTableSQL); err != nil {
		log.Fatal("db Query error", err)
	} else {
		for rows.Next() {
			var name string
			rows.Scan(&name)
			AllTable[name] = TableMeta{Name: name, Columns: loadColumnMeta(name)}
		}
		logger.Infof("LoadTableMeta END, Table size: %d", len(AllTable))
	}
	if _, exists := AllTable["Access"]; exists {
		var accessList []Access
		if err = db.Select(&accessList, accessAliasSQL); err != nil {
			log.Fatal(err)
		}
		for _, a := range accessList {
			AllTable[a.Alias] = AllTable[a.Name]
			delete(AllTable, a.Name)
			logger.Infof("scan Access, alias %s -> %s", a.Name, a.Alias)
		}
	}
}

func loadColumnMeta(name string) map[string]ColumnMeta {
	var columns []ColumnMeta
	columnMap := make(map[string]ColumnMeta)
	err := db.Select(&columns, "desc "+name)
	if err != nil {
		return nil
	}
	keys := make([]string, len(columns))
	for i, c := range columns {
		keys[i] = c.Field
		columnMap[c.Field] = c
	}
	logger.Infof("LoadTableMeta %s [%s]", name, strings.Join(keys, ","))
	return columnMap
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

func Insert(sql string, args ...interface{}) (int64, error) {
	r, err := db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	var id int64
	id, err = r.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func Update(sql string, args ...interface{}) error {
	_, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func Delete(sql string, args ...interface{}) (int64, error) {
	if r, err := db.Exec(sql, args...); err != nil {
		return 0, err
	} else {
		var rows int64
		if rows, err = r.RowsAffected(); err != nil {
			return 0, err
		} else {
			return rows, nil
		}
	}
}
