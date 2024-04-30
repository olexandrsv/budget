package iupcontrols

import (
	"github.com/ying32/dylib"
)

/*
 */
import "C"

var (
	dll             *dylib.LazyDLL  = dylib.NewLazyDLL("iupcontrols.dll")
	iupControlsOpen *dylib.LazyProc = dll.NewProc("IupControlsOpen")
)

func Open(){
	iupControlsOpen.Call()
}
