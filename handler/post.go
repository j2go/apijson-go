package handler

import (
	"fmt"
	"github.com/j2go/apijson/db"
	"github.com/j2go/apijson/logger"
	"net/http"
	"strings"
)

// PostHandler 暂时仅支持单条数据插入
func PostHandler(w http.ResponseWriter, r *http.Request) {
	commonHandle(w, r, postDataHandler)
}

func postDataHandler(bodyMap map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for k, v := range bodyMap {
		if _, exists := db.AllTable[k]; exists {
			if kvs, ok := v.(map[string]interface{}); ok {
				if id, err := insertOne(k, kvs); err != nil {
					resultMap["code"] = http.StatusBadRequest
					resultMap["message"] = err.Error()
					return resultMap
				} else {
					data, _ := db.QueryOne("select * from "+k+" where id=?", id)
					resultMap[k] = data
				}
			} else {
				resultMap["code"] = http.StatusBadRequest
				resultMap["message"] = fmt.Sprintf("参数格式错误，key: %s, value: %v", k, v)
				return resultMap
			}
		} else {
			logger.Warnf("PostHandler %s not exists", k)
		}
	}
	return resultMap
}

func insertOne(table string, kvs map[string]interface{}) (int64, error) {
	size := len(kvs)
	keys := make([]string, size)
	values := make([]string, size)
	args := make([]interface{}, size)
	i := 0
	for field, value := range kvs {
		keys[i] = field
		values[i] = "?"
		args[i] = value
		i++
	}
	sql := "insert into " + table + "(" + strings.Join(keys, ",") + ") values(" + strings.Join(values, ",") + ")"
	logger.Debugf("sql: %s, args: %v", sql, args)
	return db.Insert(sql, args...)
}
