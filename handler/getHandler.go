package handler

import (
	"encoding/json"
	"github.com/keepfoo/apijson/db"
	"io/ioutil"
	"log"
	"net/http"
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
	respMap := getResponse(bodyMap)
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

func getResponse(bodyMap map[string]interface{}) map[string]interface{} {
	for key, fields := range bodyMap {
		if fields == nil {
			bodyMap["code"] = http.StatusBadRequest
			bodyMap["msg"] = "value cannot be nil, key: %s" + key
			return bodyMap
		}
		if fieldMap, ok := fields.(map[string]interface{}); !ok {
			bodyMap["code"] = http.StatusBadRequest
			bodyMap["msg"] = "field type error, only support object"
		} else {
			parseObj := db.SQLParseObject{Src: bodyMap}
			if err := parseObj.From(key, fieldMap); err != nil {
				bodyMap["code"] = http.StatusBadRequest
				bodyMap["msg"] = err.Error()
				return bodyMap
			} else {
				if parseObj.QueryFirst {
					bodyMap[key], err = db.QueryOne(parseObj.ToSQL(), parseObj.Values...)
				} else {
					bodyMap[key], err = db.QueryAll(parseObj.ToSQL(), parseObj.Values...)
				}
				if err != nil {
					bodyMap["msg"] = err.Error()
				}
			}
		}
		log.Println("get:query table: ", key, ", fields: ", fields)
	}
	return bodyMap
}
