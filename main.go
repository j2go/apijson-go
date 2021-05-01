package main

import (
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
	respMap := make(map[string]interface{})
	for key, fields := range bodyMap {
		if fields != nil {
			if strings.HasSuffix(key, "[]") {

			}
			if fieldMap, ok := fields.(map[string]interface{}); !ok {
				respMap["code"] = http.StatusBadRequest
				respMap[key] = fmt.Errorf("field type error, only support object")
			} else {
				parseObj := db.SQLParseObject{}
				if err := parseObj.From(key, fieldMap); err != nil {
					respMap["code"] = http.StatusBadRequest
					respMap[key] = err.Error()
				} else {
					if parseObj.QueryFirst {
						respMap[key], err = db.QueryOne(parseObj.ToSQL(), parseObj.Values...)
					} else {
						respMap[key], err = db.QueryAll(parseObj.ToSQL(), parseObj.Values...)
					}
					if err != nil {
						respMap[key] = err.Error()
					}
				}
			}
		}
		log.Println("get:query table: ", key, ", fields: ", fields)
	}
	return respMap
}
