package mmu

import (
	"log"
)

type VramDmaManager struct {
	mmu *MMU

	srcAddr uint16
	dstAddr uint16

	transferLength uint16
	typeHBlank     bool

	transferActive bool
}

func (vdma *VramDmaManager) Read(addr uint16) uint8 {
	switch {
	// FF51 - HDMA1 - CGB Mode Only - New DMA Source, High
	case addr == 0xFF51:
		return uint8(vdma.srcAddr >> 8)
	// FF52 - HDMA2 - CGB Mode Only - New DMA Source, Low
	case addr == 0xFF52:
		return uint8(vdma.srcAddr & 0x00FF)
	// FF53 - HDMA3 - CGB Mode Only - New DMA Destination, High
	case addr == 0xFF53:
		return uint8(vdma.dstAddr >> 8)
	// FF54 - HDMA4 - CGB Mode Only - New DMA Destination, Low
	case addr == 0xFF54:
		return uint8(vdma.dstAddr & 0x00FF)
	// FF55 - HDMA5 - CGB Mode Only - New DMA Length/Mode/Start
	case addr == 0xFF55:
		if vdma.transferActive {
			return 0b10000000
		}
		return 0
	default:
		log.Fatalf("READ MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func (vdma *VramDmaManager) Write(addr uint16, value uint8) {
	switch {
	// FF51 - HDMA1 - CGB Mode Only - New DMA Source, High
	case addr == 0xFF51:
		vdma.srcAddr = uint16(value)<<8 | (vdma.srcAddr & 0x00FF)
	// FF52 - HDMA2 - CGB Mode Only - New DMA Source, Low
	case addr == 0xFF52:
		vdma.srcAddr = uint16(value)&0xF0 | (vdma.srcAddr & 0xFF00)
	// FF53 - HDMA3 - CGB Mode Only - New DMA Destination, High
	case addr == 0xFF53:
		vdma.dstAddr = uint16(value&0x1F|0x80)<<8 | (vdma.dstAddr & 0x00FF)
	// FF54 - HDMA4 - CGB Mode Only - New DMA Destination, Low
	case addr == 0xFF54:
		vdma.dstAddr = uint16(value)&0xF0 | (vdma.dstAddr & 0xFF00)
	// FF55 - HDMA5 - CGB Mode Only - New DMA Length/Mode/Start
	case addr == 0xFF55:
		vdma.typeHBlank = value>>7 == 1
		vdma.transferLength = ((uint16(value) & 0x7F) + 1) * 0x10
		vdma.transferActive = true
	default:
		log.Fatalf("WRITE MEMORY UNREACHABLE : (0x%04X) = 0x%0X", addr, value)
	}
}

func (vdma *VramDmaManager) Tick(cycles uint8) {
	if vdma.transferActive {
		for i := uint16(0); i < vdma.transferLength; i++ {
			vdma.mmu.Write(vdma.dstAddr+i, vdma.mmu.Read(vdma.srcAddr+i))
		}
		vdma.transferActive = false
	}
}
