package handler

import (
	"github.com/j2go/apijson/logger"
	"net/http"
)

//TODO
func HeadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(""))
	if err != nil {
		logger.Error(err.Error())
	}
}
