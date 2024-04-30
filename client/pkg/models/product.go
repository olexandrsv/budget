package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Product struct {
	ID        int
	Name      string
	Price     float64
	Countable bool
	Timestamp string
}

func (p *Product) Title() string{
	return p.Name
}

func (p *Product) ToRow() []string {
	countable := "ні"
	if p.Countable {
		countable = "так"
	}
	return []string{
		strconv.Itoa(p.ID),
		p.Name,
		fmt.Sprintf("%v", p.Price),
		countable,
	}
}

func (p *Product) SelectQuery() string {
	return "select id, name, price, countable, updatedAt from products"
}

func (p *Product) InsertQuery() string {
	return fmt.Sprintf("insert into products (name, price, countable, updatedAt) values (%q, %f, %t, %q)",
		p.Name, p.Price, p.Countable, p.Timestamp)
}

func (p *Product) UpdateQuery() string {
	return fmt.Sprintf("update products set name=%q, price=%v, countable=%t where id=%d",
		p.Name, p.Price, p.Countable, p.ID)
}

func (p *Product) DeleteQuery(id int) string {
	return fmt.Sprintf("delete from products where id=%d", id)
}

func (p *Product) Table() string {
	return "products"
}

func (p *Product) Index() int {
	return p.ID
}

func (p *Product) SetIndex(id int) {
	p.ID = id
}

func (p *Product) UpdatedAt() string {
	return p.Timestamp
}

func (p *Product) SetUpdatedAt(time string) {
	p.Timestamp = time
}

func (p *Product) Parse(s string) (*Product, error) {
	if len(s) < 2 {
		return nil, errors.New("invalid data")
	}
	s = s[1 : len(s)-1]
	fields := strings.Split(s, ",")

	if len(fields) < 5 {
		return nil, errors.New("invalid data")
	}

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, errors.New("invalid id")
	}

	name := fields[1]
	name = strings.Trim(name, `"`)

	price, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, errors.New("invalid price")
	}

	countable, err := strconv.ParseBool(fields[3])
	if err != nil {
		return nil, errors.New("invalid countable")
	}

	timestamp := fields[4]
	timestamp = strings.Trim(timestamp, `"`)

	return &Product{id, name, price, countable, timestamp}, nil
}

func (p *Product) CacheKey() Key {
	return productsKey
}
