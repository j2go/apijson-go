package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var db *sqlx.DB

func main() {
	http.HandleFunc("/get", GetHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

func init() {
	database, err := sqlx.Open("mysql", "apijson:1234qqqq@tcp(y.tadev.cn:53306)/sys")
	if err != nil {
		log.Fatal("db connect error", err)
	}
	db = database
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		log.Println("read request body error", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(data, &bodyMap); err != nil {
			log.Println("parse request body json error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		respMap := make(map[string]interface{})
		for table, fields := range bodyMap {
			if fields != nil {
				respMap[table] = QueryTable(table, fields)
			}
			log.Println("get:query table: ", table, ", fields: ", fields)
		}
		if respBody, err := json.Marshal(respMap); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(respBody)
		}
	}
}

func QueryTable(table string, fields interface{}) interface{} {
	var buffer bytes.Buffer
	buffer.WriteString("select * from ")
	buffer.WriteString(table)
	buffer.WriteString(" where ")
	if fieldMap, ok := fields.(map[string]interface{}); !ok {
		return "fields error, only support object."
	} else {
		size := len(fieldMap)
		cols := make([]string, size)
		values := make([]interface{}, size)
		i := 0
		for col, value := range fieldMap {
			if value == nil {
				return "field value error, " + col + " is nil"
			}
			cols[i] = col + "=?"
			values[i] = value
		}
		buffer.WriteString(strings.Join(cols, " and "))
		sql := buffer.String()
		if rows, err := db.Query(sql, values...); err != nil {
			return err.Error()
		} else {
			if rows.Next() {
				if columns, err := rows.Columns(); err != nil {
					return "get rows error: " + err.Error()
				} else {
					values := make([]interface{}, len(columns))
					for k := range columns {
						str := ""
						values[k] = &str
					}
					err = rows.Scan(values...)
					if err != nil {
						return "rows.Scan error: " + err.Error()
					}
					resultMap := make(map[string]interface{})
					for k, colName := range columns {
						resultMap[colName] = values[k]
					}
					return resultMap
				}
			} else {
				return ""
			}
		}
	}
}
