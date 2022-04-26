package handler

import (
	"bytes"
	"encoding/json"
	"github.com/j2go/apijson/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

func commonHandle(w http.ResponseWriter, r *http.Request, bodyHandler func(map[string]interface{}) map[string]interface{}) {
	if r.Method == http.MethodOptions {
		//logger.Infof("%v", r.Header)
		cors(w, r)
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
		cors(w, r)
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

func cors(w http.ResponseWriter, r *http.Request) {
	allowHeaders := []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"}
	for k, _ := range r.Header {
		if len(k) > 0 {
			if notContains(allowHeaders, k) {
				allowHeaders = append(allowHeaders, k)
			}
		}
	}
	host := r.Header.Get("Origin")
	if len(host) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", host)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "http://apijson.cn")
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", ")) // 无效 "*")
	w.Header().Set("Access-Control-Request-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
}

func notContains(arr []string, k string) bool {
	for _, h := range arr {
		if h == k {
			return false
		}
	}
	return true
}

func genPlaceholder(n int) string {
	if n == 1 {
		return "?"
	} else {
		buf := bytes.Buffer{}
		buf.WriteString("?")
		for i := 1; i < n; i++ {
			buf.WriteString(",?")
		}
		return buf.String()
	}
}
