package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Agent struct {
	ID        int
	Name      string
	Type      string
	Tel       string
	Timestamp string
}

func (a *Agent) Title() string{
	return a.Name
}

func (a *Agent) ToRow() []string {
	return []string{
		strconv.Itoa(a.ID),
		a.Name,
		a.Type,
		a.Tel,
		a.Timestamp,
	}
}

func (a *Agent) SelectQuery() string {
	return fmt.Sprintf("select id, name, type, tel, updatedAt from %s", a.Table())
}

func (a *Agent) InsertQuery() string {
	return fmt.Sprintf("insert into %s (name, type, tel, updatedAt) values (%q, %q, %q, %q)",
		a.Table(), a.Name, a.Type, a.Tel, a.Timestamp)
}

func (a *Agent) UpdateQuery() string {
	return fmt.Sprintf("update %s set name=%q, type=%q, tel=%q, updatedAt=%q where id=%d",
		a.Table(), a.Name, a.Type, a.Tel, a.Timestamp, a.ID)
}

func (a *Agent) DeleteQuery(id int) string {
	return fmt.Sprintf("delete from %s where id=%d", a.Table(), id)
}

func (a *Agent) Table() string {
	return "agents"
}

func (a *Agent) SetIndex(id int) {
	a.ID = id
}

func (a *Agent) UpdatedAt() string {
	return a.Timestamp
}

func (a *Agent) SetUpdatedAt(time string) {
	a.Timestamp = time
}

func (a *Agent) Index() int {
	return a.ID
}

func (a *Agent) Parse(s string) (*Agent, error) {
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

	t := fields[2]
	t = strings.Trim(t, `"`)

	tel := fields[3]
	tel = strings.Trim(tel, `"`)

	timestamp := fields[4]
	timestamp = strings.Trim(timestamp, `"`)

	return &Agent{
		ID:        id,
		Name:      name,
		Type:      t,
		Tel:       tel,
		Timestamp: timestamp,
	}, nil
}

func (a *Agent) CacheKey() Key {
	return agentsKey
}
