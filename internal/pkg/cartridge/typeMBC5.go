package cartridge

import (
	"fmt"
	"log"
)

type mbc5 struct {
	romBanks [][]uint8
	romBank  int

	ramBanks [][]uint8
	ramBank  int

	ramEnable bool
}

func (c *mbc5) Read(addr uint16) uint8 {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // ROM CART FIXED
		return c.romBanks[0][addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // ROM CART BANK N
		return c.romBanks[c.romBank][addr-0x4000]
	case addr >= 0xA000 && addr <= 0xBFFF: // CART RAM
		if c.ramEnable {
			return c.ramBanks[c.ramBank][addr-0xA000]
		}
		return 0x00
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func (c *mbc5) Write(addr uint16, value uint8) {
	switch {
	// BANK CONTROLLER

	// 0000-1FFF - RAM and Timer Enable (Write Only)
	case addr >= 0x0000 && addr <= 0x1FFF:
		c.ramEnable = value&0xF == 0x0A

	// 2000-3FFF - ROM Bank Number (Write Only)
	case addr >= 0x2000 && addr <= 0x3FFF:
		c.changeROMBank(value)

	// 4000-5FFF - RAM Bank Number - or - RTC Register Select (Write Only)
	case addr >= 0x4000 && addr <= 0x5FFF:
		if value < 0x10 { // 16 banks max
			c.ramBank = int(value)
		} else {
			log.Fatalf("RAM BANK INVALID : 0x%02X", addr)
		}

	// 6000-7FFF - Latch Clock Data (Write Only)
	case addr >= 0x6000 && addr <= 0x7FFF:
		fmt.Printf("TODO : Latch Clock Data %X : %X\n", addr, value)

	// CART RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		if c.ramEnable {
			c.ramBanks[c.ramBank][addr-0xA000] = value
		}
	// OFF RANGE
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func (c *mbc5) changeROMBank(v uint8) {
	value := int(v & 0x7F)
	if value > len(c.romBanks) {
		panic(fmt.Sprintf("ROM BANK 0x%02X DOES NOT EXIST. MAX : 0x%02X", value, len(c.romBanks)))
	}
	if value == 0 {
		c.romBank = 1
		return
	}
	c.romBank = value & 0x7F
}

func newMBC5(data []byte) (Cartridge, error) {
	cartridge := &mbc5{
		romBank: 1,
	}

	romBanks, err := splitROMBanks(data[0x148], data)
	if err != nil {
		return nil, err
	}
	cartridge.romBanks = romBanks

	ramBanks, err := makeRAMBanks(ReadRAMSize(cartridge))
	if err != nil {
		return nil, err
	}
	cartridge.ramBanks = ramBanks

	return cartridge, nil
}
