package cartridge

import (
	"fmt"
	"log"
)

type romOnly struct {
	data []uint8
}

func (c *romOnly) Read(addr uint16) uint8 {
	switch {
	case addr >= 0x0000 && addr < 0x8000: // ROM CART
		return c.data[addr]
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func (*romOnly) Write(addr uint16, value uint8) {
	switch {
	case addr >= 0x0000 && addr < 0x8000: // ROM CART
		fmt.Printf("CART ROM IS READ ONLY !!! %X : %X\n", addr, value)
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func newROMOnly(data []uint8) Cartridge {
	return &romOnly{
		data: data,
	}
}
