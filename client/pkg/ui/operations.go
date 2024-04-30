package ui

import (
	"fmt"
	"budget/iup"
	"budget/pkg/models"
	"budget/pkg/service"
	"strconv"
)

func (ui *UI) createLendingWindow(onDocCreated func(d *models.Document)) {
	authors, err := ui.service.GetAuthors()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/operations/lending.led")

	w := iup.NewWindow("lendingWindow", iup.NewSize(200, 130))

	txtSum := iup.NewText("txtSum", w)
	cmbAuthor := iup.NewComboBox[*models.Author]("cmbAuthor", w)
	datePick := iup.NewDatePick("datePick", w)
	btnSubmit := iup.NewButton("btnSubmit", w)

	cmbAuthor.SetValues(authors)

	btnSubmit.OnClick(func() {
		sum, err := strconv.ParseFloat(txtSum.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		doc := &models.Document{
			Type:       service.Lending,
			Sum:        sum,
			Date:       datePick.Date(),
			AuthorID:   cmbAuthor.Index(),
			Operations: models.NewStore[*models.Operation](),
		}

		ui.service.CreateDocument(doc)
		onDocCreated(doc)
	})

	w.Show()
}

func (ui *UI) createReceptionWindow(onDocCreated func(d *models.Document)) {
	authors, err := ui.service.GetAuthors()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	agents, err := ui.service.GetAgents()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/operations/reception.led")

	w := iup.NewWindow("receptionWindow", iup.NewSize(200, 150))

	txtSum := iup.NewText("txtSum", w)
	cmbAuthor := iup.NewComboBox[*models.Author]("cmbAuthor", w)
	datePick := iup.NewDatePick("datePick", w)
	btnSubmit := iup.NewButton("btnSubmit", w)
	txtAgent := iup.NewText("txtAgent", w)
	cmbAgent := iup.NewComboBox[*models.Agent]("cmbAgent", w)

	cmbAuthor.SetValues(authors)
	cmbAgent.SetValues(agents)

	txtAgent.OnChange(func() {
		idx := service.FindInStore[*models.Agent](txtAgent.GetText(), agents)
		cmbAgent.SetIndex(idx + 1)
	})

	btnSubmit.OnClick(func() {
		sum, err := strconv.ParseFloat(txtSum.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		doc := &models.Document{
			Type:       service.Reception,
			Sum:        sum,
			AuthorID:   cmbAuthor.Index(),
			Date:       datePick.Date(),
			AgentID:    cmbAgent.Index(),
			Operations: models.NewStore[*models.Operation](),
		}

		ui.service.CreateDocument(doc)
		onDocCreated(doc)
	})

	w.Show()
}

func (ui *UI) createAcquisitionWindow(onDocCreated func(d *models.Document)) {
	d := &models.Document{
		Operations: models.NewStore[*models.Operation](),
	}

	authors, err := ui.service.GetAuthors()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	shops, err := ui.service.GetShops()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	products, err := ui.service.GetProducts()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/operations/acquisition.led")

	w := iup.NewWindow("acquisitionWindow", iup.NewSize(650, 300))

	txtSum := iup.NewText("txtSum", w)
	cmbAuthor := iup.NewComboBox[*models.Author]("cmbAuthor", w)
	datePick := iup.NewDatePick("datePick", w)
	btnSubmit := iup.NewButton("btnSubmit", w)
	txtShop := iup.NewText("txtShop", w)
	cmbShop := iup.NewComboBox[*models.Shop]("cmbShop", w)

	txtProduct := iup.NewText("txtProduct", w)
	cmbProduct := iup.NewComboBox[*models.Product]("cmbProduct", w)
	txtPrice := iup.NewText("txtPrice", w)
	txtCount := iup.NewText("txtCount", w)
	// txtOperationSum := iup.NewText("txtOperationSum", w)

	btnCreate := iup.NewButton("btnCreate", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)

	table := iup.NewTable[*models.Operation]("table", w)

	cmbProduct.SetValues(products)
	cmbAuthor.SetValues(authors)
	cmbShop.SetValues(shops)

	table.OnClick(func(id int) {
		o := d.Operations.Get(id)
		cmbProduct.SetValue(products.Get(o.ProdID).Name)
		txtPrice.SetText(fmt.Sprintf("%v", o.Price))
		txtCount.SetText(fmt.Sprintf("%v", o.Number))
	})

	txtShop.OnChange(func() {
		idx := service.FindInStore[*models.Shop](txtShop.GetText(), shops)
		cmbShop.SetIndex(idx + 1)
	})

	txtProduct.OnChange(func() {
		idx := service.FindInStore[*models.Product](txtProduct.GetText(), products)
		cmbProduct.SetIndex(idx + 1)
	})

	btnCreate.OnClick(func() {
		count, err := strconv.Atoi(txtCount.GetText())
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		price, err := strconv.ParseFloat(txtPrice.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		operation := &models.Operation{
			ProdID: cmbProduct.Index(),
			Number: count,
			Price:  price,
		}

		ui.service.CreateOperation(operation)
		d.Operations.Add(operation)
		table.Add(operation)
	})

	btnChange.OnClick(func() {
		id := table.Index()
		count, err := strconv.Atoi(txtCount.GetText())
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		price, err := strconv.ParseFloat(txtPrice.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		operation := &models.Operation{
			ID:     id,
			ProdID: cmbProduct.Index(),
			Number: count,
			Price:  price,
		}

		ui.service.CreateOperation(operation)
		d.Operations.Put(operation.ID, operation)
		table.Change(operation)
	})

	btnDelete.OnClick(func() {
		d.Operations.Delete(table.Index())
		table.Delete()
	})

	btnSubmit.OnClick(func() {
		sum, err := strconv.ParseFloat(txtSum.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		d.Type = service.Aquisiton
		d.Sum = sum
		d.AuthorID = cmbAuthor.Index()
		d.ShopID = cmbShop.Index()
		d.Date = datePick.Date()

		ui.service.CreateDocument(d)
		onDocCreated(d)
	})

	w.Show()
}

func (ui *UI) createReceptionChangeWindow(doc *models.Document, onDocChanged func(*models.Document), onDocDeleted func()) {
	authors, err := ui.service.GetAuthors()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	agents, err := ui.service.GetAgents()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/operations/receptionChange.led")

	w := iup.NewWindow("receptionChangeWindow", iup.NewSize(200, 150))

	txtSum := iup.NewText("txtSum", w)
	cmbAuthor := iup.NewComboBox[*models.Author]("cmbAuthor", w)
	datePick := iup.NewDatePick("datePick", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)
	txtAgent := iup.NewText("txtAgent", w)
	cmbAgent := iup.NewComboBox[*models.Agent]("cmbAgent", w)

	cmbAuthor.SetValues(authors)
	cmbAuthor.SetValue(ui.service.GetAuthor(doc.AuthorID).Name)
	cmbAgent.SetValues(agents)
	cmbAgent.SetValue(ui.service.GetAgent(doc.AgentID).Name)

	txtSum.SetText(fmt.Sprintf("%v", doc.Sum))
	datePick.SetValue(doc.Date)

	txtAgent.OnChange(func() {
		idx := service.FindInStore[*models.Agent](txtAgent.GetText(), agents)
		cmbAgent.SetIndex(idx + 1)
	})

	btnChange.OnClick(func() {
		sum, err := strconv.ParseFloat(txtSum.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		doc := &models.Document{
			ID:         doc.ID,
			Type:       service.Reception,
			Sum:        sum,
			AuthorID:   cmbAuthor.Index(),
			Date:       datePick.Date(),
			AgentID:    cmbAgent.Index(),
			Operations: doc.Operations,
			GetData:    doc.GetData,
		}

		err = ui.service.UpdateDocument(doc)
		if err != nil {
			ui.ErrMessage(err)
		}
		onDocChanged(doc)
	})

	btnDelete.OnClick(func() {
		err := ui.service.DeleteDocument(doc.ID)
		if err != nil {
			ui.ErrMessage(err)
		}
		onDocDeleted()
	})

	w.Show()
}

func (ui *UI) createLendingChangeWindow(doc *models.Document, onDocChanged func(*models.Document), onDocDeleted func()) {
	authors, err := ui.service.GetAuthors()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/operations/lendingChange.led")

	w := iup.NewWindow("lendingChangeWindow", iup.NewSize(200, 150))

	txtSum := iup.NewText("txtSum", w)
	cmbAuthor := iup.NewComboBox[*models.Author]("cmbAuthor", w)
	datePick := iup.NewDatePick("datePick", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)

	cmbAuthor.SetValues(authors)
	txtSum.SetText(fmt.Sprintf("%v", doc.Sum))
	datePick.SetValue(doc.Date)

	btnChange.OnClick(func() {
		sum, err := strconv.ParseFloat(txtSum.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		doc := &models.Document{
			ID:         doc.ID,
			Type:       service.Lending,
			Sum:        sum,
			AuthorID:   cmbAuthor.Index(),
			Date:       datePick.Date(),
			Operations: doc.Operations,
			GetData:    doc.GetData,
		}

		err = ui.service.UpdateDocument(doc)
		if err != nil {
			ui.ErrMessage(err)
		}
		onDocChanged(doc)
	})

	btnDelete.OnClick(func() {
		err := ui.service.DeleteDocument(doc.ID)
		if err != nil {
			ui.ErrMessage(err)
		}
		onDocDeleted()
	})

	w.Show()
}

func (ui *UI) createAquisitionChangeWindow(d *models.Document, onDocChanged func(*models.Document), onDocDeleted func()) {
	authors, err := ui.service.GetAuthors()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	shops, err := ui.service.GetShops()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	products, err := ui.service.GetProducts()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/operations/aquisitionChange.led")

	w := iup.NewWindow("aquisitionChangeWindow", iup.NewSize(650, 300))

	txtSum := iup.NewText("txtSum", w)
	cmbAuthor := iup.NewComboBox[*models.Author]("cmbAuthor", w)
	datePick := iup.NewDatePick("datePick", w)
	btnChangeDoc := iup.NewButton("btnChangeDoc", w)
	btnDeleteDoc := iup.NewButton("btnDeleteDoc", w)
	txtShop := iup.NewText("txtShop", w)
	cmbShop := iup.NewComboBox[*models.Shop]("cmbShop", w)

	txtProduct := iup.NewText("txtProduct", w)
	cmbProduct := iup.NewComboBox[*models.Product]("cmbProduct", w)
	txtPrice := iup.NewText("txtPrice", w)
	txtCount := iup.NewText("txtCount", w)
	// txtOperationSum := iup.NewText("txtOperationSum", w)

	btnCreate := iup.NewButton("btnCreate", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)

	table := iup.NewTable[*models.Operation]("table", w)

	txtSum.SetText(fmt.Sprintf("%v", d.Sum))
	datePick.SetValue(d.Date)
	cmbProduct.SetValues(products)
	cmbAuthor.SetValues(authors)
	cmbShop.SetValues(shops)

	cmbAuthor.SetValue(ui.service.GetAuthor(d.AuthorID).Name)
	cmbShop.SetValue(ui.service.GetShop(d.ShopID).Name)

	table.OnClick(func(id int) {
		o := d.Operations.Get(id)
		cmbProduct.SetValue(products.Get(o.ProdID).Name)
		txtPrice.SetText(fmt.Sprintf("%v", o.Price))
		txtCount.SetText(fmt.Sprintf("%v", o.Number))
	})

	table.FillTable(d.Operations)

	txtShop.OnChange(func() {
		idx := service.FindInStore[*models.Shop](txtShop.GetText(), shops)
		cmbShop.SetIndex(idx + 1)
	})

	txtProduct.OnChange(func() {
		idx := service.FindInStore[*models.Product](txtProduct.GetText(), products)
		cmbProduct.SetIndex(idx + 1)
	})

	btnCreate.OnClick(func() {
		count, err := strconv.Atoi(txtCount.GetText())
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		price, err := strconv.ParseFloat(txtPrice.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		operation := &models.Operation{
			ProdID: cmbProduct.Index(),
			Number: count,
			Price:  price,
			DocID:  d.ID,
		}

		ui.service.CreateOperation(operation)
		d.Operations.Add(operation)
		table.Add(operation)
	})

	btnChange.OnClick(func() {
		id := table.Index()
		count, err := strconv.Atoi(txtCount.GetText())
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		price, err := strconv.ParseFloat(txtPrice.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		operation := &models.Operation{
			ID:     id,
			ProdID: cmbProduct.Index(),
			Number: count,
			Price:  price,
			DocID:  d.ID,
		}

		ui.service.CreateOperation(operation)
		d.Operations.Put(operation.ID, operation)
		table.Change(operation)
	})

	btnDelete.OnClick(func() {
		d.Operations.Delete(table.Index())
		table.Delete()
	})

	btnChangeDoc.OnClick(func() {
		sum, err := strconv.ParseFloat(txtSum.GetText(), 64)
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		d.Type = service.Aquisiton
		d.Sum = sum
		d.AuthorID = cmbAuthor.Index()
		d.ShopID = cmbShop.Index()
		d.Date = datePick.Date()

		err = ui.service.UpdateDocument(d)
		if err != nil {
			ui.ErrMessage(err)
			return
		}
		onDocChanged(d)
	})

	btnDeleteDoc.OnClick(func() {
		onDocDeleted()
	})

	w.Show()
}
