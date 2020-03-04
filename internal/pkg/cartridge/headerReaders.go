package cartridge

import (
	"strings"
)

// ReadTitle reads cartridge title from its header.
func ReadTitle(c Cartridge) string {
	s := make([]uint8, 0, 0xF)
	for i := uint16(0x134); i < 0x143; i++ {
		s = append(s, c.Read(i))
	}
	return strings.Trim(string(s), " ")
}

// ReadCGBCompatible reads from headers if a cartridge is "GameBoy Color compatible"
func ReadCGBCompatible(c Cartridge) bool { return c.Read(0x143) == 0x80 || c.Read(0x143) == 0xC0 }

// ReadType reads from headers which memory controller is used
func ReadType(c Cartridge) uint8 { return c.Read(0x147) }

// ReadROMSize reads from headers the rom size
func ReadROMSize(c Cartridge) uint8 { return c.Read(0x148) }

// ReadRAMSize reads from headers the ram size
func ReadRAMSize(c Cartridge) uint8 { return c.Read(0x149) }
