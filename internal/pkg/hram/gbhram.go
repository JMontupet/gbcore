package hram

const (
	// AddrStart is the first address of IOPorts
	AddrStart uint16 = 0xFF80
	// AddrEnd is the last address of IOPorts
	AddrEnd uint16 = 0xFFFE
)

// HRAM emulate Gameboy hram.
type HRAM struct {
	_data [0x7F]uint8
}

func (hram *HRAM) Read(addr uint16) uint8 {
	return hram._data[addr-AddrStart]
}

func (hram *HRAM) Write(addr uint16, value uint8) {
	hram._data[addr-AddrStart] = value
}

// NewGBHRAM returns new HRAM implementation
func NewGBHRAM() *HRAM {
	return &HRAM{}
}
