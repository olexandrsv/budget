package ui

import (
	"fmt"
	"budget/iup"
	"budget/pkg/models"
	"strconv"
)

func (ui *UI) createProductsWindow() {
	products, err := ui.service.GetProducts()
	if err != nil {
		fmt.Println(err)
	}

	ui.load("/digest/products.led")
	w := iup.NewWindow("productsWindow", iup.NewSize(310, 250))

	btnCreate := iup.NewButton("btnCreateProduct", w)
	btnChange := iup.NewButton("btnChangeProduct", w)
	btnDelete := iup.NewButton("btnDeleteProduct", w)

	lblID := iup.NewLabel("lblID", w)

	txtName := iup.NewText("txtName", w)
	txtPrice := iup.NewText("txtPrice", w)

	chbCountable := iup.NewCheckBox("chbCountable", w)
	tblProducts := iup.NewTable[*models.Product]("tblProducts", w)

	tblProducts.FillTable(products)

	tblProducts.OnClick(func(id int) {
		product := ui.service.GetProduct(id)

		lblID.SetText(strconv.Itoa(product.ID))
		txtName.SetText(product.Name)
		txtPrice.SetText(fmt.Sprintf("%v", product.Price))
		chbCountable.SetValue(product.Countable)
	})

	btnCreate.OnClick(func() {
		price, err := strconv.ParseFloat(txtPrice.GetText(), 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		product := &models.Product{
			Name:      txtName.GetText(),
			Price:     price,
			Countable: chbCountable.Value(),
		}

		id, _ := ui.service.CreateProduct(product)
		lblID.SetText(strconv.Itoa(id))
		tblProducts.Add(product)
	})

	btnChange.OnClick(func() {
		price, err := strconv.ParseFloat(txtPrice.GetText(), 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			fmt.Println(err)
			return
		}

		product := &models.Product{
			ID:        id,
			Name:      txtName.GetText(),
			Price:     price,
			Countable: chbCountable.Value(),
		}

		ui.service.UpdateProduct(product)
		tblProducts.Change(product)
	})

	btnDelete.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			fmt.Println(err)
			return
		}

		ui.service.DeleteProduct(id)
		tblProducts.Delete()
	})

	w.Show()
}

func (ui *UI) createAuthorsWindow() {
	authors, err := ui.service.GetAuthors()
	if err != nil {
		fmt.Println(err)
		return
	}

	ui.load("/digest/authors.led")
	w := iup.NewWindow("authorsWindow", iup.NewSize(280, 250))

	btnCreate := iup.NewButton("btnCreate", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)

	txtName := iup.NewText("txtName", w)
	txtCard := iup.NewText("txtCard", w)

	lblID := iup.NewLabel("lblID", w)

	tblAuthors := iup.NewTable[*models.Author]("tblAuthors", w)
	tblAuthors.FillTable(authors)

	tblAuthors.OnClick(func(id int) {
		author := ui.service.GetAuthor(id)

		lblID.SetText(strconv.Itoa(author.ID))
		txtName.SetText(author.Name)
		txtCard.SetText(author.Card)
	})

	btnCreate.OnClick(func() {
		author := &models.Author{
			Name: txtName.GetText(),
			Card: txtCard.GetText(),
		}

		id, _ := ui.service.CreateAuthor(author)
		lblID.SetText(strconv.Itoa(id))
		tblAuthors.Add(author)
	})

	btnChange.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			fmt.Println(err)
			return
		}

		author := &models.Author{
			ID:   id,
			Name: txtName.GetText(),
			Card: txtCard.GetText(),
		}

		ui.service.UpdateAuthor(author)
		tblAuthors.Change(author)
	})

	btnDelete.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			fmt.Println(err)
			return
		}

		ui.service.DeleteAuthor(id)
		tblAuthors.Delete()
	})

	w.Show()
}

func (ui *UI) createShopsWindow() {
	shops, err := ui.service.GetShops()
	if err != nil {
		fmt.Println(err)
		return
	}

	ui.load("/digest/shops.led")
	w := iup.NewWindow("shopsWindow", iup.NewSize(280, 250))

	btnCreate := iup.NewButton("btnCreate", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)

	txtName := iup.NewText("txtName", w)

	lblID := iup.NewLabel("lblID", w)

	table := iup.NewTable[*models.Shop]("table", w)
	table.FillTable(shops)

	table.OnClick(func(id int) {
		shop := ui.service.GetShop(id)

		lblID.SetText(strconv.Itoa(shop.ID))
		txtName.SetText(shop.Name)
	})

	btnCreate.OnClick(func() {
		shop := &models.Shop{
			Name: txtName.GetText(),
		}

		id, _ := ui.service.CreateShop(shop)
		lblID.SetText(strconv.Itoa(id))
		table.Add(shop)
	})

	btnChange.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			fmt.Println(err)
			return
		}

		shop := &models.Shop{
			ID:   id,
			Name: txtName.GetText(),
		}

		ui.service.UpdateShop(shop)
		table.Change(shop)
	})

	btnDelete.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			fmt.Println(err)
			return
		}

		ui.service.DeleteShop(id)
		table.Delete()
	})

	w.Show()
}

func (ui *UI) createAgentsWindow() {
	agents, err := ui.service.GetAgents()
	if err != nil {
		fmt.Println(err)
		return
	}

	ui.load("/digest/agents.led")
	w := iup.NewWindow("agentsWindow", iup.NewSize(310, 250))

	btnCreate := iup.NewButton("btnCreate", w)
	btnChange := iup.NewButton("btnChange", w)
	btnDelete := iup.NewButton("btnDelete", w)

	txtName := iup.NewText("txtName", w)
	txtType := iup.NewText("txtType", w)
	txtTel := iup.NewText("txtTel", w)

	lblID := iup.NewLabel("lblID", w)

	table := iup.NewTable[*models.Agent]("table", w)
	table.FillTable(agents)

	table.OnClick(func(id int) {
		agent := ui.service.GetAgent(id)

		lblID.SetText(strconv.Itoa(agent.ID))
		txtName.SetText(agent.Name)
		txtType.SetText(agent.Type)
		txtTel.SetText(agent.Tel)
	})

	btnCreate.OnClick(func() {
		agent := &models.Agent{
			Name: txtName.GetText(),
			Type: txtType.GetText(),
			Tel:  txtTel.GetText(),
		}

		id, _ := ui.service.CreateAgent(agent)
		lblID.SetText(strconv.Itoa(id))
		table.Add(agent)
	})

	btnChange.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		agent := &models.Agent{
			ID:   id,
			Name: txtName.GetText(),
			Type: txtType.GetText(),
			Tel:  txtTel.GetText(),
		}

		ui.service.UpdateAgent(agent)
		table.Change(agent)
	})

	btnDelete.OnClick(func() {
		id, err := strconv.Atoi(lblID.GetText())
		if err != nil {
			ui.ErrMessage(err)
			return
		}

		ui.service.DeleteAgent(id)
		table.Delete()
	})

	w.Show()
}

func (ui *UI) createDocNamesWindow() {
	docNames, err := ui.service.GetDocNames()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/digest/docNames.led")
	w := iup.NewWindow("docNamesWindow", iup.NewSize(180, 180))

	txtName := iup.NewText("txtName", w)

	lblID := iup.NewLabel("lblID", w)

	table := iup.NewTable[*models.DocName]("table", w)
	table.FillTable(docNames)

	table.OnClick(func(id int) {
		docName := ui.service.GetDocName(id)

		lblID.SetText(strconv.Itoa(docName.ID))
		txtName.SetText(docName.Name)
	})

	w.Show()
}
