package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Author struct {
	ID        int
	Name      string
	Card      string
	Timestamp string
}

func (a *Author) Title() string{
	return a.Name
}

func (a *Author) ToRow() []string {
	return []string{
		strconv.Itoa(a.ID),
		a.Name,
		a.Card,
	}
}

func (a *Author) SelectQuery() string {
	return fmt.Sprintf("select id, name, card, updatedAt from %s", a.Table())
}

func (a *Author) InsertQuery() string {
	return fmt.Sprintf("insert into %s (name, card, updatedAt) values (%q, %q, %q)",
		a.Table(), a.Name, a.Card, a.Timestamp)
}

func (a *Author) UpdateQuery() string {
	return fmt.Sprintf("update %s set name=%q, card=%q, updatedAt=%q where id=%d",
		a.Table(), a.Name, a.Card, a.Timestamp, a.ID)
}

func (a *Author) DeleteQuery(id int) string {
	return fmt.Sprintf("delete from %s where id=%d", a.Table(), id)
}

func (a *Author) Table() string {
	return "authors"
}

func (a *Author) SetIndex(id int) {
	a.ID = id
}

func (a *Author) UpdatedAt() string {
	return a.Timestamp
}

func (a *Author) SetUpdatedAt(time string) {
	a.Timestamp = time
}

func (a *Author) Index() int {
	return a.ID
}

func (a *Author) Parse(s string) (*Author, error) {
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

	card := fields[2]
	card = strings.Trim(card, `"`)

	timestamp := fields[3]
	timestamp = strings.Trim(timestamp, `"`)

	return &Author{
		ID:        id,
		Name:      name,
		Card:      card,
		Timestamp: timestamp,
	}, nil
}

func (a *Author) CacheKey() Key {
	return authorsKey
}
