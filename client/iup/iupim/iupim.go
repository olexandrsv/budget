package iupim

import (
	iup "budget/iup/wrapper"
	"unsafe"

	"github.com/ying32/dylib"
)

/*
 */
import "C"

var (
	dll          *dylib.LazyDLL  = dylib.NewLazyDLL("iupim.dll")
	iupLoadImage *dylib.LazyProc = dll.NewProc("IupLoadImage")
)

func LoadImage(s string) iup.Ihandle {
	p1, _, _ := iupLoadImage.Call(uintptr(unsafe.Pointer(C.CString(s))))
	return iup.Ihandle(p1)
}
