package ioports

type MaskedPtr struct {
	ptr  Ptr
	mask uint8
}

func (mptr *MaskedPtr) Get() uint8 {
	return mptr.ptr.Get() & mptr.mask
}
func (mptr *MaskedPtr) Set(value uint8) {
	// Remove old values and add new filtered values
	mptr.ptr.Set(mptr.ptr.Get() & ^mptr.mask | value&mptr.mask)
}
