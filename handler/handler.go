package handler

import (
	"encoding/json"
	"github.com/j2go/apijson/logger"
	"io/ioutil"
	"net/http"
)

func commonHandle(w http.ResponseWriter, r *http.Request, bodyHandler func(map[string]interface{}) map[string]interface{}) {
	if r.Method == http.MethodOptions {
		//logger.Infof("%v", r.Header)
		cors(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		logger.Error("请求参数有问题: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		logger.Infof("request: %s", string(data))
		var bodyMap map[string]interface{}
		if err = json.Unmarshal(data, &bodyMap); err != nil {
			logger.Error("请求体 JSON 格式有问题: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cors(w)
		dataMap := bodyHandler(bodyMap)
		var response []byte
		if response, err = json.Marshal(dataMap); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			//logger.Debugf("返回数据 %s", string(respBody))
			if _, err = w.Write(response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func cors(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "http://apijson.cn/")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Request-Method", "POST")
}
