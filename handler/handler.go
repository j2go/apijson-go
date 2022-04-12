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
  host := r.Header.Get("Origin")
  headers := r.Header
  hs := []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"}
  if len(headers) > 0 {
    for k, _ := range headers {
      if len(k) > 0 {
        find := false
        for _, h := range hs {
          if h == k {
            find = true
            break
          }
        }

        if find {
          continue
        }

        hs = append(hs, k)
      }
    }
  }

  if len(host) > 0 {
    w.Header().Set("Access-Control-Allow-Origin", host)
  } else {
    w.Header().Set("Access-Control-Allow-Origin", "http://apijson.cn")
  }
  w.Header().Set("Access-Control-Allow-Credentials", "true")
  w.Header().Set("Access-Control-Allow-Headers", strings.Join(hs, ", ")) // 无效 "*")
  //w.Header().Set("Access-Control-Request-Method", "POST")
  w.Header().Set("Access-Control-Request-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
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
