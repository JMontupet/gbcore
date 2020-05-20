package cartridge

import (
	"fmt"

	"github.com/jmontupet/gbcore/internal/pkg/memory"
)

// Cartridge emulate GameBoy cartridge. Different memory controller can be implemented.
type Cartridge interface {
	memory.Memory
}

// NewCartridge creates a cartridge based on its header data
func NewCartridge(data []byte) (Cartridge, error) {
	if data == nil || len(data) < 0x0150 {
		return nil, fmt.Errorf("CARTRIDGE MIN SIZE IS 0x%04X, GOT 0x%04X", 0x0150, len(data))
	}
	switch cType := data[0x147]; cType {
	case 0x00: // ROM_Only
		return newROMOnly(data), nil
	case 0x01: // ROM_MBC1
		return newMBC1(data), nil
	case 0x03: // ROM_MBC1_RAM_Batt
		return newMBC1(data), nil
	case 0x10: // ROM_MBC3_Timer_RAM_Batt
		return newMBC3(data)
	case 0x13: // ROM_MBC3_RAM_Batt
		return newMBC3(data)
	case 0x19: // ROM_MBC5
		return newMBC5(data)
	case 0x1B: // ROM_MBC5_RAM_Batt
		return newMBC5(data)
	default:
		return nil, fmt.Errorf("CARTRIDGE TYPE NOT IMPLEMENTED : 0x%02X", cType)
	}
}

// 00h  ROM ONLY
// 01h  MBC1
// 02h  MBC1+RAM
// 03h  MBC1+RAM+BATTERY
// 05h  MBC2
// 06h  MBC2+BATTERY
// 08h  ROM+RAM
// 09h  ROM+RAM+BATTERY
// 0Bh  MMM01
// 0Ch  MMM01+RAM
// 0Dh  MMM01+RAM+BATTERY
// 0Fh  MBC3+TIMER+BATTERY
// 10h  MBC3+TIMER+RAM+BATTERY
// 11h  MBC3
// 12h  MBC3+RAM
// 13h  MBC3+RAM+BATTERY
// 19h  MBC5
// 1Ah  MBC5+RAM
// 1Bh  MBC5+RAM+BATTERY
// 1Ch  MBC5+RUMBLE
// 1Dh  MBC5+RUMBLE+RAM
// 1Eh  MBC5+RUMBLE+RAM+BATTERY
// 20h  MBC6
// 22h  MBC7+SENSOR+RUMBLE+RAM+BATTERY
// FCh  POCKET CAMERA
// FDh  BANDAI TAMA5
// FEh  HuC3
// FFh  HuC1+RAM+BATTERY

const romBankSize = 0x4000 // 16KB
const ramBankSize = 0x2000 // 8KB
