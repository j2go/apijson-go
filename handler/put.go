package handler

import (
	"fmt"
	"github.com/j2go/apijson/db"
	"github.com/j2go/apijson/logger"
	"net/http"
	"strings"
)

// PutHandler 更新接口，暂时仅支持单条更新
func PutHandler(w http.ResponseWriter, r *http.Request) {
	commonHandle(w, r, putDataHandler)
}

func putDataHandler(bodyMap map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for k, v := range bodyMap {
		if _, exists := db.AllTable[k]; exists {
			if kvs, ok := v.(map[string]interface{}); ok {
				if id, err := updateOne(k, kvs); err != nil {
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
			logger.Warnf("PutHandler %s not exists", k)
			resultMap[k] = "table " + k + " not exists"
		}
	}
	return resultMap
}

func updateOne(table string, kvs map[string]interface{}) (int64, error) {
	if id, exists := kvs["id"]; exists {
		if _, ok := id.(float64); !ok {
			logger.Warnf("id: %+v", id)
			return -1, fmt.Errorf("'id' type is not num, key: %s， kvs: %v", table, kvs)
		}
		size := len(kvs) - 1
		fields := make([]string, size)
		args := make([]interface{}, size)
		i := 0
		for k, v := range kvs {
			if k != "id" {
				fields[i] = "`" + k + "`=?"
				args[i] = v
				i++
			}
		}
		sql := fmt.Sprintf("update %s set %s where id=%v", table, strings.Join(fields, ","), id)
		logger.Debugf("sql: %s, args: %v", sql, args)
		if err := db.Update(sql, args...); err != nil {
			return -2, err
		}
		return int64(id.(float64)), nil
	} else {
		return -100, fmt.Errorf("data update must have 'id' field, key: %s， kvs: %v", table, kvs)
	}
}
