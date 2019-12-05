package ioports

type BitPtr struct {
	memory *IOPorts
	addr   uint16
	mask   uint8
}

func (ptr *BitPtr) Get() bool { return ptr.memory.Read(ptr.addr)&ptr.mask != 0 }
func (ptr *BitPtr) Set(set bool) {
	if set {
		ptr.memory.Write(ptr.addr, ptr.memory.Read(ptr.addr)|ptr.mask)
	} else {
		ptr.memory.Write(ptr.addr, ptr.memory.Read(ptr.addr) & ^ptr.mask)
	}
}
