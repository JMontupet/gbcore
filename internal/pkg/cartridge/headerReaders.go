package cartridge

import (
	"strings"
)

func ReadTitle(c Cartridge) string {
	s := make([]uint8, 0, 0xF)
	for i := uint16(0x134); i < 0x143; i++ {
		s = append(s, c.Read(i))
	}
	return strings.Trim(string(s), " ")
}

func ReadCGBCompatible(c Cartridge) bool { return c.Read(0x143) == 0x80 || c.Read(0x143) == 0xC0 }
func ReadType(c Cartridge) uint8         { return c.Read(0x147) }
func ReadROMSize(c Cartridge) uint8      { return c.Read(0x148) }
func ReadRAMSize(c Cartridge) uint8      { return c.Read(0x149) }
