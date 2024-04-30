package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Shop struct {
	ID        int
	Name      string
	Timestamp string
}

func (s *Shop) Title() string{
	return s.Name
}

func (s *Shop) ToRow() []string {
	return []string{
		strconv.Itoa(s.ID),
		s.Name,
		s.Timestamp,
	}
}

func (s *Shop) SelectQuery() string {
	return fmt.Sprintf("select id, name, updatedAt from %s", s.Table())
}

func (s *Shop) InsertQuery() string {
	return fmt.Sprintf("insert into %s (name, updatedAt) values (%q, %q)",
		s.Table(), s.Name, s.Timestamp)
}

func (s *Shop) UpdateQuery() string {
	return fmt.Sprintf("update %s set name=%q, updatedAt=%q where id=%d",
		s.Table(), s.Name, s.Timestamp, s.ID)
}

func (s *Shop) DeleteQuery(id int) string {
	return fmt.Sprintf("delete from %s where id=%d", s.Table(), id)
}

func (s *Shop) Table() string {
	return "shops"
}

func (s *Shop) Index() int {
	return s.ID
}

func (s *Shop) SetIndex(id int) {
	s.ID = id
}

func (s *Shop) UpdatedAt() string {
	return s.Timestamp
}

func (s *Shop) SetUpdatedAt(time string) {
	s.Timestamp = time
}

func (s *Shop) Parse(data string) (*Shop, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid data")
	}
	data = data[1 : len(data)-1]
	fields := strings.Split(data, ",")

	if len(fields) < 3 {
		return nil, errors.New("invalid data")
	}

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, errors.New("invalid id")
	}

	name := fields[1]
	name = strings.Trim(name, `"`)

	timestamp := fields[2]
	timestamp = strings.Trim(timestamp, `"`)

	return &Shop{id, name, timestamp}, nil
}

func (s *Shop) CacheKey() Key {
	return shopsKey
}
