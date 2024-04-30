package ui

import (
	"fmt"
	"budget/iup"
	"budget/pkg/models"
	"budget/pkg/service"
)

func (ui *UI) createAgentsReportWindow() {
	docs, err := ui.service.GetAgentsDocuments()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	agents, err := ui.service.GetAgents()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/reports/report.led")
	w := iup.NewWindow("reportWindow", iup.NewSize(600, 300))

	table := iup.NewTable[*models.Document]("table", w)
	cmbAgents := iup.NewComboBox[*models.Agent]("cmbSortedElement", w)

	table.FillTable(docs)
	cmbAgents.SetValues(agents)

	cmbAgents.OnChange(func(id int) {
		docs, err := ui.service.FilterAgentDocuments(id)
		if err != nil {
			ui.ErrMessage(err)
			return
		}
		table.Clear()
		table.FillTable(docs)
	})

	w.Show()
}

func (ui *UI) createShopsReportWindow() {
	docs, err := ui.service.GetAquisitionDocuments()
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	shops, err := ui.service.GetShops()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/reports/report.led")
	w := iup.NewWindow("reportWindow", iup.NewSize(600, 300))

	table := iup.NewTable[*models.Document]("table", w)
	cmbShops := iup.NewComboBox[*models.Shop]("cmbSortedElement", w)

	cmbShops.SetValues(shops)
	table.FillTable(docs)

	cmbShops.OnChange(func(id int) {
		docs, err := ui.service.FilterShopsDocuments(id)
		if err != nil {
			ui.ErrMessage(err)
			return
		}
		table.Clear()
		table.FillTable(docs)
	})

	w.Show()
}

func (ui *UI) createProductsReportWindow() {
	operations, err := ui.service.GetOperations()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	products, err := ui.service.GetProducts()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/reports/operation.led")
	w := iup.NewWindow("operationWindow", iup.NewSize(600, 300))

	table := iup.NewTable[*models.Operation]("table", w)
	cmbProducts := iup.NewComboBox[*models.Product]("cmbSortElement", w)

	table.FillTable(operations)
	cmbProducts.SetValues(products)

	cmbProducts.OnChange(func(id int) {
		filteredOperations := ui.service.FilterOperations(cmbProducts.Index(), operations)
		table.FillTable(filteredOperations)
	})

	w.Show()
}

func (ui *UI) createIncomeReportWindow() {
	var year string
	income, total, err := ui.service.Income(year)
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	ui.load("/reports/income.led")
	w := iup.NewWindow("incomeWindow", iup.NewSize(950, 300))

	txtYear := iup.NewText("txtYear", w)
	btnShow := iup.NewButton("btnShow", w)

	table := iup.NewTable[*models.Row]("table", w)
	tableTotal := iup.NewTable[*models.Row]("tableTotal", w)

	table.FillTable(income)
	tableTotal.SetHeader(convertToHeader(total, "Всього:"))

	btnShow.OnClick(func() {
		year = txtYear.GetText()
		income, total, err := ui.service.Income(year)
		if err != nil {
			ui.ErrMessage(err)
			return
		}
		table.FillTable(income)
		tableTotal.SetHeader(convertToHeader(total, "Всього:"))
	})

	w.Show()
}

func convertToHeader(list []float64, items ...string) []string {
	res := make([]string, 0, len(list)+len(items))
	res = append(res, items...)
	for _, el := range list {
		res = append(res, fmt.Sprintf("%v", el))
	}
	return res
}

func (ui *UI) createExpendituresReportWindow() {
	products, err := ui.service.GetProducts()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	var prodID int
	var startDate, endDate string
	operations, total, err := ui.service.GetOperationsReport(prodID, startDate, endDate)
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	ui.load("/reports/expenditures.led")
	w := iup.NewWindow("expendituresWindow", iup.NewSize(500, 300))

	txtSearch := iup.NewText("txtSearch", w)
	cmbProduct := iup.NewComboBox[*models.Product]("cmbProduct", w)
	startDatePicker := iup.NewDatePick("startDate", w)
	endDatePicker := iup.NewDatePick("endDate", w)
	btnShow := iup.NewButton("btnShow", w)
	table := iup.NewTable[*models.Row]("table", w)
	tableTotal := iup.NewTable[*models.Row]("tableTotal", w)

	table.FillTable(operations)
	tableTotal.SetHeader(convertToHeader(total, "Всього:", ""))
	cmbProduct.SetValues(products)

	txtSearch.OnChange(func() {
		idx := service.FindInStore[*models.Product](txtSearch.GetText(), products)
		cmbProduct.SetIndex(idx+1)
	})

	btnShow.OnClick(func() {
		prodID = cmbProduct.Index()
		startDate = startDatePicker.Date()
		endDate = endDatePicker.Date()

		operations, total, err := ui.service.GetOperationsReport(prodID, startDate, endDate)
		if err != nil {
			ui.ErrMessage(err)
			return
		}
		table.FillTable(operations)
		tableTotal.SetHeader(convertToHeader(total, "Всього:", ""))
	})

	w.Show()
}
