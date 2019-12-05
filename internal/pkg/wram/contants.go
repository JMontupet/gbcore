package wram

// To bind to MMU
const (
	// WRamStart is the WRAM Memory Start Addr
	WRamStart uint16 = 0xC000
	// WRamEnd is the WRAM Memory End Addr
	WRamEnd uint16 = 0xDFFF

	// MirrorWramStart is the Mirror of C000~DDFF Start Addr
	MirrorWramStart uint16 = 0xE000
	// MirrorWramEnd is the Mirror of C000~DDFF Start Addr
	MirrorWramEnd uint16 = 0xFDFF
)

// Internals
const (
	// Fixed RAM Memory End Addr
	fixedEnd uint16 = 0xCFFF
	// Banked RAM Memory Start Addr
	bankedStart uint16 = 0xD000
	// WRamEnd is the WRAM Memory End Addr
)
