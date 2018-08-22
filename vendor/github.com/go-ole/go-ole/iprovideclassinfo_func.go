// +build !windows

package ole

func getClassInfo(disp *IProviDEWHlassInfo) (tinfo *ITypeInfo, err error) {
	return nil, NewError(E_NOTIMPL)
}
