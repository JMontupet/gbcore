package wram

import (
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/ioports"
)

// WRam is the internal Working RAM of the Gameboy
type WRam struct {
	// _fixedRAM Contains Fixed Ram values
	_fixedRAM [1024 * 4]uint8
	// _bankedRAM Contains Banked Ram values
	_bankedRAM [7][1024 * 4]uint8

	// ff70 - SVBK - CGB Mode Only - WRAM Bank (bit 0-2)
	ff70 *ioports.MaskedPtr
}

func (io *WRam) getBank() uint8 {
	bank := io.ff70.Get()
	if bank != 0 {
		bank-- // Translate to array index
	}
	// Bank 0 is interpreted as bank 1. No translation needed
	return bank
}

func (io *WRam) Read(addr uint16) uint8 {
	switch {
	case addr >= WRamStart && addr <= fixedEnd:
		return io._fixedRAM[addr-WRamStart]
	case addr >= bankedStart && addr <= WRamEnd:
		return io._bankedRAM[io.getBank()][addr-bankedStart]
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0
	}
}

func (io *WRam) Write(addr uint16, value uint8) {
	switch {
	case addr >= WRamStart && addr <= fixedEnd:
		io._fixedRAM[addr-WRamStart] = value
	case addr >= bankedStart && addr <= WRamEnd:
		io._bankedRAM[io.getBank()][addr-bankedStart] = value
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

// NewWram create new WRam instance
//
// ioports is required to read current wram bank (FF70)
func NewWram(io *ioports.IOPorts) *WRam {
	return &WRam{
		ff70: io.NewMaskedPtr(0xFF70, 0x03),
	}
}
