package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func TestGetJson(t *testing.T) {
	json, err := getJSON("show tables")
	if err != nil {
		t.Error(err)
	}
	log.Println(json)
}
