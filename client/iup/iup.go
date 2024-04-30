package iup

import (
	"fmt"
	"budget/iup/iupim"
	iup "budget/iup/wrapper"
	"budget/pkg/models"
	"log"
	"slices"
	"strconv"
)

/*
 */
import "C"

func Open() {
	iup.Open()
}

func Loop() {
	iup.MainLoop()
}

func Load(path string) {
	iup.Load(path)
}

func Message(title, msg string) {
	iup.Message(title, msg)
}

type Size struct {
	width  int
	height int
}

func NewSize(width, height int) Size {
	return Size{width, height}
}

type Window struct {
	ptr iup.Ihandle
}

func NewWindow(name string, size Size) *Window {
	w := Window{
		ptr: iup.GetHandle(name),
	}
	width := strconv.Itoa(size.width)
	height := strconv.Itoa(size.height)
	iup.SetAttribute(w.ptr, "SIZE", width+"X"+height)
	return &w
}

func (w *Window) SetIcon(path string) {
	img := iupim.LoadImage(path)
	iup.SetAttributeHandle(w.ptr, "ICON", img)
}

func (w *Window) Show() {
	iup.Show(w.ptr)
}

func (w *Window) Ptr() iup.Ihandle {
	return w.ptr
}

type Button struct {
	Ptr iup.Ihandle
}

func NewButton(name string, w *Window) *Button {
	btn := iup.GetChild(w.ptr, name)
	if iup.GetAttribute(btn, "FLAT") == "YES" {
		path := iup.GetAttribute(btn, "IMG")
		img := iupim.LoadImage("../assets/" + path)
		iup.SetAttributeHandle(btn, "IMAGE", img)
	}
	return &Button{
		Ptr: btn,
	}
}

func (btn *Button) OnClick(fn func()) {
	iup.SetCallback(btn.Ptr, "ACTION", func(h iup.Ihandle) int {
		fn()
		return 0
	})
}

type Text struct {
	ptr iup.Ihandle
}

func NewText(name string, w *Window) *Text {
	return &Text{
		ptr: iup.GetChild(w.ptr, name),
	}
}

func (txt *Text) SetText(text string) {
	iup.SetAttribute(txt.ptr, "VALUE", text)
}

func (txt *Text) GetText() string {
	return iup.GetAttribute(txt.ptr, "VALUE")
}

func (txt *Text) OnChange(fn func()) {
	iup.SetCallback(txt.ptr, "VALUECHANGED_CB", func(h iup.Ihandle) int {
		fn()
		return 0
	})
}

type Label struct {
	ptr iup.Ihandle
}

func NewLabel(name string, w *Window) *Label {
	return &Label{
		ptr: iup.GetChild(w.ptr, name),
	}
}

func (lbl *Label) SetText(text string) {
	iup.SetAttribute(lbl.ptr, "TITLE", text)
}

func (lbl *Label) GetText() string {
	return iup.GetAttribute(lbl.ptr, "TITLE")
}

type ComboValue[T any] interface {
	Index() int
	Title() string
	models.CacheObject[T]
}

type ComboBox[T ComboValue[T]] struct {
	Ptr     iup.Ihandle
	indexes []int
}

func NewComboBox[T ComboValue[T]](name string, w *Window) *ComboBox[T] {
	return &ComboBox[T]{
		Ptr: iup.GetChild(w.ptr, name),
	}
}

func (cmb *ComboBox[T]) OnChange(fn func(int)) {
	iup.SetCallback(cmb.Ptr, "ACTION", func(h iup.Ihandle, _ *C.char, item, state int) int {
		fn(cmb.Index())
		return 0
	})
}

func (cmb *ComboBox[T]) SetValues(values models.Store[T]) {
	indexes := make([]int, 0, values.Len())

	var i int
	values.ForEach(func(t T) error {
		iup.SetAttribute(cmb.Ptr, strconv.Itoa(i+1), t.Title())
		indexes = append(indexes, t.Index())
		i++
		return nil
	})

	cmb.indexes = indexes
}

func (cmb *ComboBox[T]) Index() int {
	value := iup.GetAttribute(cmb.Ptr, "VALUE")
	idx, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return cmb.indexes[idx-1]
}

func (cmb *ComboBox[T]) SetIndex(idx int) {
	iup.SetAttribute(cmb.Ptr, "VALUE", strconv.Itoa(idx))
}

func (cmb *ComboBox[T]) SetValue(value string) {
	iup.SetAttribute(cmb.Ptr, "VALUESTRING", value)
}

type DatePick struct {
	Ptr iup.Ihandle
}

func NewDatePick(name string, w *Window) *DatePick {
	return &DatePick{
		Ptr: iup.GetChild(w.ptr, name),
	}
}

func (d *DatePick) SetDate(day, month, year int) {
	iup.SetAttribute(d.Ptr, "VALUE", fmt.Sprintf("%d/%d/%d", year, day, month))
}

func (d *DatePick) SetValue(v string) {
	iup.SetAttribute(d.Ptr, "VALUE", v)
}

func (d *DatePick) Date() string {
	return iup.GetAttribute(d.Ptr, "VALUE")
}

type TableRow[T any] interface {
	ToRow() []string
	models.Object[T]
}

type Table[T TableRow[T]] struct {
	Ptr           iup.Ihandle
	onClick       func(int)
	onDblClick    func(int)
	onRightClick  func(int)
	onHeaderClick func(int)
	indexes       []int
}

func NewTable[T TableRow[T]](name string, w *Window) *Table[T] {
	tbl := &Table[T]{
		Ptr: iup.GetChild(w.ptr, name),
	}
	tbl.Ptr.SetAttribute("INDEX", "-1")
	tbl.setHandles()

	return tbl
}

func (tbl *Table[T]) setHandles() {
	iup.SetCallback(tbl.Ptr, "CLICK_CB", func(ptr iup.Ihandle, line int, col int, status *C.char) int {
		id := tbl.indexes[line-1]

		if line == 0 {
			if tbl.onHeaderClick != nil {
				tbl.onHeaderClick(id)
			}
			return 0
		}
		tbl.highlighLine(line)

		switch s := C.GoString(status); {
		case s[5] == 'D':
			if tbl.onDblClick != nil {
				tbl.onDblClick(id)
			}
		case s[2] == '1':
			if tbl.onClick != nil {
				tbl.onClick(id)
			}
		case s[4] == '3':
			if tbl.onRightClick != nil {
				tbl.onRightClick(id)
			}
		}

		return 0
	})
}

func (tbl *Table[T]) highlighLine(line int) {
	index := iup.GetAttribute(tbl.Ptr, "INDEX")
	if index != "-1" {
		iup.SetAttribute(tbl.Ptr, "BGCOLOR"+index+":*", "255 255 255")
		iup.SetAttribute(tbl.Ptr, "REDRAW", "L"+index)
	}
	l := strconv.Itoa(line)
	iup.SetAttribute(tbl.Ptr, "INDEX", l)
	iup.SetAttribute(tbl.Ptr, "BGCOLOR"+l+":*", "145 201 247")
	iup.SetAttribute(tbl.Ptr, "REDRAW", "L"+l)
}

func (tbl *Table[T]) OnClick(fn func(int)) {
	tbl.onClick = fn
}

func (tbl *Table[T]) OnDblClick(fn func(int)) {
	tbl.onDblClick = fn
}

func (tbl *Table[T]) OnHeaderClick(fn func(int)) {
	tbl.onHeaderClick = fn
}

func (tbl *Table[T]) OnRightClick(fn func(int)) {
	tbl.onRightClick = fn
}

func (tbl *Table[T]) FillTable(data models.Store[T]) {
	n, _ := strconv.Atoi(iup.GetAttribute(tbl.Ptr, "NUMCOL"))
	iup.SetAttribute(tbl.Ptr, "NUMLIN", strconv.Itoa(data.Len()))
	indexes := make([]int, 0, data.Len())

	var i int
	data.ForEach(func(t T) error {
		row := t.ToRow()
		for j := 0; j < n; j++ {
			iup.SetAttribute(tbl.Ptr, strconv.Itoa(i+1)+":"+strconv.Itoa(j+1), row[j])
		}
		indexes = append(indexes, t.Index())
		i++
		return nil
	})
	tbl.indexes = indexes

	iup.SetAttribute(tbl.Ptr, "REDRAW", "true")

	tbl.highlighLine(1)
	iup.SetAttribute(tbl.Ptr, "INDEX", "1")
}

func (tbl *Table[T]) SetHeader(header []string) {
	for i, v := range header {
		iup.SetAttribute(tbl.Ptr, "0:"+strconv.Itoa(i+1), v)
	}
	iup.SetAttribute(tbl.Ptr, "REDRAW", "true")
}

func (tbl *Table[T]) Clear() {
	len := iup.GetAttribute(tbl.Ptr, "NUMLIN")
	iup.SetAttribute(tbl.Ptr, "DELLIN", "0-"+len)
}

func (tbl *Table[T]) Index() int {
	if tbl.indexes == nil {
		return -1
	}
	idx, err := strconv.Atoi(iup.GetAttribute(tbl.Ptr, "INDEX"))
	if err != nil {
		log.Fatal(err)
	}
	return tbl.indexes[idx-1]
}

func (tbl *Table[T]) Add(row T) {
	tbl.indexes = append(tbl.indexes, row.Index())
	len := iup.GetAttribute(tbl.Ptr, "NUMLIN")
	i, _ := strconv.Atoi(len)
	iup.SetAttribute(tbl.Ptr, "ADDLIN", len)
	for j, value := range row.ToRow() {
		iup.SetAttribute(tbl.Ptr, strconv.Itoa(i+1)+":"+strconv.Itoa(j+1), value)
		iup.SetAttribute(tbl.Ptr, "REDRAW", "L"+strconv.Itoa(i+1))
	}
	tbl.highlighLine(i + 1)
}

func (tbl *Table[T]) Change(row T) {
	index := iup.GetAttribute(tbl.Ptr, "INDEX")
	idx, _ := strconv.Atoi(index)
	for j, value := range row.ToRow() {
		iup.SetAttribute(tbl.Ptr, strconv.Itoa(idx)+":"+strconv.Itoa(j+1), value)
		iup.SetAttribute(tbl.Ptr, "REDRAW", "L"+strconv.Itoa(idx+1))
	}
}

func (tbl *Table[T]) Delete() {
	index := iup.GetAttribute(tbl.Ptr, "INDEX")
	idx, _ := strconv.Atoi(index)
	tbl.indexes = slices.Delete[[]int](tbl.indexes, idx-1, idx)
	iup.SetAttribute(tbl.Ptr, "DELLIN", index)
}

type SubMenu struct {
	Ptr iup.Ihandle
}

func NewSubmenu(name string, w *Window) *SubMenu {
	return &SubMenu{
		Ptr: iup.GetChild(w.Ptr(), name),
	}
}

func (menu *SubMenu) OnClick(fn func()) {
	iup.SetCallback(menu.Ptr, "ACTION", func(h iup.Ihandle) int {
		fn()
		return 0
	})
}

type CheckBox struct {
	Ptr iup.Ihandle
}

func NewCheckBox(name string, w *Window) *CheckBox {
	return &CheckBox{
		Ptr: iup.GetChild(w.Ptr(), name),
	}
}

func (chb *CheckBox) Value() bool {
	value := iup.GetAttribute(chb.Ptr, "VALUE")
	return value == "ON"
}

func (chb *CheckBox) SetValue(value bool) {
	v := "OFF"
	if value {
		v = "ON"
	}
	iup.SetAttribute(chb.Ptr, "VALUE", v)
}
