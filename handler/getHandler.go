package handler

import (
	"encoding/json"
	"fmt"
	"github.com/keepfoo/apijson/db"
	"github.com/keepfoo/apijson/logger"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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
	respMap := NewSQLParseContext(bodyMap).getResponse()
	w.WriteHeader(respMap["code"].(int))
	if respBody, err := json.Marshal(respMap); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err := w.Write(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

type SQLParseContext struct {
	req           map[string]interface{}
	resp          map[string]interface{}
	waitKeys      map[string]bool
	completedKeys map[string]bool
	end           bool
}

func NewSQLParseContext(bodyMap map[string]interface{}) *SQLParseContext {
	logger.Debugf("NewSQLParseContext %v", bodyMap)
	return &SQLParseContext{req: bodyMap, resp: make(map[string]interface{}), waitKeys: make(map[string]bool), completedKeys: make(map[string]bool)}
}

func (c *SQLParseContext) getResponse() map[string]interface{} {
	for key, _ := range c.req {
		if !c.completedKeys[key] {
			c.parseResponse(key)
			if c.end {
				return c.resp
			}
		}
	}
	return c.resp
}

func (c *SQLParseContext) parseResponse(key string) {
	c.waitKeys[key] = true
	logger.Debugf("开始解析 %s", key)
	if c.req[key] == nil {
		c.End(http.StatusBadRequest, "值不能为空, key: "+key)
		return
	}
	if fieldMap, ok := c.req[key].(map[string]interface{}); !ok {
		c.End(http.StatusBadRequest, "值类型不对，只支持 Object 类型")
	} else {
		parseObj := db.SQLParseObject{LoadFunc: c.queryResp}
		if c.end {
			return
		}
		if err := parseObj.From(key, fieldMap); err != nil {
			c.End(http.StatusBadRequest, err.Error())
			return
		} else {
			if parseObj.QueryFirst {
				c.resp[key], err = db.QueryOne(parseObj.ToSQL(), parseObj.Values...)
			} else {
				c.resp[key], err = db.QueryAll(parseObj.ToSQL(), parseObj.Values...)
			}
			if err != nil {
				c.End(http.StatusInternalServerError, err.Error())
			} else {
				c.resp["code"] = http.StatusOK
			}
		}
	}
	c.waitKeys[key] = false
	//log.Println("get:query table: ", key, ", fields: ", fields)
}

func (c *SQLParseContext) queryResp(queryString string) interface{} {
	paths := strings.Split(queryString, "/")
	var targetValue interface{}
	for _, x := range paths {
		if targetValue == nil {
			if c.waitKeys[x] {
				c.End(http.StatusBadRequest, "关联查询有循环依赖，queryString: "+queryString)
				return nil
			} else if c.completedKeys[x] {
				targetValue = c.resp[x]
			} else {
				c.parseResponse(x)
				targetValue = c.resp[x]
			}
		} else {
			targetValue = targetValue.(map[string]interface{})[x]
		}
		if targetValue == nil {
			c.End(http.StatusBadRequest, fmt.Sprintf("关联查询未发现相应值，queryString: %s", queryString))
		}
	}
	return targetValue
}

func (c *SQLParseContext) End(code int, msg string) {
	c.resp["code"] = code
	c.resp["msg"] = msg
	c.end = true
	logger.Debugf("发生错误，终止请求，code: %d, msg: %s", code, msg)
}
