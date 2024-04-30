package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Operation struct {
	ID        int
	ProdID    int
	Number    int
	Price     float64
	DocID     int
	Timestamp string
	GetData   func(DataRequest) DataResponse
}

func (o *Operation) ToRow() []string {
	documentReq := DataRequest{
		DocID: o.DocID,
	}
	document := o.GetData(documentReq).Document
	req := DataRequest{
		ProdID:   o.ProdID,
		AuthorID: document.AuthorID,
		ShopID:   document.ShopID,
	}
	resp := o.GetData(req)
	fmt.Println(resp)
	return []string{
		resp.Product.Name,
		strconv.Itoa(o.Number),
		fmt.Sprintf("%v", o.Price),
		fmt.Sprintf("%v", o.Price*float64(o.Number)),
		resp.Author.Name,
		resp.Shop.Name,
		document.Date,
	}
}

func (o *Operation) Index() int {
	return o.ID
}

func (o *Operation) SetIndex(idx int) {
	o.ID = idx
}

func (o *Operation) SelectQuery(id int) string {
	return fmt.Sprintf("select id, prodId, number, price, docId, updatedAt from %s where docId=%d", o.Table(), id)
}

func (o *Operation) InsertQuery(items string) string {
	return fmt.Sprintf("insert into %s (prodId, number, price, docId, updatedAt) values %s", o.Table(), items)
}

func (o *Operation) UpdateQuery(docID int, items string) string {
	return fmt.Sprintf("delete from %s where docId=%d; insert into %s (prodId, number, price, docId, updatedAt) values %s",
		o.Table(), docID, o.Table(), items)
}

func (o *Operation) Table() string {
	return "operations"
}

func (o *Operation) ToSqlRow() string {
	return fmt.Sprintf("(%d, %d, %v, %d, %q)", o.ProdID, o.Number, o.Price, o.DocID, o.Timestamp)
}

func (o *Operation) Parse(s string) (*Operation, error) {
	if len(s) < 2 {
		return nil, errors.New("invalid data")
	}
	s = s[1 : len(s)-1]
	fields := strings.Split(s, ",")

	if len(fields) < 6 {
		return nil, errors.New("invalid data")
	}

	c := NewConverter()

	operation := &Operation{
		ID:        c.ToInt(fields[0]),
		ProdID:    c.ToInt(fields[1]),
		Number:    c.ToInt(fields[2]),
		Price:     c.ToFloat(fields[3]),
		DocID:     c.ToInt(fields[4]),
		Timestamp: c.ToString(fields[5]),
	}

	if c.Err() != nil {
		return nil, c.Err()
	}

	return operation, nil
}

func (o *Operation) CacheKey() Key {
	return operationsKey
}
