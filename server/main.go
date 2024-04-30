package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	repo := NewRepository()
	server := NewServer(repo)

	http.HandleFunc("/insert", server.Insert)
	http.HandleFunc("/select", server.Select)
	http.HandleFunc("/update_ex", server.UpdateEx)
	http.HandleFunc("/update", server.Update)
	http.HandleFunc("/delete", server.Delete)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println(err)
	}
}
