package mmu

import (
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/hram"
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/mmu/memorymap"
	"github.com/jmontupet/gbcore/internal/pkg/unusableaddr"
	"github.com/jmontupet/gbcore/internal/pkg/wram"
)

func (m *MMU) Read(addr uint16) uint8 {
	switch {
	////// Cartridge bank 00 + Cartridge bank 01~NN //////
	case addr >= memorymap.FixedRomStart && addr <= memorymap.SwitchableRomEnd:
		return m.cartridge.Read(addr)

	////// VRAM  Switchable 0~1 on CGB//////
	case addr >= memorymap.VRamStart && addr <= memorymap.VRamEnd:
		return m.gpu.Read(addr)

	////// Cartridge RAM //////
	case addr >= memorymap.ExternalRamStart && addr <= memorymap.ExternalRamEnd:
		return m.cartridge.Read(addr)

	////// WRAM bank 0 + WRAM bank 1 ( ~ 7 on CGB ) //////
	case addr >= wram.WRamStart && addr <= wram.WRamEnd:
		return m.wram.Read(addr)

	////// Mirrored WRAM 0 & N //////
	case addr >= wram.MirrorWramStart && addr <= wram.MirrorWramEnd:
		return m.mirrorWram.Read(addr)

	////// OAM //////
	case addr >= memorymap.OAMStart && addr <= memorymap.OAMEnd:
		return m.gpu.Read(addr)

	////// NOT USABLE BUT USED //////
	case addr >= unusableaddr.AddrStart && addr <= unusableaddr.AddrEnd:
		return m.unusableAddr.Read(addr)

	////// IO Registers //////
	// Delegate control to Joypad
	case addr == 0xFF00:
		return m.joypad.Read(addr)
	// Delegate control to oamDMA for dma tranfer
	case addr == 0xFF46:
		return m.oamDMA.Read(addr)
	// Delegate control to vramDMA for dma tranfer
	case addr == 0xFF51, addr == 0xFF52,
		addr == 0xFF53, addr == 0xFF54, addr == 0xFF55:
		return m.vramDMA.Read(addr)

	// Delegate control to gpu for colors palettes
	case addr == 0xFF68,
		addr == 0xFF69,
		addr == 0xFF6A,
		addr == 0xFF6B:
		return m.gpu.Read(addr)
	case addr >= ioports.AddrStart && addr <= ioports.AddrEnd:
		return m.io.Read(addr)

	////// HRAM //////
	case addr >= hram.AddrStart && addr <= hram.AddrEnd:
		return m.hram.Read(addr)

	////// Interrupts Enable (IE) //////
	case addr == memorymap.Interrupts:
		return m.interrupt.Read(addr)
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}
