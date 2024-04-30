package service

import (
	"errors"
	"fmt"
	"budget/pkg/models"
	"budget/pkg/repository"
	"strconv"
	"strings"
	"time"
)

const (
	Reception = 1
	Lending   = 2
	Aquisiton = 3
)

type Service struct {
	repo *repository.Repository
}

func New(r *repository.Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateProduct(p *models.Product) (int, error) {
	p.Timestamp = time.Now().String()
	return s.repo.CreateProduct(p)
}

func (s *Service) GetProducts() (models.Store[*models.Product], error) {
	return s.repo.GetProducts()
}

func (s *Service) GetProduct(id int) *models.Product {
	return s.repo.GetProduct(id)
}

func (s *Service) UpdateProduct(p *models.Product) error {
	return s.repo.UpdateProduct(p)
}

func (s *Service) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}

func (s *Service) CreateAuthor(a *models.Author) (int, error) {
	a.Timestamp = time.Now().String()
	return s.repo.CreateAuthor(a)
}

func (s *Service) GetAuthors() (models.Store[*models.Author], error) {
	return s.repo.GetAuthors()
}

func (s *Service) GetAuthor(id int) *models.Author {
	return s.repo.GetAuthor(id)
}

func (s *Service) UpdateAuthor(a *models.Author) error {
	return s.repo.UpdateAuthor(a)
}

func (s *Service) DeleteAuthor(id int) error {
	return s.repo.DeleteAuthor(id)
}

func (s *Service) CreateShop(shop *models.Shop) (int, error) {
	shop.Timestamp = time.Now().String()
	return s.repo.CreateShop(shop)
}

func (s *Service) GetShops() (models.Store[*models.Shop], error) {
	return s.repo.GetShops()
}

func (s *Service) GetShop(id int) *models.Shop {
	return s.repo.GetShop(id)
}

func (s *Service) UpdateShop(shop *models.Shop) error {
	return s.repo.UpdateShop(shop)
}

func (s *Service) DeleteShop(id int) error {
	return s.repo.DeleteShop(id)
}

func (s *Service) CreateAgent(a *models.Agent) (int, error) {
	a.Timestamp = time.Now().String()
	return s.repo.CreateAgent(a)
}

func (s *Service) GetAgents() (models.Store[*models.Agent], error) {
	return s.repo.GetAgents()
}

func (s *Service) GetAgent(id int) *models.Agent {
	return s.repo.GetAgent(id)
}

func (s *Service) UpdateAgent(a *models.Agent) error {
	return s.repo.UpdateAgent(a)
}

func (s *Service) DeleteAgent(id int) error {
	return s.repo.DeleteAgent(id)
}

func (s *Service) GetDocNames() (models.Store[*models.DocName], error) {
	return s.repo.GetDocNames()
}

func (s *Service) GetDocName(id int) *models.DocName {
	return s.repo.GetDocName(id)
}

func (s *Service) CreateDocument(d *models.Document) (int, error) {
	d.Timestamp = time.Now().String()
	d.GetData = s.getData
	id, err := s.repo.CreateDocument(d)
	if err != nil {
		return 0, err
	}

	if d.Operations.Len() == 0 {
		return id, nil
	}

	d.Operations.ForEach(func(o *models.Operation) error {
		o.DocID = id
		return nil
	})

	err = s.repo.CreateOperations(d.Operations)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetDocuments() (models.Store[*models.Document], error) {
	docs, err := s.repo.GetDocuments()
	if err != nil {
		return nil, err
	}

	err = docs.ForEach(func(doc *models.Document) error {
		doc.GetData = s.getData
		operations, err := s.GetDocOperations(doc.ID)
		if err != nil {
			return err
		}
		operations.ForEach(func(o *models.Operation) error {
			o.GetData = s.getData
			return nil
		})

		doc.Operations = operations
		return nil
	})

	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (s *Service) getData(req models.DataRequest) models.DataResponse {
	var resp models.DataResponse
	if req.DocID != 0 {
		resp.Document = s.GetDocument(req.DocID)
	}
	if req.DocNameID != 0 {
		resp.DocName = s.GetDocName(req.DocNameID)
	}
	if req.AuthorID != 0 {
		resp.Author = s.GetAuthor(req.AuthorID)
	}
	if req.ShopID != 0 {
		resp.Shop = s.GetShop(req.ShopID)
	}
	if req.AgentID != 0 {
		resp.Agent = s.GetAgent(req.AgentID)
	}
	if req.ProdID != 0 {
		resp.Product = s.GetProduct(req.ProdID)
	}
	return resp
}

func (s *Service) GetDocument(id int) *models.Document {
	return s.repo.GetDocument(id)
}

func (s *Service) UpdateDocument(d *models.Document) error {
	err := s.repo.UpdateDocument(d)
	if err != nil {
		return err
	}

	if d.Operations.Len() == 0 {
		return nil
	}
	err = s.repo.UpdateOperations(d.ID, d.Operations)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteDocument(id int) error {
	return s.repo.DeleteDocument(id)
}

func (s *Service) GetDocOperations(docID int) (models.Store[*models.Operation], error) {
	return s.repo.GetOperations(docID)
}

func (s *Service) CreateOperation(o *models.Operation) {
	o.GetData = s.getData
}

func (s *Service) GetAgentsDocuments() (models.Store[*models.Document], error) {
	docs, err := s.GetDocuments()
	if err != nil {
		return nil, err
	}

	agentDocs := models.NewStore[*models.Document]()
	docs.ForEach(func(d *models.Document) error {
		if d.Type == Reception {
			agentDocs.Put(d.Index(), d)
		}
		return nil
	})

	return agentDocs, nil
}

func (s *Service) FilterAgentDocuments(id int) (models.Store[*models.Document], error) {
	docs, err := s.GetDocuments()
	if err != nil {
		return nil, err
	}

	filteredDocs := models.NewStore[*models.Document]()
	docs.ForEach(func(d *models.Document) error {
		if d.AgentID == id {
			filteredDocs.Put(d.Index(), d)
		}
		return nil
	})

	return filteredDocs, nil
}

func (s *Service) GetAquisitionDocuments() (models.Store[*models.Document], error) {
	docs, err := s.GetDocuments()
	if err != nil {
		return nil, err
	}

	filteredDocs := models.NewStore[*models.Document]()
	docs.ForEach(func(d *models.Document) error {
		if d.Type == Aquisiton {
			filteredDocs.Put(d.Index(), d)
		}
		return nil
	})

	return filteredDocs, nil
}

func (s *Service) FilterShopsDocuments(id int) (models.Store[*models.Document], error) {
	docs, err := s.GetDocuments()
	if err != nil {
		return nil, err
	}

	filteredDocs := models.NewStore[*models.Document]()
	docs.ForEach(func(d *models.Document) error {
		if d.ShopID == id {
			filteredDocs.Put(d.Index(), d)
		}
		return nil
	})

	return filteredDocs, nil
}

func (s *Service) GetOperations() (models.Store[*models.Operation], error) {
	docs, err := s.GetDocuments()
	if err != nil {
		return nil, err
	}

	operations := models.NewStore[*models.Operation]()
	docs.ForEach(func(d *models.Document) error {
		if d.Type == Aquisiton {
			d.Operations.ForEach(func(o *models.Operation) error {
				operations.Add(o)
				return nil
			})
		}
		return nil
	})

	return operations, nil
}

func (s *Service) FilterOperations(id int, operations models.Store[*models.Operation]) models.Store[*models.Operation] {
	filteredOperations := models.NewStore[*models.Operation]()
	operations.ForEach(func(o *models.Operation) error {
		if o.ProdID == id {
			filteredOperations.Put(o.Index(), o)
		}
		return nil
	})

	return filteredOperations
}

func (s *Service) Income(year string) (models.Store[*models.Row], []float64, error) {
	documents, err := s.GetDocuments()
	if err != nil {
		return nil, nil, err
	}

	months := make(map[int][]float64)
	err = documents.ForEach(func(d *models.Document) error {
		if d.Type != Reception {
			return nil
		}

		y, err := getYear(d.Date)
		if err != nil {
			return err
		}

		if y != year && year != "" {
			return nil
		}

		month, err := getMonth(d.Date)
		if err != nil {
			return err
		}

		if _, ok := months[d.AuthorID]; !ok {
			months[d.AuthorID] = make([]float64, 12)
		}
		months[d.AuthorID][month] += d.Sum

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	store := models.NewStore[*models.Row]()
	for k, values := range months {
		row := []string{
			s.GetAuthor(k).Name,
		}
		for _, v := range values {
			row = append(row, fmt.Sprintf("%v", v))
		}
		data := &models.Row{
			Idx: k,
			Row: row,
		}
		store.Add(data)
	}

	total := make([]float64, 12)
	for _, v := range months {
		for i, n := range v {
			total[i] += n
		}
	}

	return store, total, nil
}

func getMonth(date string) (int, error) {
	splitted := strings.Split(date, "/")
	if len(splitted) < 2 {
		return 0, errors.New("wrong date")
	}
	return strconv.Atoi(splitted[1])
}

func getYear(date string) (string, error) {
	splitted := strings.Split(date, "/")
	if len(splitted) < 2 {
		return "", errors.New("wrong date")
	}
	return splitted[0], nil
}

func (s *Service) GetOperationsReport(prodID int, startDate, endDate string) (models.Store[*models.Row], []float64, error) {
	documents, err := s.GetDocuments()
	if err != nil {
		return nil, nil, err
	}

	total := make([]float64, 3)
	store := models.NewStore[*models.Row]()
	documents.ForEach(func(d *models.Document) error {
		if d.Type == Aquisiton {
			d.Operations.ForEach(func(o *models.Operation) error {
				if operationCorresponds(o, d.Date, prodID, startDate, endDate) {
					sum := float64(o.Number) * o.Price
					row := []string{
						d.Date,
						s.GetProduct(o.ProdID).Name,
						strconv.Itoa(o.Number),
						fmt.Sprintf("%v", o.Price),
						fmt.Sprintf("%v", sum),
					}
					total[0] += float64(o.Number)
					total[1] += o.Price
					total[2] += sum
					store.Add(&models.Row{
						Row: row,
					})
				}
				return nil
			})

		}
		return nil
	})

	return store, total, nil
}

func operationCorresponds(o *models.Operation, date string, prodID int, startDate, endDate string) bool {
	if prodID == 0 && startDate == "" && endDate == "" {
		return true
	}
	if o.ProdID != prodID {
		return false
	}
	if endDate < date || date < startDate {
		return false
	}
	return true
}

func FindInStore[T models.DataObject[T]](text string, store models.Store[T]) int {
	var i int
	var j int
	store.ForEach(func(t T) error {
		if titleCorresponds(t.Title(), text) {
			i = j
			return errors.New("end")
		}
		j++
		return nil
	})
	return i
}

func titleCorresponds(title, text string) bool {
	if strings.HasPrefix(title, text) {
		return true
	}

	textSplitted := strings.Split(text, " ")
	titleSplitted := strings.Split(title, " ")
	if len(textSplitted) < 2 || len(titleSplitted) < 2 {
		return false
	}

	b := strings.HasPrefix(titleSplitted[1], textSplitted[1])
	if textSplitted[0] == "" {
		return b
	}

	return strings.HasPrefix(titleSplitted[0], textSplitted[0]) && b
}

func (s *Service) GetOverallReport() (float64, float64, error){
	documents, err := s.GetDocuments()
	if err != nil{
		return 0, 0, err
	}

	var received float64
	var expended float64
	documents.ForEach(func(d *models.Document) error {
		if d.Type == Aquisiton{
			expended += d.Sum
		} else if d.Type == Reception{
			received += d.Sum
		}
		return nil
	})

	return received, expended, nil
}