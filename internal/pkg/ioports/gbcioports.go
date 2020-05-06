package ioports

const (
	// AddrStart is the first address of IOPorts
	AddrStart uint16 = 0xFF00
	// AddrEnd is the last address of IOPorts
	AddrEnd uint16 = 0xFF7F
)

// IOPorts emulate Gameboy Color IO ports
type IOPorts struct {
	_data [0x80]uint8
}

func (io *IOPorts) Read(addr uint16) uint8 {
	return io._data[addr-AddrStart]
}

func (io *IOPorts) Write(addr uint16, value uint8) {
	io._data[addr-AddrStart] = value
}

func (io *IOPorts) NewPtr(addr uint16) *Ptr {
	return &Ptr{
		addr: addr,
		io:   io,
	}
}

func (io *IOPorts) NewMaskedPtr(addr uint16, mask uint8) *MaskedPtr {
	return &MaskedPtr{
		ptr: Ptr{
			addr: addr,
			io:   io,
		},
		mask: mask,
	}
}

func (io *IOPorts) NewBit0Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x01} }
func (io *IOPorts) NewBit1Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x02} }
func (io *IOPorts) NewBit2Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x04} }
func (io *IOPorts) NewBit3Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x08} }
func (io *IOPorts) NewBit4Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x10} }
func (io *IOPorts) NewBit5Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x20} }
func (io *IOPorts) NewBit6Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x40} }
func (io *IOPorts) NewBit7Ptr(addr uint16) *BitPtr { return &BitPtr{io, addr, 0x80} }

// NewGBIOPorts create new IOPorts, basically a simple Memory implementation.
func NewGBIOPorts() *IOPorts {
	return &IOPorts{}
}
