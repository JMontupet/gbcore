package ioports

type Ptr struct {
	io   *IOPorts
	addr uint16
}

func (ptr *Ptr) Get() uint8 {
	return ptr.io.Read(ptr.addr)
}
func (ptr *Ptr) Set(value uint8) {
	ptr.io.Write(ptr.addr, value)
}
func (ptr *Ptr) getBitN(mask uint8) bool { return ptr.io.Read(ptr.addr)&mask != 0 }

func (ptr *Ptr) GetBit0() bool { return ptr.getBitN(0x01) }
func (ptr *Ptr) GetBit1() bool { return ptr.getBitN(0x02) }
func (ptr *Ptr) GetBit2() bool { return ptr.getBitN(0x04) }
func (ptr *Ptr) GetBit3() bool { return ptr.getBitN(0x08) }
func (ptr *Ptr) GetBit4() bool { return ptr.getBitN(0x10) }
func (ptr *Ptr) GetBit5() bool { return ptr.getBitN(0x20) }
func (ptr *Ptr) GetBit6() bool { return ptr.getBitN(0x40) }
func (ptr *Ptr) GetBit7() bool { return ptr.getBitN(0x80) }

func (ptr *Ptr) setBitN(set bool, mask uint8) {
	if set {
		ptr.io.Write(ptr.addr, ptr.io.Read(ptr.addr)|mask)
	} else {
		ptr.io.Write(ptr.addr, ptr.io.Read(ptr.addr) & ^mask)
	}
}
func (ptr *Ptr) SetBit0(set bool) { ptr.setBitN(set, 0x01) }
func (ptr *Ptr) SetBit1(set bool) { ptr.setBitN(set, 0x02) }
func (ptr *Ptr) SetBit2(set bool) { ptr.setBitN(set, 0x04) }
func (ptr *Ptr) SetBit3(set bool) { ptr.setBitN(set, 0x08) }
func (ptr *Ptr) SetBit4(set bool) { ptr.setBitN(set, 0x10) }
func (ptr *Ptr) SetBit5(set bool) { ptr.setBitN(set, 0x20) }
func (ptr *Ptr) SetBit6(set bool) { ptr.setBitN(set, 0x40) }
func (ptr *Ptr) SetBit7(set bool) { ptr.setBitN(set, 0x80) }
