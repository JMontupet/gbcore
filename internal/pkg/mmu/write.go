package mmu

import (
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/hram"
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/mmu/memorymap"
	"github.com/jmontupet/gbcore/internal/pkg/unusableaddr"
	"github.com/jmontupet/gbcore/internal/pkg/wram"
)

func (m *MMU) Write(addr uint16, value uint8) {
	switch {
	// Writing any value to this register resets it to 00h.
	case addr == 0xFF04:
		m.io.Write(0xFF04, 0)

	////// Cartridge bank 00 + Cartridge bank 01~NN //////
	case addr >= memorymap.FixedRomStart && addr <= memorymap.SwitchableRomEnd:
		m.cartridge.Write(addr, value)

	////// VRAM  Switchable 0~1 on CGB//////
	case addr >= memorymap.VRamStart && addr <= memorymap.VRamEnd:
		m.gpu.Write(addr, value)

	////// Cartridge RAM //////
	case addr >= memorymap.ExternalRamStart && addr <= memorymap.ExternalRamEnd:
		m.cartridge.Write(addr, value)

	////// WRAM bank 0 + WRAM bank 1 ( ~ 7 on CGB ) //////
	case addr >= wram.WRamStart && addr <= wram.WRamEnd:
		m.wram.Write(addr, value)

	////// Mirrored WRAM 0 & N //////
	case addr >= wram.MirrorWramStart && addr <= wram.MirrorWramEnd:
		m.mirrorWram.Write(addr, value)

	////// OAM //////
	case addr >= memorymap.OAMStart && addr <= memorymap.OAMEnd:
		m.gpu.Write(addr, value)

	////// NOT USABLE BUT USED //////
	case addr >= unusableaddr.AddrStart && addr <= unusableaddr.AddrEnd:
		m.unusableAddr.Write(addr, value)

	////// IO Registers //////
	// Delegate control to Joypad
	case addr == 0xFF00:
		m.joypad.Write(addr, value)
	// Delegate control to oamDMA for dma tranfer
	case addr == 0xFF46:
		m.oamDMA.Write(addr, value)
	// Delegate control to vramDMA for dma tranfer
	case addr == 0xFF51, addr == 0xFF52,
		addr == 0xFF53, addr == 0xFF54, addr == 0xFF55:
		m.vramDMA.Write(addr, value)
	// Delegate control to gpu for colors palettes
	case addr == 0xFF68,
		addr == 0xFF69,
		addr == 0xFF6A,
		addr == 0xFF6B:
		m.gpu.Write(addr, value)
	case addr >= ioports.AddrStart && addr <= ioports.AddrEnd:
		m.io.Write(addr, value)

	////// HRAM //////
	case addr >= hram.AddrStart && addr <= hram.AddrEnd:
		m.hram.Write(addr, value)

	////// Interrupts Enable (IE) //////
	case addr == memorymap.Interrupts:
		m.interrupt.Write(addr, value)
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}
