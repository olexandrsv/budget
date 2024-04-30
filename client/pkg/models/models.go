package models

import (
	"strconv"
	"strings"
)

type Key int

const (
	authorsKey Key = iota
	agentsKey
	productsKey
	docNamesKey
	shopsKey
	documentsKey
	operationsKey
)

type DataObject[T any] interface {
	SqlObject[T]
	CacheObject[T]
}

type SqlObject[T any] interface {
	SelectQuery() string
	UpdateQuery() string
	InsertQuery() string
	DeleteQuery(int) string
	Table() string
	UpdatedAt() string
	SetUpdatedAt(string)
	Title() string
}

type CacheObject[T any] interface {
	CacheKey() Key
	Parse(string) (T, error)
	Object[T]
}

type Object[T any] interface {
	Index() int
	SetIndex(int)
}

type Converter struct {
	err error
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ToString(s string) string {
	if c.err != nil {
		return ""
	}
	return strings.Trim(s, `"`)
}

func (c *Converter) ToInt(s string) int {
	if c.err != nil {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		c.err = err
	}
	return n
}

func (c *Converter) ToFloat(s string) float64 {
	if c.err != nil {
		return 0
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		c.err = err
	}
	return n
}

func (c *Converter) Err() error {
	return c.err
}

type DataRequest struct {
	DocID     int
	DocNameID int
	AuthorID  int
	ShopID    int
	AgentID   int
	ProdID    int
}

type DataResponse struct {
	Document *Document
	DocName  *DocName
	Author   *Author
	Shop     *Shop
	Agent    *Agent
	Product  *Product
}
