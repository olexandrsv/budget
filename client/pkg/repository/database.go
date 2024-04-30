package repository

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Database struct{}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) Insert(sql, table string) (int, error) {
	sql = strings.ReplaceAll(sql, " ", "_")
	query := fmt.Sprintf("http://localhost:3333/insert?sql=%s&tbl=%s", sql, table)
	r, err := http.Get(query)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, errors.New(string(resp))
	}

	id, err := strconv.Atoi(string(resp))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *Database) Get(sql string) (string, error) {
	sql = strings.ReplaceAll(sql, " ", "_")
	query := fmt.Sprintf("http://localhost:3333/select?sql=%s", sql)
	r, err := http.Get(query)
	if err != nil {
		fmt.Println("http error: ", err)
		return "", err
	}
	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	if r.StatusCode == 200 {
		return string(resp), nil
	}

	return "", errors.New(string(resp))
}

func (db *Database) UpdateEx(sql, table, updatedAt string, id int) error {
	sql = strings.ReplaceAll(sql, " ", "_")
	updatedAt = strings.ReplaceAll(updatedAt, " ", "_")

	tmpl := "http://localhost:3333/update_ex?sql=%s&id=%d&updatedAt=%s&table=%s"
	query := fmt.Sprintf(tmpl, sql, id, updatedAt, table)
	r, err := http.Get(query)
	if err != nil {
		return err
	}

	if r.StatusCode == 200 {
		return nil
	}

	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return errors.New(string(resp))
}

func (db *Database) Update(sql string) error {
	fmt.Println("sql:", sql)
	sql = strings.ReplaceAll(sql, " ", "_")
	sql = strings.ReplaceAll(sql, ";", "|")

	tmpl := "http://localhost:3333/update?sql=%s"
	query := fmt.Sprintf(tmpl, sql)
	r, err := http.Get(query)
	if err != nil {
		return err
	}

	if r.StatusCode == 200 {
		return nil
	}

	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return errors.New(string(resp))
}

func (db *Database) Delete(sql string) error {
	sql = strings.ReplaceAll(sql, " ", "_")
	query := fmt.Sprintf("http://localhost:3333/delete?sql=%s", sql)

	r, err := http.Get(query)
	if err != nil {
		return err
	}

	if r.StatusCode == 200 {
		return nil
	}

	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return errors.New(string(resp))
}
