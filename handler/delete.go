package handler

import (
	"fmt"
	"github.com/j2go/apijson/db"
	"github.com/j2go/apijson/logger"
	"net/http"
)

// DeleteHandler 根据 id 参数删除数据
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	commonHandle(w, r, delDataHandler)
}

func delDataHandler(bodyMap map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for k, v := range bodyMap {
		if _, exists := db.AllTable[k]; exists {
			if kvs, ok := v.(map[string]interface{}); ok {
				if n, err := doDelete(k, kvs); err != nil {
					resultMap[k] = map[string]interface{}{"code": http.StatusBadRequest, "message": err}
				} else {
					resultMap[k] = map[string]int64{"code": http.StatusOK, "count": n}
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

func doDelete(table string, kvs map[string]interface{}) (int64, error) {
	if idValue, exists := kvs["id"]; exists {
		if _, ok := idValue.(float64); !ok {
			logger.Warnf("id: %+v", idValue)
			return -1, fmt.Errorf("'id' type is not num, key: %s， kvs: %v", table, kvs)
		}
		sql := fmt.Sprintf("delete from %s where id=?", table)
		logger.Debugf("sql: %s", sql)
		if rows, err := db.Delete(sql, idValue); err != nil {
			return -2, err
		} else {
			return rows, nil
		}
	} else if idValue, exists = kvs["id{}"]; exists {
		if _, ok := idValue.([]interface{}); !ok {
			logger.Warnf("wrong id array: %+v", idValue)
			return -1, fmt.Errorf("'id{}' type is not num array, key: %s， kvs: %v", table, kvs)
		}
		idArr := idValue.([]interface{})
		sql := fmt.Sprintf("delete from %s where id in(%s)", table, genPlaceholder(len(idArr)))
		logger.Debugf("sql: %s", sql)
		if rows, err := db.Delete(sql, idArr...); err != nil {
			return -3, err
		} else {
			return rows, nil
		}
	} else {
		return -100, fmt.Errorf("data delete must have field 'id' or 'id{}' , key: %s， kvs: %v", table, kvs)
	}
}
