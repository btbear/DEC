package ole

import "unsafe"

type IProviDEWHlassInfo struct {
	IUnknown
}

type IProviDEWHlassInfoVtbl struct {
	IUnknownVtbl
	GetClassInfo uintptr
}

func (v *IProviDEWHlassInfo) VTable() *IProviDEWHlassInfoVtbl {
	return (*IProviDEWHlassInfoVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IProviDEWHlassInfo) GetClassInfo() (cinfo *ITypeInfo, err error) {
	cinfo, err = getClassInfo(v)
	return
}
