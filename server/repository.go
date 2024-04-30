package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	db, err := sql.Open("sqlite3", "budget.db")
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{db}
}

func (r *Repository) Select(sql string) (string, error) {
	res, err := r.db.Query(sql)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("[")
	for res.Next() {
		sb.WriteString("[")
		columns, err := res.Columns()
		if err != nil {
			return "", err
		}
		n := len(columns)
		types, err := res.ColumnTypes()
		if err != nil {
			return "", err
		}

		data := make([]interface{}, n)
		for i := 0; i < n; i++ {
			var v any
			data[i] = &v
		}
		err = res.Scan(data...)
		if err != nil {
			return "", err
		}

		for i := 0; i < n; i++ {
			info := *(data[i].(*any))
			s := fmt.Sprint(info)

			if types[i].ScanType().String() == "sql.NullString" {
				s = `"` + s + `"`
			}
			sb.WriteString(s)
			if i != n-1 {
				sb.WriteString(",")
			}
		}

		sb.WriteString("]")
	}
	sb.WriteString("]")

	return sb.String(), nil
}

func (r *Repository) Insert(sql string, table string) (int, error) {
	fmt.Println(sql)
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(sql)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	row := tx.QueryRow(fmt.Sprintf("select max(id) from %s", table))
	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	tx.Commit()

	return id, nil
}

func (r *Repository) UpdateEx(sql string, id int, updatedAt, table string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	tx.Exec(sql)
	row := tx.QueryRow(fmt.Sprintf("select count(id) from %s where id=$1 and updatedAt=$2", table), id, updatedAt)

	var count int
	err = row.Scan(&count)

	if err != nil {
		tx.Rollback()
		return err
	}

	if count == 0 {
		tx.Rollback()
		return errors.New("row was modified")
	}
	tx.Commit()

	return nil
}

func (r *Repository) Update(sql string) error {
	_, err := r.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(sql string) error {
	_, err := r.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
