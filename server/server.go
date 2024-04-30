package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	repo *Repository
}

func NewServer(repo *Repository) *Server {
	return &Server{
		repo: repo,
	}
}

func (s Server) Insert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	sql := query.Get("sql")
	tbl := query.Get("tbl")

	sql = strings.ReplaceAll(sql, "_", " ")

	id, err := s.repo.Insert(sql, tbl)
	if err != nil {
		s.writeResponse(w, 500, err.Error())
		return
	}
	s.writeResponse(w, 200, strconv.Itoa(id))
}

func (s Server) Select(w http.ResponseWriter, r *http.Request) {
	sql := r.URL.Query().Get("sql")
	sql = strings.ReplaceAll(sql, "_", " ")
	
	res, err := s.repo.Select(sql)
	if err != nil {
		s.writeResponse(w, 500, err.Error())
		return
	}
	s.writeResponse(w, 200, res)
}

func (s Server) UpdateEx(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	sql := query.Get("sql")
	paramId := query.Get("id")
	updatedAt := query.Get("updatedAt")
	table := query.Get("table")

	sql = strings.ReplaceAll(sql, "_", " ")
	updatedAt = strings.ReplaceAll(updatedAt, "_", " ")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Println(err)
		s.writeResponse(w, 400, err.Error())
		return
	}

	err = s.repo.UpdateEx(sql, id, updatedAt, table)
	if err != nil {
		fmt.Println(err)
		s.writeResponse(w, 500, err.Error())
		return
	}
	
	s.writeResponse(w, 200, "")
}

func (s Server) Update(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Println("query:", query)
	sql := query.Get("sql")
	sql = strings.ReplaceAll(sql, "_", " ")
	sql = strings.ReplaceAll(sql, "|", ";")
	fmt.Println("sql:", sql)

	err := s.repo.Update(sql)
	if err != nil {
		fmt.Println(err)
		s.writeResponse(w, 500, err.Error())
		return
	}
	s.writeResponse(w, 200, "")
}

func (s Server) Delete(w http.ResponseWriter, r *http.Request) {
	sql := r.URL.Query().Get("sql")
	sql = strings.ReplaceAll(sql, "_", " ")

	err := s.repo.Delete(sql)
	if err != nil {
		fmt.Println(err)
		s.writeResponse(w, 500, err.Error())
		return
	}
	s.writeResponse(w, 200, "")
}

func (s Server) writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
