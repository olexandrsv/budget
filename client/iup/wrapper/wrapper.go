package wrapper

import (
	"budget/iup/iupcontrols"
	"syscall"
	"unsafe"

	"github.com/ying32/dylib"
)

/*
 */
import "C"

type Ihandle uintptr

const (
	CENTER       = 0xFFFF // iup.Popup and iup.ShowXY parameter value
	LEFT         = 0xFFFE // iup.Popup and iup.ShowXY parameter value
	RIGHT        = 0xFFFD // iup.Popup and iup.ShowXY parameter value
	MOUSEPOS     = 0xFFFC // iup.Popup and iup.ShowXY parameter value
	CURRENT      = 0xFFFB // iup.Popup and iup.ShowXY parameter value
	CENTERPARENT = 0xFFFA // iup.Popup and iup.ShowXY parameter value
	TOP          = LEFT   // iup.Popup and iup.ShowXY parameter value
	BOTTOM       = RIGHT  // iup.Popup and iup.ShowXY parameter value
)

var (
	dll *dylib.LazyDLL

	iupOpen *dylib.LazyProc
	//iupControlsOpen       *dylib.LazyProc
	iupAppend             *dylib.LazyProc
	iupDestroy            *dylib.LazyProc
	iupClose              *dylib.LazyProc
	iupSetGlobal          *dylib.LazyProc
	iupShow               *dylib.LazyProc
	iupPopup              *dylib.LazyProc
	iupLoad               *dylib.LazyProc
	iupMainLoop           *dylib.LazyProc
	iupSetCallback        *dylib.LazyProc
	iupSetHandle          *dylib.LazyProc
	iupSetAttributeHandle *dylib.LazyProc
	iupGetAttribute       *dylib.LazyProc
	iupGetInt             *dylib.LazyProc
	iupGetHandle          *dylib.LazyProc
	iupGetDialogChild     *dylib.LazyProc
	iupSetStrAttribute    *dylib.LazyProc
	iupSetAttributes      *dylib.LazyProc
	iupSetFocus           *dylib.LazyProc

	iupMessage *dylib.LazyProc

	iupHbox *dylib.LazyProc
	iupVbox *dylib.LazyProc

	iupDialog     *dylib.LazyProc
	iupButton     *dylib.LazyProc
	iupDatePick   *dylib.LazyProc
	iupFileDlg    *dylib.LazyProc
	iupFlatButton *dylib.LazyProc
	iupItem       *dylib.LazyProc
	iupLabel      *dylib.LazyProc
	iupList       *dylib.LazyProc
	iupMenu       *dylib.LazyProc
	iupSeparator  *dylib.LazyProc
	iupSubmenu    *dylib.LazyProc
	iupText       *dylib.LazyProc
	iupTabs       *dylib.LazyProc
	iupToggle     *dylib.LazyProc
	iupTree       *dylib.LazyProc
)

func Open() {
	initApi()
	iupOpen.Call(0, 0)
	SetGlobal("UTF8MODE", "YES")
	iupcontrols.Open()
}
func Append(h Ihandle, ch Ihandle) {
	iupAppend.Call(uintptr(h), uintptr(ch))
}
func Destroy(h Ihandle) {
	iupDestroy.Call(uintptr(h))
}
func Close() {
	iupClose.Call()
}
func SetGlobal(p string, v string) {
	iupSetGlobal.Call(uintptr(unsafe.Pointer(C.CString(p))), uintptr(unsafe.Pointer(C.CString(v))))
}
func MainLoop() {
	iupMainLoop.Call()
}
func Show(h Ihandle) {
	iupShow.Call(uintptr(h))
}
func Popup(h Ihandle, x int, y int) {
	iupPopup.Call(uintptr(h), uintptr(x), uintptr(y))
}

func Load(s string) {
	iupLoad.Call(uintptr(uintptr(unsafe.Pointer(C.CString(s)))))
}

func SetCallback(h Ihandle, name string, fn interface{}) {
	c_name := uintptr(unsafe.Pointer(C.CString(name)))
	switch fn.(type) {
	case uintptr:
		iupSetCallback.Call(uintptr(h), c_name, fn.(uintptr))
	default:
		iupSetCallback.Call(uintptr(h), c_name, syscall.NewCallbackCDecl(fn))
	}
}
func SetHandle(name string, h Ihandle) {
	iupSetHandle.Call(uintptr(unsafe.Pointer(C.CString(name))), uintptr(h))
}
func SetAttributeHandle(h Ihandle, name string, ch Ihandle) {
	iupSetAttributeHandle.Call(uintptr(h), uintptr(unsafe.Pointer(C.CString(name))), uintptr(ch))
}
func GetAttribute(h Ihandle, name string) string {
	c_name := uintptr(unsafe.Pointer(C.CString(name)))
	r, _, _ := iupGetAttribute.Call(uintptr(h), c_name)
	c := (*C.char)(unsafe.Pointer(r))
	return C.GoString(c)
}

func GetChild(h Ihandle, value string) Ihandle {
	c_val := uintptr(unsafe.Pointer(C.CString(value)))
	r, _, _ := iupGetDialogChild.Call(uintptr(h), c_val)
	return Ihandle(r)
}

func GetInt(h Ihandle, name string) int {
	c_name := uintptr(unsafe.Pointer(C.CString(name)))
	r, _, _ := iupGetInt.Call(uintptr(h), c_name)
	return int(r)
}

func GetHandle(value string) Ihandle {
	c_val := uintptr(unsafe.Pointer(C.CString(value)))
	r, _, _ := iupGetHandle.Call(c_val)
	return Ihandle(r)
}
func SetAttribute(h Ihandle, name string, value string) Ihandle {
	c_name := uintptr(unsafe.Pointer(C.CString(name)))
	c_val := uintptr(unsafe.Pointer(C.CString(value)))
	r, _, _ := iupSetStrAttribute.Call(uintptr(h), c_name, c_val)
	return Ihandle(r)
}
func SetAttributes(h Ihandle, value string) Ihandle {
	c_val := uintptr(unsafe.Pointer(C.CString(value)))
	r, _, _ := iupSetAttributes.Call(uintptr(h), c_val)
	return Ihandle(r)
}
func SetFocus(h Ihandle) {
	iupSetFocus.Call(uintptr(h))
}

func Message(t string, s string) {
	iupMessage.Call(uintptr(unsafe.Pointer(C.CString(t))), uintptr(unsafe.Pointer(C.CString(s))))
}

func Menu(childs ...Ihandle) Ihandle {
	var arr []uintptr
	arr = []uintptr{}
	for _, el := range childs {
		arr = append(arr, uintptr(el))
	}
	arr = append(arr, 0)
	p1, _, _ := iupMenu.Call(arr...)
	return Ihandle(p1)
}
func Hbox(childs ...Ihandle) Ihandle {
	var arr []uintptr
	arr = []uintptr{}
	for _, el := range childs {
		arr = append(arr, uintptr(el))
	}
	arr = append(arr, 0)
	p1, _, _ := iupHbox.Call(arr...)
	return Ihandle(p1)
}
func Vbox(childs ...Ihandle) Ihandle {
	var arr []uintptr
	arr = []uintptr{}
	for _, el := range childs {
		arr = append(arr, uintptr(el))
	}
	arr = append(arr, 0)
	p1, _, _ := iupVbox.Call(arr...)
	return Ihandle(p1)
}
func Tabs(childs ...Ihandle) Ihandle {
	var arr []uintptr
	arr = []uintptr{}
	for _, el := range childs {
		arr = append(arr, uintptr(el))
	}
	arr = append(arr, 0)
	p1, _, _ := iupTabs.Call(arr...)
	return Ihandle(p1)
}

func Dialog(h Ihandle) Ihandle {
	p1, _, _ := iupDialog.Call(uintptr(h))
	return Ihandle(p1)
}

func Button(s string) Ihandle {
	p1, _, _ := iupButton.Call(uintptr(unsafe.Pointer(C.CString(s))), uintptr(unsafe.Pointer(C.CString(""))))
	return Ihandle(p1)
}
func FlatButton(s string) Ihandle {
	p1, _, _ := iupFlatButton.Call(uintptr(unsafe.Pointer(C.CString(s))))
	return Ihandle(p1)
}
func FileDlg() Ihandle {
	p1, _, _ := iupFileDlg.Call()
	return Ihandle(p1)
}

func DatePick() Ihandle {
	p1, _, _ := iupDatePick.Call()
	return Ihandle(p1)
}

func Item(s string, a string) Ihandle {
	p1, _, _ := iupItem.Call(uintptr(unsafe.Pointer(C.CString(s))), uintptr(unsafe.Pointer(C.CString(a))))
	return Ihandle(p1)
}
func Label(s string) Ihandle {
	p1, _, _ := iupLabel.Call(uintptr(unsafe.Pointer(C.CString(s))))
	return Ihandle(p1)
}
func List(s string) Ihandle {
	p1, _, _ := iupList.Call(uintptr(unsafe.Pointer(C.CString(s))))
	return Ihandle(p1)
}
func Separator() Ihandle {
	p1, _, _ := iupSeparator.Call()
	return Ihandle(p1)
}
func Submenu(s string, ch Ihandle) Ihandle {
	p1, _, _ := iupSubmenu.Call(uintptr(unsafe.Pointer(C.CString(s))), uintptr(ch))
	return Ihandle(p1)
}
func Text(s string) Ihandle {
	p1, _, _ := iupText.Call(uintptr(unsafe.Pointer(C.CString(s))))
	return Ihandle(p1)
}
func Toggle(s string, a string) Ihandle {
	p1, _, _ := iupToggle.Call(uintptr(unsafe.Pointer(C.CString(s))), uintptr(unsafe.Pointer(C.CString(a))))
	return Ihandle(p1)
}
func Tree() Ihandle {
	p1, _, _ := iupTree.Call()
	return Ihandle(p1)
}

func initApi() {
	dll = dylib.NewLazyDLL("iup.dll")
	iupOpen = dll.NewProc("IupOpen")
	iupClose = dll.NewProc("IupClose")
	iupAppend = dll.NewProc("IupAppend")
	iupDestroy = dll.NewProc("IupDesttoy")
	iupSetGlobal = dll.NewProc("IupSetGlobal")
	iupShow = dll.NewProc("IupShow")
	iupPopup = dll.NewProc("IupPopup")
	iupLoad = dll.NewProc("IupLoad")
	iupMainLoop = dll.NewProc("IupMainLoop")
	iupGetAttribute = dll.NewProc("IupGetAttribute")
	iupGetHandle = dll.NewProc("IupGetHandle")
	iupGetInt = dll.NewProc("IupGetInt")
	iupGetDialogChild = dll.NewProc("IupGetDialogChild")
	iupSetCallback = dll.NewProc("IupSetCallback")
	iupSetHandle = dll.NewProc("IupSetHandle")
	iupSetAttributeHandle = dll.NewProc("IupSetAttributeHandle")
	iupSetStrAttribute = dll.NewProc("IupSetStrAttribute")
	iupSetAttributes = dll.NewProc("IupSetAttributes")
	iupSetFocus = dll.NewProc("IupSetFocus")

	iupMessage = dll.NewProc("IupMessage")

	iupHbox = dll.NewProc("IupHbox")
	iupVbox = dll.NewProc("IupVbox")
	iupMenu = dll.NewProc("IupMenu")
	iupTabs = dll.NewProc("IupTabs")

	iupDialog = dll.NewProc("IupDialog")
	iupButton = dll.NewProc("IupButton")
	iupFileDlg = dll.NewProc("IupFileDlg")
	iupFlatButton = dll.NewProc("IupFlatButton")
	iupDatePick = dll.NewProc("IupDatePick")
	iupItem = dll.NewProc("IupItem")
	iupLabel = dll.NewProc("IupLabel")
	iupList = dll.NewProc("IupList")
	iupSeparator = dll.NewProc("IupSeparator")
	iupSubmenu = dll.NewProc("IupSubmenu")
	iupText = dll.NewProc("IupText")
	iupToggle = dll.NewProc("IupToggle")
	iupTree = dll.NewProc("IupTree")

}
func (h Ihandle) GetAttribute(name string) string {
	return GetAttribute(h, name)
}

func (h Ihandle) GetChild(name string) Ihandle {
	return GetChild(h, name)
}

func (h Ihandle) SetAttribute(name string, value string) Ihandle {
	return SetAttribute(h, name, value)
}

func (h Ihandle) SetAttributes(name string) Ihandle {
	return SetAttributes(h, name)
}
func (h Ihandle) SetCallback(name string, fn interface{}) Ihandle {
	SetCallback(h, name, fn)
	return h
}
func (h Ihandle) SetAttributeHandle(name string, ch Ihandle) {
	SetAttributeHandle(h, name, ch)
}
