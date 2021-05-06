package main

import (
	"github.com/keepfoo/apijson/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/head", handler.HeadHandler)
	http.HandleFunc("/get", handler.GetHandler)
	http.HandleFunc("/post", handler.PostHandler)
	http.HandleFunc("/put", handler.PutHandler)
	http.HandleFunc("/delete", handler.DeleteHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
