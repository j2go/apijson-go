package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keepfoo/apijson/db"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/get", GetHandler)
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		log.Println("read request body error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		var bodyMap map[string]interface{}
		if err = json.Unmarshal(data, &bodyMap); err != nil {
			log.Println("parse request body json error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		respMap := make(map[string]interface{})
		for table, fields := range bodyMap {
			if fields != nil {
				var records []map[string]interface{}
				if records, err = QueryTable(table, fields); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					respMap[table] = err.Error()
				} else {
					if returnArray(table) {
						respMap[table] = records
					} else {
						respMap[table] = records[0]
					}
				}
			}
			log.Println("get:query table: ", table, ", fields: ", fields)
		}
		if respBody, err := json.Marshal(respMap); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(respBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func QueryTable(table string, fields interface{}) ([]map[string]interface{}, error) {
	var buffer bytes.Buffer
	buffer.WriteString("select * from ")
	if returnArray(table) {
		buffer.WriteString(table[0 : len(table)-2])
	} else {
		buffer.WriteString(table)
	}
	buffer.WriteString(" where ")
	if fieldMap, ok := fields.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("field type error, only support object")
	} else {
		size := len(fieldMap)
		cols := make([]string, size)
		values := make([]interface{}, size)
		i := 0
		for col, value := range fieldMap {
			if value == nil {
				return nil, fmt.Errorf("field value error, %s is nil", col)
			}
			cols[i] = col + "=?"
			values[i] = value
		}
		buffer.WriteString(strings.Join(cols, " and "))
		if returnArray(table) {
			return db.QueryAll(buffer.String(), values...)
		}
		buffer.WriteString(" limit 1")
		return db.QueryAll(buffer.String(), values...)
	}
}

func returnArray(table string) bool {
	return strings.HasSuffix(table, "[]")
}
