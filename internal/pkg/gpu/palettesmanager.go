package gpu

import (
	"log"

	"github.com/jmontupet/gbcore/pkg/coreio"
)

// FF68 & FF6A : Bit 0-5   Index (00-3F)
//               Bit 7     Auto Increment  (0=Disabled, 1=Increment after Writing)

type palettesManager struct {
	cgb bool

	bgPaletteIndex   uint8
	bgPaletteAutoInc bool
	bgPaletteData    [2 * 4 * 8]uint8 // 2 bytes * 4 colors * 8 palettes. Access via FF69

	spritePaletteIndex   uint8
	spritePaletteAutoInc bool
	spritePaletteData    [2 * 4 * 8]uint8 // 2 bytes * 4 colors * 8 palettes. Access via FF6B
}

func (pm *palettesManager) setPalettes(colors *coreio.FrameColors) {
	if pm.cgb {
		var palettePointer uint16 = 0b100000 * 3
		for i := 0; i < 64; i += 2 {
			var color = uint16(pm.bgPaletteData[i]) | uint16(pm.bgPaletteData[i+1])<<8
			colors[palettePointer] = (uint8(color) & 0x1F) << 3
			palettePointer++
			colors[palettePointer] = (uint8(color>>5) & 0x1F) << 3
			palettePointer++
			colors[palettePointer] = (uint8(color>>10) & 0x1F) << 3
			palettePointer++
		}

		palettePointer = 0b000000 * 3
		for i := 0; i < 64; i += 2 {
			var color = uint16(pm.spritePaletteData[i]) | uint16(pm.spritePaletteData[i+1])<<8
			colors[palettePointer] = (uint8(color) & 0x1F) << 3
			palettePointer++
			colors[palettePointer] = (uint8(color>>5) & 0x1F) << 3
			palettePointer++
			colors[palettePointer] = (uint8(color>>10) & 0x1F) << 3
			palettePointer++
		}
	}
	//  else {
	// 	panic("TODO : NON CGB PALETTES")
	// }
}

func (pm *palettesManager) Write(addr uint16, value uint8) {
	switch addr {
	case 0xFF68:
		pm.bgPaletteIndex = value & 0x3F
		pm.bgPaletteAutoInc = value>>7 == 1
	case 0xFF69:
		pm.bgPaletteData[pm.bgPaletteIndex] = value
		if pm.bgPaletteAutoInc {
			pm.bgPaletteIndex = (pm.bgPaletteIndex + 1) & 0x3F
		}
	case 0xFF6A:
		pm.spritePaletteIndex = value & 0x3F
		pm.spritePaletteAutoInc = value>>7 == 1
	case 0xFF6B:
		pm.spritePaletteData[pm.spritePaletteIndex] = value
		if pm.spritePaletteAutoInc {
			pm.spritePaletteIndex = (pm.spritePaletteIndex + 1) & 0x3F
		}
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func (pm *palettesManager) Read(addr uint16) uint8 {
	switch addr {
	case 0xFF68:
		if pm.bgPaletteAutoInc {
			return pm.bgPaletteIndex & 0x80
		}
		return pm.bgPaletteIndex
	case 0xFF69:
		return pm.bgPaletteData[pm.bgPaletteIndex]
	case 0xFF6A:
		if pm.spritePaletteAutoInc {
			return pm.spritePaletteIndex & 0x80
		}
		return pm.spritePaletteIndex
	case 0xFF6B:
		return pm.spritePaletteData[pm.spritePaletteIndex]
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func newPalettesManager(cgb bool) *palettesManager {
	pm := new(palettesManager)
	pm.cgb = cgb
	for i := range pm.bgPaletteData {
		pm.bgPaletteData[i] = 0xFF
		pm.spritePaletteData[i] = 0xFF
	}
	return pm
}
