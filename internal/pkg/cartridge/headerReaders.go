package cartridge

import (
	"strings"
)

// ReadNintendoLogo reads Nintendo Logo from its header.
func ReadNintendoLogo(c Cartridge) (data [48]byte) {
	var addrStart uint16 = 0x104
	for i := range data {
		data[i] = c.Read(addrStart + uint16(i))
	}
	return data
}

// ReadTitle reads cartridge title from its header.
func ReadTitle(c Cartridge) string {
	s := make([]uint8, 0, 15)
	for i := uint16(0x134); i < 0x143; i++ {
		s = append(s, c.Read(i))
	}
	return strings.Trim(string(s), " ")
}

// ReadCGBCompatible reads from headers if a cartridge is "GameBoy Color compatible"
func ReadCGBCompatible(c Cartridge) bool { return c.Read(0x143) == 0x80 || c.Read(0x143) == 0xC0 }

// ReadManufacturerCode reads Manufacturer Code from headers
func ReadManufacturerCode(c Cartridge) string {
	return string([]byte{
		c.Read(0x013F), c.Read(0x0140),
		c.Read(0x0141), c.Read(0x0142),
	})
}

// ReadCGBFlag reads CGB Flag from headers
// 	 80h - Game supports CGB functions, but works on old gameboys also.
//	 C0h - Game works on CGB only (physically the same as 80h).
// Values with Bit 7 set, and either Bit 2 or 3 set, will switch the gameboy
// into a special non-CGB-mode with uninitialized palettes.
func ReadCGBFlag(c Cartridge) uint8 { return c.Read(0x143) }

// ReadNewLicenseeCode reads New Licensee Code from headers
//	Specifies a two character ASCII licensee code.
func ReadNewLicenseeCode(c Cartridge) string {
	return string([]byte{c.Read(0x0144), c.Read(0x0145)})
}

// ReadSGBFlag reads SGB Flag from headers
// 	00h = No SGB functions (Normal Gameboy or CGB only game)
// 	03h = Game supports SGB functions
func ReadSGBFlag(c Cartridge) uint8 { return c.Read(0x146) }

// ReadType reads from headers which memory controller is used
func ReadType(c Cartridge) uint8 { return c.Read(0x147) }

// ReadCartridgeType reads from headers which memory controller is used
func ReadCartridgeType(c Cartridge) uint8 { return c.Read(0x147) }

// ReadROMSize reads from headers the rom size
func ReadROMSize(c Cartridge) uint8 { return c.Read(0x148) }

// ReadRAMSize reads from headers the ram size
func ReadRAMSize(c Cartridge) uint8 { return c.Read(0x149) }

// ReadDestinationCode reads Destination Code from headers
// 	00h - Japanese
// 	01h - Non-Japanese
func ReadDestinationCode(c Cartridge) uint8 { return c.Read(0x14A) }

// ReadOldLicenseeCode reads Old Licensee Code from headers
//  - A value of 33h signalizes that the New License Code in header bytes 0144-0145 is used instead.
// 	- Super GameBoy functions won't work if <> $33.
func ReadOldLicenseeCode(c Cartridge) uint8 { return c.Read(0x14B) }

// ReadROMVersion reads ROM Version from headers
func ReadROMVersion(c Cartridge) uint8 { return c.Read(0x14C) }

// ReadHeaderChecksum reads Header Checksum from headers
func ReadHeaderChecksum(c Cartridge) uint8 { return c.Read(0x14D) }

func CalcHeaderChecksum(c Cartridge) uint8 {
	//  x=0:FOR i=0134h TO 014Ch:x=x-MEM[i]-1:NEXT
	var cs uint8
	for i := uint16(0x134); i < 0x14D; i++ {
		cs = cs - c.Read(i) - 1
	}
	return cs
}

// ReadROMChecksum reads ROM Checksum from headers
func ReadROMChecksum(c Cartridge) uint16 {
	return uint16(c.Read(0x014E)) | (uint16(c.Read(0x014F)) << 8)
}
