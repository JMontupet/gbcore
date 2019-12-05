package gpu

import (
	"log"

	"github.com/jmontupet/gbcore/pkg/coreio"

	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/mmu/memorymap"
)

const (
	nbFrameLines = 144
	nbFrameRow   = 160
	oamLength    = 80   // 77-83 (80)
	drawLength   = 172  // 169-175 (172)
	hBlankLength = 204  // 201-207 (204)
	vBlankLength = 4560 // 4560

	// Each Line
	oamEnd     = oamLength
	drawEnd    = oamEnd + drawLength
	hBlanckEnd = drawEnd + hBlankLength

	// Each Fram
	frameEnd  = nbFrameLines * hBlanckEnd
	vBlankEnd = frameEnd + vBlankLength
)

type gpuMode uint8

const frameBufferSize = nbFrameRow * nbFrameLines

const (
	ModeHBlank = iota
	ModeVBlank
	ModeOAM
	ModeVRAM
)

type GPU struct {
	// Color ?
	cgb bool

	// Palettes Manager
	palettesManager *palettesManager

	// Internal VRAM
	_vram gbVRAM

	// OAM Memory
	_oam oam

	// Store & Send rendered pixels
	renderer    coreio.FrameDrawer
	frameBuffer *coreio.FrameBuffer
	frameColors *coreio.FrameColors
	// frameBufferSprite coreio.FrameBuffer
	// frameBufferBg     coreio.FrameBuffer

	// Cycles counters
	frameCycles int
	lineCycles  int

	// STAT
	statInterruptOnCoincidence *ioports.BitPtr    // Bit 6 - LYC=LY Coincidence Interrupt (1=Enable) (Read/Write)
	statInterruptOnOAM         *ioports.BitPtr    // Bit 5 - Mode 2 OAM Interrupt         (1=Enable) (Read/Write)
	statInterruptOnVBlank      *ioports.BitPtr    // Bit 4 - Mode 1 V-Blank Interrupt     (1=Enable) (Read/Write)
	statInterruptOnHBlank      *ioports.BitPtr    // Bit 3 - Mode 0 H-Blank Interrupt     (1=Enable) (Read/Write)
	coincidenceFlag            *ioports.BitPtr    // Bit 2 - Coincidence Flag  (0:LYC<>LY, 1:LYC=LY) (Read Only)
	mode                       *ioports.MaskedPtr // Bit 1-0 - Mode Flag       (Mode 0-3, see below)
	// 								0: During H-Blank
	// 								1: During V-Blank
	// 								2: During Searching OAM-RAM
	// 								3: During Transfering Data to LCD Driver

	// Display registers
	ly  *ioports.Ptr // LY - LCDC Y-Coordinate
	lyc *ioports.Ptr // LYC - LY Compare

	// Background position
	scy *ioports.Ptr
	scx *ioports.Ptr

	// Colors palettes
	bgPalette      *ioports.Ptr // BGP - BG Palette Data
	spritePalette0 *ioports.Ptr // OBP0 - Object Palette 0 Data
	spritePalette1 *ioports.Ptr // OBP1 - Object Palette 1 Data

	// Window position
	wy *ioports.Ptr // WY - Window Y Position
	wx *ioports.Ptr // WX - Window X Position minus 7

	// LCDC Values
	displayEnable *ioports.BitPtr // Bit 7 - LCD Display Enable 			(0=Off, 1=On)
	winTileMap    *ioports.BitPtr // Bit 6 - Window Tile Map Display Select	(0=9800-9BFF, 1=9C00-9FFF)
	winEnable     *ioports.BitPtr // Bit 5 - Window Display Enable			(0=Off, 1=On)
	tileData      *ioports.BitPtr // Bit 4 - BG & Window Tile Data Select	(0=8800-97FF, 1=8000-8FFF)
	bgTileMap     *ioports.BitPtr // Bit 3 - BG Tile Map Display Select		(0=9800-9BFF, 1=9C00-9FFF)
	spriteSize    *ioports.BitPtr // Bit 2 - OBJ (Sprite) Size				(0=8x8, 1=8x16)
	spriteEnable  *ioports.BitPtr // Bit 1 - OBJ (Sprite) Display Enable	(0=Off, 1=On)
	bgEnable      *ioports.BitPtr // Bit 0 - BG Display 			 		(0=Off, 1=On)

	// Interrupts
	vBlankInterrupt *ioports.BitPtr
	statInterrupt   *ioports.BitPtr
}

func (gpu *GPU) getMode() gpuMode     { return gpuMode(gpu.mode.Get()) }
func (gpu *GPU) setMode(mode gpuMode) { gpu.mode.Set(uint8(mode)) }

func (gpu *GPU) setMonoColorPalette(offset int, palette uint8) {
	for i, c := range [...]uint8{
		0: palette & 0x03,
		1: palette >> 2 & 0x03,
		2: palette >> 4 & 0x03,
		3: palette >> 6 & 0x03,
	} {
		var intencity uint8 = 0
		switch c {
		case 0:
			intencity = 0xED
		case 1:
			intencity = 0x99
		case 2:
			intencity = 0x66
		case 3:
			intencity = 0x21
		default:
			log.Fatal("COLOR NOT DEFINED")

		}
		gpu.frameColors[(offset+i)*3] = intencity
		gpu.frameColors[(offset+i)*3+1] = intencity
		gpu.frameColors[(offset+i)*3+2] = intencity
	}
}

func (gpu *GPU) FlushFrameBuffer() {
	// Prepare colors palette
	if gpu.cgb {
		gpu.palettesManager.setPalettes(gpu.frameColors)
	} else {
		gpu.setMonoColorPalette(0b100000, gpu.bgPalette.Get())
		gpu.setMonoColorPalette(0b000000, gpu.spritePalette0.Get())
		gpu.setMonoColorPalette(0b000100, gpu.spritePalette1.Get())
	}

	// for i := 0; i < frameBufferSize; i++ {
	// 	if gpu.frameBufferSprite[i] != 255 {
	// 		gpu.frameBuffer[i] = gpu.frameBufferSprite[i]
	// 	} else {
	// 		gpu.frameBuffer[i] = gpu.frameBufferBg[i]
	// 	}
	// 	gpu.frameBufferSprite[i] = 255
	// 	gpu.frameBufferBg[i] = 0
	// }

	gpu.frameBuffer, gpu.frameColors =
		gpu.renderer.SwapFrameBuffer(gpu.frameBuffer, gpu.frameColors)
}

func (gpu *GPU) Tick(cycles uint8) uint8 {
	gpu.frameCycles += int(cycles)
	gpu.lineCycles += int(cycles)
	mode := gpu.getMode()

	if gpu.frameCycles > frameEnd { // IN V-BLANK
		if mode != ModeVBlank {
			// Last Line
			gpu.drawLine(gpu.ly.Get())
			// VBlank interrupt
			gpu.vBlankInterrupt.Set(true)
			if gpu.statInterruptOnVBlank.Get() {
				gpu.statInterrupt.Set(true)
			}
			gpu.FlushFrameBuffer()
			gpu.setMode(ModeVBlank)
		}
		if gpu.frameCycles >= vBlankEnd { // END : Reset cycles
			// gpu.frameCycles - vBlankEnd should not exceed 80 (OAM END)
			gpu.lineCycles = gpu.frameCycles - vBlankEnd
			gpu.frameCycles = gpu.frameCycles - vBlankEnd
			gpu.ly.Set(0)
			gpu.setMode(ModeOAM)
		}

	} else { // IN FRAME
		switch {
		case gpu.lineCycles < oamEnd:
			if mode != ModeOAM {
				gpu.setMode(ModeOAM)
				if gpu.statInterruptOnOAM.Get() {
					gpu.statInterrupt.Set(true)
				}
			}
			break
		case gpu.lineCycles < drawEnd:
			if mode != ModeVRAM {
				gpu.setMode(ModeVRAM)
				gpu.drawLine(gpu.ly.Get())
			}
			break
		case gpu.lineCycles < hBlanckEnd:
			if mode != ModeHBlank {
				gpu.setMode(ModeHBlank)
				if gpu.statInterruptOnHBlank.Get() {
					gpu.statInterrupt.Set(true)
				}
			}
			break
		}
	}

	// ALLWAYS UPDATE LY IN FRAME AND V-BLANK
	if gpu.lineCycles >= hBlanckEnd {
		gpu.ly.Set(gpu.ly.Get() + 1)
		gpu.lineCycles = gpu.lineCycles - hBlanckEnd
	}

	// Compare ly - lyc
	if gpu.ly.Get() == gpu.lyc.Get() {
		gpu.coincidenceFlag.Set(true)
		if gpu.statInterruptOnCoincidence.Get() {
			gpu.statInterrupt.Set(true)
		}
	} else {
		gpu.coincidenceFlag.Set(false)
	}
	return gpu.ly.Get()
}

func (gpu *GPU) drawLine(screenLine uint8) {
	if screenLine >= nbFrameLines { // if memory is corrupted...
		return
	}
	if gpu.bgEnable.Get() {
		gpu.drawBgLine(screenLine)
	}
	if gpu.winEnable.Get() {
		gpu.drawWindowLine(screenLine)
	}
	if gpu.spriteEnable.Get() {
		gpu.drawSpriteLine(screenLine)
	}
}
func (gpu *GPU) drawWindowLine(screenLine uint8) {
	bgScrollY := screenLine + gpu.wy.Get()
	bgScrollX := gpu.wx.Get() - 7
	if screenLine < gpu.wy.Get() {
		return
	}
	firstTileOffset := bgScrollX & 7
	var winMap uint8 = 0
	if gpu.winTileMap.Get() {
		winMap = 1
	}
	tileDataLine := bgScrollY & 7

	drawPointer := uint32(screenLine) * nbFrameRow

	tileDataTable := gpu.tileData.Get()

	nbPixelsDraw := 0
	var nbPixelsToDraw uint8 = 8
drawloop:
	for i := uint8(0); ; i++ { // Open range : break with label "drawloop"
		tileMapInfo := gpu._vram.GetTileInfo(winMap, (bgScrollX>>3+i)&0x1F, bgScrollY>>3)
		palettePrefix := 0b100000 | (tileMapInfo.palette)<<2
		gpu._vram.GetTileData(tileDataTable, tileMapInfo.tileID, tileMapInfo.bank).
			appendPixelsLine(gpu.frameBuffer[drawPointer:drawPointer], tileDataLine, firstTileOffset, nbPixelsToDraw, palettePrefix, tileMapInfo.hFlip == 1, tileMapInfo.vFlip == 1)

		drawPointer += (8 - uint32(firstTileOffset))
		nbPixelsDraw += (8 - int(firstTileOffset))
		if diff := nbFrameRow - nbPixelsDraw; diff <= 0 { // END of the line : 160 pixels draw
			break drawloop
		} else if diff < 8 {
			nbPixelsToDraw = uint8(diff)
		}
		firstTileOffset = 0
	}
}
func (gpu *GPU) drawSpriteLine(screenLine uint8) {

	// println("DRAW sprites")
	var maxHeight uint8 = 8
	mode8x16 := gpu.spriteSize.Get()
	if mode8x16 {
		maxHeight = 16
	}

	for i := 0; i < 40; i++ { // 40 sprites * 4 bytes
		sprite := gpu._oam.GetSprite(uint8(i))
		if sprite.Y == 0 || sprite.Y >= 160 { // Offscreen
			continue
		}
		topLeftX := int(sprite.X) - 8
		topLeftY := int(sprite.Y) - 16
		if topLeftY > int(screenLine) || topLeftY+int(maxHeight) <= int(screenLine) { // Not on the line
			continue
		}

		offscreenX := sprite.X == 0 || sprite.X >= 168 // Offscreen but affects the priority
		if !offscreenX {

			spriteLine := int(screenLine) - topLeftY

			for j, p := range sprite.GetPixels(&gpu._vram, uint8(spriteLine), 0, mode8x16) {
				pixelX := topLeftX + j
				if pixelX >= 160 {
					break
				}
				if pixelX < 0 {
					continue
				}
				// Sprite data 00 is transparent
				if p == 0 {
					continue
				}
				var colorPrefix uint8 = 0b000000
				if gpu.cgb {
					colorPrefix |= sprite.ColorPalette << 2
				} else {
					if sprite.PaletteNumber == 1 {
						colorPrefix = 0b000100
					}
				}

				// BG color 0 is always behind OBJ
				if gpu.frameBuffer[int(screenLine)*160+pixelX]&0x03 == 0 {
					gpu.frameBuffer[int(screenLine)*160+pixelX] = colorPrefix + p
				}
				// OBJ-to-BG Priority (0=OBJ Above BG, 1=OBJ Behind BG color 1-3)
				if !sprite.ObjToBgPriority {
					gpu.frameBuffer[int(screenLine)*160+pixelX] = colorPrefix + p
				}
			}
		}
	}

}

func (gpu *GPU) drawBgLine(screenLine uint8) {
	bgScrollY := screenLine + gpu.scy.Get()
	bgScrollX := gpu.scx.Get()
	firstTileOffset := bgScrollX & 7

	var bgMap uint8 = 0
	if gpu.bgTileMap.Get() {
		bgMap = 1
	}

	tileDataLine := bgScrollY & 7 // modulus 8

	drawPointer := uint32(screenLine) * nbFrameRow

	tileDataTable := gpu.tileData.Get()
	nbPixelsDraw := 0
	nbPixelsToDraw := uint8(8)
	// var buff [8]uint8
drawloop:
	for i := uint8(0); ; i++ { // Open range : break with label "drawloop"
		tileMapInfo := gpu._vram.GetTileInfo(bgMap, (bgScrollX>>3+i)&0x1F, bgScrollY>>3)
		palettePrefix := 0b100000 | (tileMapInfo.palette)<<2

		gpu._vram.GetTileData(tileDataTable, tileMapInfo.tileID, tileMapInfo.bank).
			appendPixelsLine(gpu.frameBuffer[drawPointer:drawPointer], tileDataLine, firstTileOffset, nbPixelsToDraw, palettePrefix, tileMapInfo.hFlip == 1, tileMapInfo.vFlip == 1)

		drawPointer += (8 - uint32(firstTileOffset))
		nbPixelsDraw += (8 - int(firstTileOffset))
		if diff := nbFrameRow - nbPixelsDraw; diff <= 0 { // END of the line : 160 pixels draw
			break drawloop
		} else if diff < 8 {
			nbPixelsToDraw = uint8(diff)
		}
		firstTileOffset = 0
	}
}

// Proxy for _vram access & OAM
func (gpu *GPU) Read(addr uint16) uint8 {
	switch {
	case
		addr == 0xFF68,
		addr == 0xFF69,
		addr == 0xFF6A,
		addr == 0xFF6B:
		return gpu.palettesManager.Read(addr)
	case addr == 0xFF46, addr >= memorymap.OAMStart && addr <= memorymap.OAMEnd:
		return gpu._oam.Read(addr)
	case addr >= memorymap.VRamStart && addr <= memorymap.VRamEnd:
		return gpu._vram.Read(addr)
	default:
		log.Fatalf("GPU MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}
func (gpu *GPU) Write(addr uint16, value uint8) {
	switch {
	case
		addr == 0xFF68,
		addr == 0xFF69,
		addr == 0xFF6A,
		addr == 0xFF6B:
		gpu.palettesManager.Write(addr, value)
	case addr == 0xFF46:
		gpu._oam.Write(addr, value)
	case addr >= memorymap.OAMStart && addr <= memorymap.OAMEnd:
		gpu._oam.internalWrite(addr, value)
	case addr >= memorymap.VRamStart && addr <= memorymap.VRamEnd:
		gpu._vram.Write(addr, value)
	default:
		log.Fatalf("GPU MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func NewGBGPU(io *ioports.IOPorts, renderer coreio.FrameDrawer, cgb bool) *GPU {
	gpu := &GPU{
		cgb:   cgb,
		_vram: newGBVRAM(io),
		_oam:  newOAM(),

		palettesManager: newPalettesManager(cgb),

		renderer: renderer,

		frameBuffer: new(coreio.FrameBuffer),
		frameColors: new(coreio.FrameColors),

		statInterruptOnCoincidence: io.NewBit6Ptr(0xFF41),
		statInterruptOnOAM:         io.NewBit5Ptr(0xFF41),
		statInterruptOnVBlank:      io.NewBit4Ptr(0xFF41),
		statInterruptOnHBlank:      io.NewBit3Ptr(0xFF41),
		coincidenceFlag:            io.NewBit2Ptr(0xFF41),
		mode:                       io.NewMaskedPtr(0xFF41, 0x03),

		ly:  io.NewPtr(0xFF44),
		lyc: io.NewPtr(0xFF45),

		wy: io.NewPtr(0xFF4A),
		wx: io.NewPtr(0xFF4B),

		scy: io.NewPtr(0xFF42),
		scx: io.NewPtr(0xFF43),

		statInterrupt:   io.NewBit1Ptr(0xFF0F),
		vBlankInterrupt: io.NewBit0Ptr(0xFF0F),

		displayEnable: io.NewBit7Ptr(0xFF40),
		winTileMap:    io.NewBit6Ptr(0xFF40),
		winEnable:     io.NewBit5Ptr(0xFF40),
		tileData:      io.NewBit4Ptr(0xFF40),
		bgTileMap:     io.NewBit3Ptr(0xFF40),
		spriteSize:    io.NewBit2Ptr(0xFF40),
		spriteEnable:  io.NewBit1Ptr(0xFF40),
		bgEnable:      io.NewBit0Ptr(0xFF40),

		bgPalette:      io.NewPtr(0xFF47), // BGP - BG Palette Data
		spritePalette0: io.NewPtr(0xFF48), // OBP0 - Object Palette 0 Data
		spritePalette1: io.NewPtr(0xFF49), // OBP1 - Object Palette 1 Data
	}
	return gpu
}
