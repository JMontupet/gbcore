package gpu

import (
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/mmu/memorymap"
)

const (
	vramOffset uint16 = memorymap.VRamStart // 0x8000

	tileData1Start uint16 = 0x8000
	tileData0Start uint16 = 0x8800
	tileData1End   uint16 = 0x8FFF
	tileData0End   uint16 = 0x97FF

	tileMap0Start uint16 = 0x9800
	tileMap0End   uint16 = 0x9BFF
	tileMap1Start uint16 = 0x9C00
	tileMap1End   uint16 = 0x9FFF

	// vramEnd uint16 = memorymap.VRamEnd // 0x9FFF

	// tileSize uint16 = 0x10
)

type gbVRAM struct {
	tileMaps [2]tileMap
	tiles    [2][384]tile
	bankFlag *ioports.MaskedPtr
}

func (vram *gbVRAM) Read(addr uint16) uint8 {
	switch {
	// Tile MAP 0
	case addr >= tileMap0Start && addr <= tileMap0End:
		if vram.bankFlag.Get() == 0 {
			return vram.tileMaps[0][addr-tileMap0Start].GetTileID()
		}
		return vram.tileMaps[0][addr-tileMap0Start].GetTileAttr()
	// Tile MAP 1
	case addr >= tileMap1Start && addr <= tileMap1End:
		if vram.bankFlag.Get() == 0 {
			return vram.tileMaps[1][addr-tileMap1Start].GetTileID()
		}
		return vram.tileMaps[1][addr-tileMap1Start].GetTileAttr()

	case addr >= tileData1Start && addr <= tileData0End:
		addr -= vramOffset
		return vram.tiles[vram.bankFlag.Get()][addr>>4][addr&15]
	default:
		log.Fatalf("GPU MEMORY UNREACHABLE : 0x%04X", addr)
		return 0
	}
}

func (vram *gbVRAM) Write(addr uint16, value uint8) {
	switch {
	// Tile MAP 0
	case addr >= tileMap0Start && addr <= tileMap0End:
		if vram.bankFlag.Get() == 0 {
			vram.tileMaps[0][addr-tileMap0Start].SetTileID(value)
		} else {
			vram.tileMaps[0][addr-tileMap0Start].SetTileAttr(value)
		}
	// Tile MAP 1
	case addr >= tileMap1Start && addr <= tileMap1End:
		if vram.bankFlag.Get() == 0 {
			vram.tileMaps[1][addr-tileMap1Start].SetTileID(value)
		} else {
			vram.tileMaps[1][addr-tileMap1Start].SetTileAttr(value)
		}

	case addr >= tileData1Start && addr <= tileData0End:
		addr -= vramOffset
		vram.tiles[vram.bankFlag.Get()][addr>>4][addr&15] = value
	default:
		log.Fatalf("GPU MEMORY UNREACHABLE : 0x%04X", addr)
	}
}
func (vram *gbVRAM) GetTileInfo(mapIndex uint8, tileX uint8, tileY uint8) tileMapInfo {
	return vram.tileMaps[mapIndex][uint16(tileY)*32+uint16(tileX)]
}
func (vram *gbVRAM) GetTileData(dataTable bool, rawID uint8, bank uint8) *tile {
	if dataTable {
		// unsigned start from 8000
		return &vram.tiles[bank][rawID]
	}
	// signed centered to 9000
	return &vram.tiles[bank][(uint16(0x9000+int32(int8(rawID))<<4)-vramOffset)>>4]
}

func newGBVRAM(io *ioports.IOPorts) gbVRAM {
	vram := gbVRAM{
		bankFlag: io.NewMaskedPtr(0xFF4F, 0x01),
	}
	return vram
}
