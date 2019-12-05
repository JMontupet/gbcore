package unusableaddr

// To bind to MMU
const (
	// AddrStart is the first address of Unused Memory
	AddrStart uint16 = 0xFEA0
	// AddrEnd is the last address of Unused Memory
	AddrEnd uint16 = 0xFEFF
)

// UnusableAddr emulate unused address range.
//
// Not used but present.
// Any game uses this range ?
type UnusableAddr struct {
	_data [0x60]uint8
}

func (unused *UnusableAddr) Read(addr uint16) uint8 {
	return unused._data[addr-AddrStart]
}

func (unused *UnusableAddr) Write(addr uint16, value uint8) {
	unused._data[addr-AddrStart] = value
}

// NewUnusableAddr returns simple Memory implementation for Unused memory range.
func NewUnusableAddr() *UnusableAddr {
	return &UnusableAddr{}
}
