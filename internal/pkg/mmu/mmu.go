package mmu

import (
	"github.com/jmontupet/gbcore/internal/pkg/hram"

	"github.com/jmontupet/gbcore/internal/pkg/unusableaddr"

	"github.com/jmontupet/gbcore/internal/pkg/gpu"

	"github.com/jmontupet/gbcore/internal/pkg/joypad"
	"github.com/jmontupet/gbcore/internal/pkg/wram"

	"github.com/jmontupet/gbcore/internal/pkg/interrupt"
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/memory"
)

type MMU struct {
	cartridge    memory.Memory
	gpu          memory.Memory
	io           memory.Memory
	hram         memory.Memory
	wram         memory.Memory
	mirrorWram   memory.Memory
	interrupt    memory.Memory
	joypad       memory.Memory
	unusableAddr memory.Memory
	oamDMA       *OamDmaManager
	vramDMA      *VramDmaManager
}

func (mmu *MMU) GetOamDMA() *OamDmaManager   { return mmu.oamDMA }
func (mmu *MMU) GetVramDMA() *VramDmaManager { return mmu.vramDMA }

func (mmu *MMU) SetCartridge(c memory.Memory) { mmu.cartridge = c }

func NewMMU(
	gbCart memory.Memory,
	gbGPU *gpu.GPU,
	gbIO *ioports.IOPorts,
	gbHRAM *hram.HRAM,
	gbRAM *wram.WRam,
	gbInterrupt interrupt.Manager,
	gbJoypad *joypad.Joypad,
	gbUnusableAddr *unusableaddr.UnusableAddr,
) *MMU {
	mmu := &MMU{
		cartridge:    gbCart,
		gpu:          gbGPU,
		io:           gbIO,
		hram:         gbHRAM,
		wram:         gbRAM,
		mirrorWram:   wram.NewTranlatedWram(gbRAM),
		interrupt:    gbInterrupt,
		joypad:       gbJoypad,
		unusableAddr: gbUnusableAddr,
	}
	mmu.oamDMA = &OamDmaManager{mmu: mmu}
	mmu.vramDMA = &VramDmaManager{mmu: mmu}
	return mmu
}
