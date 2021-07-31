package main

import (
	"github.com/keepfoo/apijson/handler"
	"github.com/keepfoo/apijson/logger"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/head", handler.HeadHandler)
	http.HandleFunc("/get", handler.GetHandler)
	http.HandleFunc("/post", handler.PostHandler)
	http.HandleFunc("/put", handler.PutHandler)
	http.HandleFunc("/delete", handler.DeleteHandler)

	addr := ":8080"
	logger.SetLevel(logger.DEBUG)
	logger.Info("server listen on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
