package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type DocName struct {
	ID   int
	Name string
}

func (d *DocName) Title() string{
	return d.Name
}


func (d *DocName) ToRow() []string {
	return []string{
		strconv.Itoa(d.ID),
		d.Name,
	}
}

func (d *DocName) SelectQuery() string {
	return fmt.Sprintf("select id, name from %s", d.Table())
}

func (d *DocName) InsertQuery() string {
	return fmt.Sprintf("insert into %s (name) values (%q)", d.Table(), d.Name)
}

func (d *DocName) UpdateQuery() string {
	return fmt.Sprintf("update %s set name=%q where id=%d", d.Table(), d.Name, d.ID)
}

func (d *DocName) DeleteQuery(id int) string {
	return fmt.Sprintf("delete from %s where id=%d", d.Table(), id)
}

func (d *DocName) Table() string {
	return "docNames"
}

func (d *DocName) SetIndex(id int) {
	d.ID = id
}

func (d *DocName) UpdatedAt() string {
	return ""
}

func (d *DocName) SetUpdatedAt(time string) {}

func (d *DocName) Index() int {
	return d.ID
}

func (d *DocName) Parse(s string) (*DocName, error) {
	if len(s) < 2 {
		return nil, errors.New("invalid data")
	}
	s = s[1 : len(s)-1]
	fields := strings.Split(s, ",")

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}

	name := fields[1]
	name = strings.Trim(name, `"`)

	return &DocName{
		ID:        id,
		Name:      name,
	}, nil
}

func (d *DocName) CacheKey() Key {
	return docNamesKey
}
