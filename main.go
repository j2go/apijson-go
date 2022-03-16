package main

import (
	"github.com/j2go/apijson/db"
	"github.com/j2go/apijson/handler"
	"github.com/j2go/apijson/logger"
	"log"
	"net/http"
)

func main() {
	db.Init("apijson", "root:1234qwer@tcp(localhost:3306)")
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
