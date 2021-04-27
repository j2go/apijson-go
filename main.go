package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keepfoo/apijson/db"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		handleRequestJson(data, w)
	}
}

func handleRequestJson(data []byte, w http.ResponseWriter) {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(data, &bodyMap); err != nil {
		log.Println("parse request body json error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respMap := make(map[string]interface{})
	for table, fields := range bodyMap {
		if fields != nil {
			if fieldMap, ok := fields.(map[string]interface{}); !ok {
				w.WriteHeader(http.StatusBadRequest)
				respMap[table] = fmt.Errorf("field type error, only support object")
			} else {
				parseObj := SQLParseObject{}
				if err := parseObj.from(table, fieldMap); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					respMap[table] = err.Error()
				} else {
					if parseObj.queryFirst {
						respMap[table], err = db.QueryOne(parseObj.toSQL(), parseObj.values...)
					} else {
						respMap[table], err = db.QueryAll(parseObj.toSQL(), parseObj.values...)
					}
					if err != nil {
						respMap[table] = err.Error()
					}
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

type SQLParseObject struct {
	columns    []string
	table      string
	where      []string
	limit      int
	page       int
	queryFirst bool
	withPage   bool
	values     []interface{}
}

func (o *SQLParseObject) from(table string, fieldMap map[string]interface{}) error {
	if strings.HasSuffix(table, "[]") {
		o.table = table[0 : len(table)-2]
		o.queryFirst = false
	} else {
		o.table = table
		o.queryFirst = true
	}
	for field, value := range fieldMap {
		if value == nil {
			return fmt.Errorf("field value error, %s is nil", field)
		} else if strings.HasPrefix(field, "@") {
			switch field[1:] {
			case "page":
				o.page = int(value.(float64))
			case "size":
				o.limit = int(value.(float64))
			case "column":
				o.columns = strings.Split(value.(string), ",")
			}
		} else {
			o.where = append(o.where, field+"=?")
			o.values = append(o.values, value)
		}
	}
	o.withPage = o.page > 0 && o.limit > 0
	return nil
}

func (o *SQLParseObject) toSQL() string {
	var buf bytes.Buffer
	buf.WriteString("select ")
	if o.columns == nil {
		buf.WriteString(" * ")
	} else {
		buf.WriteString(strings.Join(o.columns, ","))
	}
	buf.WriteString(" from ")
	buf.WriteString(o.table)
	buf.WriteString(" where ")
	buf.WriteString(strings.Join(o.where, " and "))
	if o.queryFirst {
		buf.WriteString(" limit 1")
	} else if o.withPage {
		buf.WriteString(" limit ")
		buf.WriteString(strconv.Itoa(o.limit))
		buf.WriteString(" offset ")
		buf.WriteString(strconv.Itoa(o.limit * (o.page - 1)))
	}
	return buf.String()
}
