package cartridge

import (
	"fmt"

	"github.com/jmontupet/gbcore/internal/pkg/memory"
)

// Cartridge emulate GameBoy cartrigde. Different memory controller can be implemented.
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
		return newMBC3(data), nil
	case 0x13: // ROM_MBC3_RAM_Batt
		return newMBC3(data), nil
	case 0x19: // ROM_MBC5
		return newMBC5(data), nil
	case 0x1B: // ROM_MBC5_RAM_Batt
		return newMBC5(data), nil
	default:
		return nil, fmt.Errorf("CARTRIDGE TYPE NOT IMPLEMENTED : 0x%02X", cType)
	}
}

// const (
// 	ROM_Only                 CartridgeType = 0x00
// 	ROM_MBC1                 CartridgeType = 0x01
// 	ROM_MBC1_RAM             CartridgeType = 0x02
// 	ROM_MBC1_RAM_Batt        CartridgeType = 0x03
// 	ROM_MBC                  CartridgeType = 0x05
// 	ROM_MBC2_Batt            CartridgeType = 0x06
// 	ROM_RAM                  CartridgeType = 0x08
// 	ROM_RAM_Batt             CartridgeType = 0x09
// 	ROM_MMM01                CartridgeType = 0x0B
// 	ROM_MMM01_RAM            CartridgeType = 0x0C
// 	ROM_MMM01_RAM_Batt       CartridgeType = 0x0D
// 	ROM_MBC3_Timer_Batt      CartridgeType = 0x0F
// 	ROM_MBC3_Timer_RAM_Batt  CartridgeType = 0x10
// 	ROM_MBC3                 CartridgeType = 0x11
// 	ROM_MBC3_RAM             CartridgeType = 0x12
// 	ROM_MBC3_RAM_Batt        CartridgeType = 0x13
// 	ROM_MBC5                 CartridgeType = 0x19
// 	ROM_MBC5_RAM             CartridgeType = 0x1A
// 	ROM_MBC5_RAM_Batt        CartridgeType = 0x1B
// 	ROM_MBC5_Rumble          CartridgeType = 0x1C
// 	ROM_MBC5_Rumble_RAM      CartridgeType = 0x1D
// 	ROM_MBC5_Rumble_RAM_Batt CartridgeType = 0x1E
// 	Pocket_Camera            CartridgeType = 0x1F
// 	Bandai_TAMA5             CartridgeType = 0xFD
// 	Hudson_HuC_3             CartridgeType = 0xFE
// 	Hudson_HuC_1             CartridgeType = 0xFF
// )
