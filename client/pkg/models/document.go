package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Document struct {
	ID         int
	Date       string
	Type       int
	Sum        float64
	AuthorID   int
	ShopID     int
	AgentID    int
	Timestamp  string
	GetData    func(DataRequest) DataResponse
	Operations Store[*Operation]
}

func (d *Document) ToRow() []string {
	req := DataRequest{
		DocNameID: d.Type,
		AuthorID:  d.AuthorID,
		ShopID:    d.ShopID,
		AgentID:   d.AgentID,
	}
	resp := d.GetData(req)
	shop, agent := "", ""
	if resp.Shop != nil {
		shop = resp.Shop.Name
	}
	if resp.Agent != nil {
		agent = resp.Agent.Name
	}
	return []string{
		strconv.Itoa(d.ID),
		d.Date,
		resp.DocName.Name,
		fmt.Sprintf("%v", d.Sum),
		resp.Author.Name,
		shop,
		agent,
	}
}

func (d *Document) Title() string {
	return ""
}

func (d *Document) SelectQuery() string {
	return fmt.Sprintf("select id, date, type, sum, authorId, shopId, agentId, updatedAt from %s", d.Table())
}

func (d *Document) InsertQuery() string {
	return fmt.Sprintf("insert into %s (date, type, sum, authorId, shopId, agentId, updatedAt) values (%q, %d, %v, %d, %d, %d, %q)",
		d.Table(), d.Date, d.Type, d.Sum, d.AuthorID, d.ShopID, d.AgentID, d.Timestamp)
}

func (d *Document) UpdateQuery() string {
	return fmt.Sprintf("update %s set date=%q, type=%d, sum=%v, authorId=%d, shopId=%d, agentId=%d, updatedAt=%q where id=%d",
		d.Table(), d.Date, d.Type, d.Sum, d.AuthorID, d.ShopID, d.AgentID, d.Timestamp, d.ID)
}

func (d *Document) DeleteQuery(id int) string {
	return fmt.Sprintf("delete from %s where id=%d", d.Table(), id)
}

func (d *Document) Table() string {
	return "documents"
}

func (d *Document) Index() int {
	return d.ID
}

func (d *Document) SetIndex(id int) {
	d.ID = id
}

func (d *Document) UpdatedAt() string {
	return d.Timestamp
}

func (d *Document) SetUpdatedAt(time string) {
	d.Timestamp = time
}

func (d *Document) Parse(s string) (*Document, error) {
	if len(s) < 2 {
		return nil, errors.New("invalid data")
	}
	s = s[1 : len(s)-1]
	fields := strings.Split(s, ",")

	if len(fields) < 8 {
		return nil, errors.New("invalid data")
	}

	c := NewConverter()

	newDoc := &Document{
		ID:        c.ToInt(fields[0]),
		Date:      c.ToString(fields[1]),
		Type:      c.ToInt(fields[2]),
		Sum:       c.ToFloat(fields[3]),
		AuthorID:  c.ToInt(fields[4]),
		ShopID:    c.ToInt(fields[5]),
		AgentID:   c.ToInt(fields[6]),
		Timestamp: c.ToString(fields[7]),
	}

	if err := c.Err(); err != nil {
		return nil, err
	}

	return newDoc, nil
}

func (d *Document) CacheKey() Key {
	return documentsKey
}
