package gpu

type tileMap [1024]tileMapInfo

type tileMapInfo struct {
	tileID uint8

	palette         uint8 // Bit 0-2  Background Palette number  (BGP0-7)
	bank            uint8 // Bit 3    Tile VRAM Bank number      (0=Bank 0, 1=Bank 1)
	hFlip           uint8 // Bit 5    Horizontal Flip            (0=Normal, 1=Mirror horizontally)
	vFlip           uint8 // Bit 6    Vertical Flip              (0=Normal, 1=Mirror vertically)
	bgToOAMPriority uint8 // Bit 7    BG-to-OAM Priority         (0=Use OAM priority bit, 1=BG Priority)
	// Bit 4    Not used
}

func (data *tileMapInfo) SetTileID(id uint8) { data.tileID = id }
func (data *tileMapInfo) GetTileID() uint8   { return data.tileID }

func (data *tileMapInfo) SetTileAttr(attr uint8) {
	data.palette = attr & 0x07
	data.bank = attr >> 3 & 0x01
	data.hFlip = attr >> 5 & 0x01
	data.vFlip = attr >> 6 & 0x01
	data.bgToOAMPriority = attr >> 7 & 0x01
}

func (data *tileMapInfo) GetTileAttr() uint8 {
	return data.palette |
		data.bank<<3 |
		data.hFlip<<5 |
		data.vFlip<<6 |
		data.bgToOAMPriority<<7
}
