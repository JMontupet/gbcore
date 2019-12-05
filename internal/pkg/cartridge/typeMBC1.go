package cartridge

import (
	"fmt"
	"log"
)

const romBankSizeInt uint = 0x4000 // 16KB
const ramBankSizeInt uint = 0x2000 // 8KB

type mbc1 struct {
	data      []uint8
	romBank   uint
	nbROMBank uint

	ram       []uint8
	ramBank   uint
	nbRAMBank uint

	modeRam bool

	ramEnable bool
}

func (c *mbc1) Read(addr uint16) uint8 {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // ROM CART FIXED
		return c.data[addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // ROM CART BANK N
		return c.data[uint(addr)-0x4000+c.romBank*romBankSizeInt]
	case addr >= 0xA000 && addr <= 0xBFFF: // CART RAM
		if !c.ramEnable {
			return 0x00
		}
		return c.ram[uint(addr)-0xA000+c.ramBank*ramBankSizeInt]
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func (c *mbc1) Write(addr uint16, value uint8) {
	switch {
	// BANK CONTROLLER

	// 0000-1FFF - RAM Enable
	case addr >= 0x0000 && addr <= 0x1FFF:
		c.ramEnable = value&0xF == 0x0A

	// 2000-3FFF - ROM Bank Number (Write Only)
	case addr >= 0x2000 && addr <= 0x3FFF:
		c.changeROMBank(value)

	// 4000-5FFF - RAM Bank Number - or - Upper Bits of ROM Bank Number (Write Only)
	case addr >= 0x4000 && addr <= 0x5FFF:
		if c.modeRam {
			panic(fmt.Sprintf("TODO : 4000-5FFF - RAM Bank Number - or - Upper Bits of ROM Bank Number (Write Only) : %x", addr))
		}
		upperBits := uint((value & 0x03) << 5)
		c.changeROMBank(uint8(c.romBank&0x1F | upperBits))

	// 6000-7FFF - ROM/RAM Mode Select (Write Only)
	// 00h = ROM Banking Mode (up to 8KByte RAM, 2MByte ROM) (default)
	// 01h = RAM Banking Mode (up to 32KByte RAM, 512KByte ROM)
	case addr >= 0x6000 && addr <= 0x7FFF:
		// fmt.Printf("CHANGE CARTRIDGE MODE TO 0x%02X\n", value)
		c.modeRam = value&0x0A == 0x0A

	// CART RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		if c.ramEnable {
			c.ram[addr-0xA000] = value
		}

	// OFF RANGE
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func (c *mbc1) changeROMBank(v uint8) {
	value := uint(v)
	if value > c.nbROMBank {
		log.Fatalf(fmt.Sprintf("ROM BANK 0x%02X DOES NOT EXIST. MAX : 0x%02X", value, c.nbROMBank))
	}
	if value == 0 || value == 0x20 || value == 0x40 || value == 0x60 {
		value++
	}
	// if value != c.romBank {
	// 	fmt.Printf("CHANGE CARTRIDGE ROM BANK TO 0x%02X\n", value)
	// }
	c.romBank = value
}

func newMBC1(data []byte) Cartridge {
	cartridge := &mbc1{
		data:    data,
		romBank: 1,
	}
	switch ramSize := ReadRAMSize(cartridge); ramSize {
	case 0x00: // 00h - None
		cartridge.nbRAMBank = 0
	case 0x02: // 02h - 8 Kbytes
		cartridge.nbRAMBank = 4
		cartridge.ram = make([]uint8, 1024*8)
	default:
		panic(fmt.Sprintf("RAM SIZE NOT MANAGED FOR MBC1 : 0x%02X", ramSize))
	}

	switch romSize := ReadROMSize(cartridge); romSize {
	case 0x00: // 00h - None
		cartridge.nbROMBank = 0
	case 0x01: // 01h -		64KByte (4 banks)
		cartridge.nbROMBank = 4
	case 0x03: // 03h - 256KByte (16 banks)
		cartridge.nbROMBank = 16
	case 0x04: // 04h - 512KByte (32 banks)
		cartridge.nbROMBank = 32
	default:
		panic(fmt.Sprintf("ROM SIZE NOT MANAGED FOR MBC1 : 0x%02X", romSize))
	}

	return cartridge
}
