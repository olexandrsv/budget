package ui

import (
	"budget/iup"
	"budget/pkg/models"
	"budget/pkg/service"
	"fmt"
	"strconv"

	"budget/iup/wrapper"
)

type UI struct {
	ledsPath string
	service  *service.Service
}

func New(path string, s *service.Service) *UI {
	return &UI{
		ledsPath: path,
		service:  s,
	}
}

func (ui *UI) Run() {
	ui.service.GetDocNames()
	ui.service.GetAgents()
	ui.service.GetAuthors()
	ui.service.GetShops()
	ui.service.GetProducts()

	docs, err := ui.service.GetDocuments()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	iup.Open()
	ui.load("/main.led")
	w := iup.NewWindow("mainWindow", iup.NewSize(650, 300))

	itemProducts := iup.NewSubmenu("itemProducts", w)
	itemAuthors := iup.NewSubmenu("itemAuthors", w)
	itemShops := iup.NewSubmenu("itemShops", w)
	itemAgents := iup.NewSubmenu("itemAgents", w)
	itemDocNames := iup.NewSubmenu("itemDocNames", w)

	itemLending := iup.NewSubmenu("itemLending", w)
	itemReception := iup.NewSubmenu("itemReception", w)
	itemAcquisition := iup.NewSubmenu("itemAcquisition", w)

	itemAgentReport := iup.NewSubmenu("itemAgentReport", w)
	itemShopReport := iup.NewSubmenu("itemShopReport", w)
	itemProductReport := iup.NewSubmenu("itemProductReport", w)
	itemIncomeReport := iup.NewSubmenu("itemIncomeReport", w)
	itemExpenditures := iup.NewSubmenu("itemExpendituresReport", w)

	btnAdd := iup.NewButton("btnAdd", w)
	btnDelete := iup.NewButton("btnDelete", w)
	btnLend := iup.NewButton("btnLend", w)
	btnRevenue := iup.NewButton("btnRevenue", w)
	btnExpenditures := iup.NewButton("btnExpenditures", w)
	_ = btnAdd
	_ = btnDelete
	_ = btnLend
	_ = btnRevenue
	_ = btnExpenditures

	table := iup.NewTable[*models.Document]("table", w)
	table.FillTable(docs)

	onDocCreated := func(d *models.Document) {
		table.Add(d)
		ui.changeStatisticForm(w, d)
	}
	onDocChanged := func(d *models.Document) {
		table.Change(d)
		ui.fillStatisticForm(w)
	}
	onDocDeleted := func() {
		table.Delete()
		ui.fillStatisticForm(w)
	}

	table.OnDblClick(func(id int) {
		doc := ui.service.GetDocument(id)
		switch doc.Type {
		case service.Reception:
			ui.createReceptionChangeWindow(doc, onDocChanged, onDocDeleted)
		case service.Lending:
			ui.createLendingChangeWindow(doc, onDocChanged, onDocDeleted)
		case service.Aquisiton:
			ui.createAquisitionChangeWindow(doc, onDocChanged, onDocDeleted)
		}
		ui.fillStatisticForm(w)
	})

	btnAdd.OnClick(func() {
		ui.createReceptionWindow(onDocCreated)
	})

	btnDelete.OnClick(func() {
		ui.createAcquisitionWindow(onDocCreated)
	})

	btnLend.OnClick(func() {
		ui.createAgentsReportWindow()
	})

	btnRevenue.OnClick(func() {
		ui.createIncomeReportWindow()
	})

	btnExpenditures.OnClick(func() {
		ui.createExpendituresReportWindow()
	})

	table.OnRightClick(func(i int) {
		ptr := wrapper.GetHandle("popupMenu")
		wrapper.Popup(ptr, 65532, 65532)
	})

	itemProducts.OnClick(func() {
		ui.createProductsWindow()
	})

	itemAuthors.OnClick(func() {
		ui.createAuthorsWindow()
	})

	itemShops.OnClick(func() {
		ui.createShopsWindow()
	})

	itemAgents.OnClick(func() {
		ui.createAgentsWindow()
	})

	itemDocNames.OnClick(func() {
		ui.createDocNamesWindow()
	})

	itemLending.OnClick(func() {
		ui.createLendingWindow(onDocCreated)
	})

	itemReception.OnClick(func() {
		ui.createReceptionWindow(onDocCreated)
	})

	itemAcquisition.OnClick(func() {
		ui.createAcquisitionWindow(onDocCreated)
	})

	itemAgentReport.OnClick(func() {
		ui.createAgentsReportWindow()
	})

	itemShopReport.OnClick(func() {
		ui.createShopsReportWindow()
	})

	itemProductReport.OnClick(func() {
		ui.createProductsReportWindow()
	})

	itemIncomeReport.OnClick(func() {
		ui.createIncomeReportWindow()
	})

	itemExpenditures.OnClick(func() {
		ui.createExpendituresReportWindow()
	})

	ui.fillStatisticForm(w)

	w.Show()
	iup.Loop()
}

func (ui *UI) fillStatisticForm(w *iup.Window) {
	received, expended, err := ui.service.GetOverallReport()
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	txtReceived := iup.NewText("txtReceived", w)
	txtExpended := iup.NewText("txtExpended", w)
	txtBalance := iup.NewText("txtBalance", w)

	txtReceived.SetText(fmt.Sprintf("%v", received))
	txtExpended.SetText(fmt.Sprintf("%v", expended))
	txtBalance.SetText(fmt.Sprintf("%v", received-expended))
}

func (ui *UI) changeStatisticForm(w *iup.Window, d *models.Document) {
	txtReceived := iup.NewText("txtReceived", w)
	txtExpended := iup.NewText("txtExpended", w)
	txtBalance := iup.NewText("txtBalance", w)

	received, err := strconv.ParseFloat(txtReceived.GetText(), 64)
	if err != nil {
		ui.ErrMessage(err)
		return
	}
	expended, err := strconv.ParseFloat(txtExpended.GetText(), 64)
	if err != nil {
		ui.ErrMessage(err)
		return
	}

	if d.Type == service.Aquisiton {
		expended += d.Sum
	} else if d.Type == service.Reception {
		received += d.Sum
	}

	txtReceived.SetText(fmt.Sprintf("%v", received))
	txtExpended.SetText(fmt.Sprintf("%v", expended))
	txtBalance.SetText(fmt.Sprintf("%v", received-expended))
}

func (ui *UI) load(path string) {
	iup.Load(ui.ledsPath + path)
}

func (ui *UI) ErrMessage(err error) {
	fmt.Println(err)
	iup.Message("Error", err.Error())
}
