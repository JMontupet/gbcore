package gpu

import (
	"github.com/jmontupet/gbcore/internal/pkg/mmu/memorymap"
)

// DMA INFO :
// - 40 sprites
// - 4 bytes / sprite
// - 40 * 4 = 160 bytes to transfer
// - 1 byte per cycle

const oamOffset uint16 = memorymap.OAMStart

type oam struct {
	_sprites [40]Sprite
}

func (oam *oam) GetSprite(id uint8) *Sprite {
	return &oam._sprites[id]
}

func (*oam) Read(_ uint16) uint8 {
	// No reason to read sprites info, but if needed :
	panic("TODO : READ SPRITE DATA")
	// return oam.Data[addr-oamOffset]
}

func (oam *oam) Write(addr uint16, value uint8) {
	oam.internalWrite(addr, value)
}

func (oam *oam) internalWrite(addr uint16, value uint8) {
	spriteID := (addr - oamOffset) >> 2
	spriteByte := (addr - oamOffset) & 3
	switch spriteByte {
	case 0:
		oam._sprites[spriteID].Y = value
	case 1:
		oam._sprites[spriteID].X = value
	case 2:
		oam._sprites[spriteID].TileID = value
	case 3:
		oam._sprites[spriteID].ObjToBgPriority = value&0x80 != 0
		oam._sprites[spriteID].YFlip = value&0x40 != 0
		oam._sprites[spriteID].XFlip = value&0x20 != 0
		oam._sprites[spriteID].PaletteNumber = value >> 4 & 0x01
		oam._sprites[spriteID].BankNumber = value >> 3 & 0x01
		oam._sprites[spriteID].ColorPalette = value & 0x03
		oam._sprites[spriteID].upToDate = false
	}
}

func newOAM() oam {
	oamRes := oam{}
	return oamRes
}

type Sprite struct {
	upToDate bool

	X      uint8
	Y      uint8
	TileID uint8

	// Bit7   OBJ-to-BG Priority 					(0=OBJ Above BG, 1=OBJ Behind BG color 1-3)
	// 												(Used for both BG and Window. BG color 0 is always behind OBJ)
	ObjToBgPriority bool
	// Bit6   Y flip          						(0=Normal, 1=Vertically mirrored)
	YFlip bool
	// Bit5   X flip          						(0=Normal, 1=Horizontally mirrored)
	XFlip bool
	// Bit4   Palette number  **Non CGB Mode Only** (0=OBP0, 1=OBP1)
	PaletteNumber uint8

	// *** CGB MODE ***
	// Bit3   Tile VRAM-Bank  **CGB Mode Only**     (0=Bank 0, 1=Bank 1)
	BankNumber uint8
	// Bit2-0 Palette number  **CGB Mode Only**     (OBP0-7)
	ColorPalette uint8
	// ***/CGB MODE ***
}

func (t *Sprite) GetPixels(vram *gbVRAM, line uint8, offset uint8, mode8x16 bool) []uint8 {
	tileID := t.TileID
	if t.YFlip {
		if mode8x16 {
			line = uint8((int16(line) - 15) * -1)
		} else {
			line = uint8((int16(line) - 7) * -1)
		}
	}

	if mode8x16 {
		if line > 7 {
			tileID |= 0x01
			line %= 8
		} else {
			tileID &= 0xFE
		}
	}
	var buff [8]uint8
	tilePixels := vram.
		GetTileData(true, tileID, t.BankNumber).
		appendPixelsLine(buff[:0], line, offset, 8, 0, false, false)
	pixels := make([]uint8, len(tilePixels))
	copy(pixels, tilePixels)
	if t.XFlip {
		for i := len(pixels)/2 - 1; i >= 0; i-- {
			opp := len(pixels) - 1 - i
			pixels[i], pixels[opp] = pixels[opp], pixels[i]
		}
	}
	return pixels
}
