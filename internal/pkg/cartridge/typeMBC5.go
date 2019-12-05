package cartridge

import (
	"fmt"
	"log"
)

type mbc5 struct {
	data      []uint8
	romBank   uint
	nbROMBank uint

	ram       []uint8
	ramBank   uint
	nbRAMBank uint

	ramTimerEnable bool

	rtcEnable bool
}

func (c *mbc5) Read(addr uint16) uint8 {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // ROM CART FIXED
		return c.data[addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // ROM CART BANK N
		return c.data[uint(addr)-0x4000+c.romBank*romBankSizeInt]
	case addr >= 0xA000 && addr <= 0xBFFF: // CART RAM
		if c.ramTimerEnable {
			return c.ram[uint(addr)-0xA000+c.ramBank*ramBankSizeInt]
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
		c.ramTimerEnable = value&0xF == 0x0A

	// 2000-3FFF - ROM Bank Number (Write Only)
	case addr >= 0x2000 && addr <= 0x3FFF:
		c.changeROMBank(value)

	// 4000-5FFF - RAM Bank Number - or - RTC Register Select (Write Only)
	case addr >= 0x4000 && addr <= 0x5FFF:
		if value <= 0x03 { // 4 banks max
			c.ramBank = uint(value)
			c.rtcEnable = false
		} else {
			c.rtcEnable = true
		}

	// 6000-7FFF - Latch Clock Data (Write Only)
	case addr >= 0x6000 && addr <= 0x7FFF:
		fmt.Printf("TODO : Latch Clock Data %X : %X\n", addr, value)

	// CART RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		if c.ramTimerEnable && !c.rtcEnable {
			c.ram[addr-0xA000] = value
		}
	// OFF RANGE
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func (c *mbc5) changeROMBank(v uint8) {
	value := uint(v & 0x7F)
	if value > c.nbROMBank {
		panic(fmt.Sprintf("ROM BANK 0x%02X DOES NOT EXIST. MAX : 0x%02X", value, c.nbROMBank))
	}
	if value == 0 {
		c.romBank = 1
		return
	}
	c.romBank = value & 0x7F
}

func newMBC5(data []byte) Cartridge {
	cartridge := &mbc5{
		data:    data,
		romBank: 1,
	}
	switch ramSize := ReadRAMSize(cartridge); ramSize {
	case 0x00: // 00h - None
		cartridge.nbRAMBank = 0
	case 0x02: // 02h - 8 Kbytes
		cartridge.nbRAMBank = 1
		cartridge.ram = make([]uint8, 1024*8)
	case 0x03: // 03h -		32 KBytes (4 banks of 8KBytes each)
		cartridge.nbRAMBank = 4
		cartridge.ram = make([]uint8, 4*1024*8)
	default:
		panic(fmt.Sprintf("RAM SIZE NOT MANAGED FOR MBC5 : 0x%02X", ramSize))
	}

	switch romSize := ReadROMSize(cartridge); romSize {

	case 0x05: // 05h -		1MByte (64 banks)
		cartridge.nbROMBank = 64
	case 0x06: // 06h -   2MByte (128 banks) - only 125 banks used by MBC1
		cartridge.nbROMBank = 128
	case 0x07: // 07h -   4MByte (256 banks)
		cartridge.nbROMBank = 256
	default:
		panic(fmt.Sprintf("ROM SIZE NOT MANAGED FOR MBC5 : 0x%02X", romSize))
	}

	return cartridge
}
