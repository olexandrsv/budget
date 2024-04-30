package repository

import (
	"bytes"
	"errors"
	"fmt"
	"budget/pkg/models"
	"strings"
)

type Repository struct {
	db    *Database
	cache *Cache
}

func (repo *Repository) Cache() *Cache {
	return repo.cache
}

func NewRepository(db *Database, c *Cache) *Repository {
	return &Repository{
		db:    db,
		cache: c,
	}
}

func get[T models.DataObject[T]](repo *Repository) (models.Store[T], error) {
	store := cacheGetStore[T](repo.cache)
	if store != nil {
		return store, nil
	}

	var item T
	s, err := repo.db.Get(item.SelectQuery())
	if err != nil {
		return nil, err
	}

	store, err = parseResponse[T](s)
	if err != nil {
		return nil, err
	}

	cachePutStore(repo.cache, store)
	return store, nil
}

func parseResponse[T models.CacheObject[T]](s string) (models.Store[T], error) {
	if len(s) < 2 {
		return nil, errors.New("invalid server response")
	}
	s = s[1 : len(s)-1]

	store := models.NewStore[T]()
	var idx int
	for {
		if idx = strings.Index(s, "]"); idx == -1 {
			break
		}
		var item T
		item, err := item.Parse(s[:idx+1])
		if err != nil {
			return nil, err
		}
		store.Put(item.Index(), item)
		s = s[idx+1:]
	}

	return store, nil
}

func create[T models.DataObject[T]](repo *Repository, item T) (int, error) {
	id, err := repo.db.Insert(item.InsertQuery(), item.Table())
	if err != nil {
		return 0, err
	}
	item.SetIndex(id)
	cachePutItem(repo.cache, item)
	return id, nil
}

func update[T models.DataObject[T]](repo *Repository, item T) error {
	cachedItem := cacheGetItem[T](repo.cache, item.Index())
	item.SetUpdatedAt(cachedItem.UpdatedAt())

	err := repo.db.UpdateEx(item.UpdateQuery(), item.Table(), item.UpdatedAt(), item.Index())
	if err != nil {
		return err
	}

	cachePutItem(repo.cache, item)
	return nil
}

func remove[T models.DataObject[T]](repo *Repository, index int) error {
	var item T
	err := repo.db.Delete(item.DeleteQuery(index))
	if err != nil {
		return err
	}

	cacheDeleteItem[T](repo.cache, index)
	return nil
}

func (repo *Repository) CreateProduct(p *models.Product) (int, error) {
	return create(repo, p)
}

func (repo *Repository) GetProducts() (models.Store[*models.Product], error) {
	return get[*models.Product](repo)
}

func (repo *Repository) GetProduct(id int) *models.Product {
	return cacheGetItem[*models.Product](repo.cache, id)
}

func (repo *Repository) UpdateProduct(p *models.Product) error {
	return update(repo, p)
}

func (repo *Repository) DeleteProduct(id int) error {
	return remove[*models.Product](repo, id)
}

func (repo *Repository) CreateAuthor(a *models.Author) (int, error) {
	return create(repo, a)
}

func (repo *Repository) GetAuthors() (models.Store[*models.Author], error) {
	return get[*models.Author](repo)
}

func (repo *Repository) GetAuthor(id int) *models.Author {
	return cacheGetItem[*models.Author](repo.cache, id)
}

func (repo *Repository) UpdateAuthor(a *models.Author) error {
	return update(repo, a)
}

func (repo *Repository) DeleteAuthor(id int) error {
	return remove[*models.Author](repo, id)
}

func (repo *Repository) CreateShop(s *models.Shop) (int, error) {
	return create(repo, s)
}

func (repo *Repository) GetShops() (models.Store[*models.Shop], error) {
	return get[*models.Shop](repo)
}

func (repo *Repository) GetShop(id int) *models.Shop {
	return cacheGetItem[*models.Shop](repo.cache, id)
}

func (repo *Repository) UpdateShop(s *models.Shop) error {
	return update(repo, s)
}

func (repo *Repository) DeleteShop(id int) error {
	return remove[*models.Shop](repo, id)
}

func (repo *Repository) CreateAgent(a *models.Agent) (int, error) {
	return create(repo, a)
}

func (repo *Repository) GetAgents() (models.Store[*models.Agent], error) {
	return get[*models.Agent](repo)
}

func (repo *Repository) GetAgent(id int) *models.Agent {
	return cacheGetItem[*models.Agent](repo.cache, id)
}

func (repo *Repository) UpdateAgent(a *models.Agent) error {
	return update(repo, a)
}

func (repo *Repository) DeleteAgent(id int) error {
	return remove[*models.Agent](repo, id)
}

func (repo *Repository) GetDocNames() (models.Store[*models.DocName], error) {
	return get[*models.DocName](repo)
}

func (repo *Repository) GetDocName(id int) *models.DocName {
	return cacheGetItem[*models.DocName](repo.cache, id)
}

func (repo *Repository) CreateDocument(d *models.Document) (int, error) {
	return create(repo, d)
}

func (repo *Repository) GetDocuments() (models.Store[*models.Document], error) {
	return get[*models.Document](repo)
}

func (repo *Repository) GetDocument(id int) *models.Document {
	return cacheGetItem[*models.Document](repo.cache, id)
}

func (repo *Repository) UpdateDocument(d *models.Document) error {
	return update(repo, d)
}

func (repo *Repository) DeleteDocument(id int) error {
	return remove[*models.Document](repo, id)
}

func (repo *Repository) GetOperations(id int) (models.Store[*models.Operation], error) {
	var o *models.Operation
	s, err := repo.db.Get(o.SelectQuery(id))
	if err != nil {
		return nil, err
	}

	items, err := parseResponse[*models.Operation](s)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *Repository) UpdateOperations(docID int, operations models.Store[*models.Operation]) error {
	operations.ForEach(func(o *models.Operation) error {
		fmt.Printf("operation: %v\n", o)
		return nil
	})
	query := createQueryItems(operations)
	var o *models.Operation
	return repo.db.Update(o.UpdateQuery(docID, query))
}

func (repo *Repository) CreateOperations(operations models.Store[*models.Operation]) error {
	query := createQueryItems(operations)
	var o *models.Operation
	return repo.db.Update(o.InsertQuery(query))
}

func createQueryItems(operations models.Store[*models.Operation]) string {
	var b bytes.Buffer
	var i int
	operations.ForEach(func(o *models.Operation) error {
		b.WriteString(o.ToSqlRow())
		if i < operations.Len()-1 {
			b.WriteString(",")
		}
		i++
		return nil
	})

	return b.String()
}
