package mmu

import (
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/mmu/memorymap"
)

type OamDmaManager struct {
	_dmaRegister uint8

	mmu *MMU

	transferSrc    uint16
	transferActive bool
	transferDst    uint16
}

func (odma *OamDmaManager) Read(addr uint16) uint8 {
	switch {
	case odma.transferActive:
		return 0xFF
	case addr == 0xFF46:
		return odma._dmaRegister
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func (odma *OamDmaManager) Write(addr uint16, value uint8) {
	switch {
	case odma.transferActive:
	case addr == 0xFF46:
		odma.transferSrc = uint16(value) << 8
		odma.transferDst = 0xFE00
		odma.transferActive = true
		odma._dmaRegister = value
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func (odma *OamDmaManager) Tick(cycles uint8) {
	if !odma.transferActive {
		return
	}
	nbByte := memorymap.OAMEnd - odma.transferDst
	addrSrc := odma.transferSrc
	addrDst := odma.transferDst
	if nbByte > uint16(cycles) {
		nbByte = uint16(cycles)
	}
	odma.transferDst += nbByte
	odma.transferSrc += nbByte
	if odma.transferDst >= memorymap.OAMEnd {
		odma.transferActive = false
	}
	for i := uint16(0); i < nbByte; i++ {
		odma.mmu.Write(addrDst+i, odma.mmu.Read(addrSrc+i))
	}
}
