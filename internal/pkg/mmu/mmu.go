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
	gpu          *gpu.GPU
	io           *ioports.IOPorts
	hram         *hram.HRAM
	wram         *wram.WRam
	mirrorWram   *wram.TranslatedWram
	interrupt    *interrupt.Manager
	joypad       *joypad.Joypad
	unusableAddr *unusableaddr.UnusableAddr
	oamDMA       *OamDmaManager
	vramDMA      *VramDmaManager
}

func (mmu *MMU) GetOamDMA() *OamDmaManager   { return mmu.oamDMA }
func (mmu *MMU) GetVramDMA() *VramDmaManager { return mmu.vramDMA }

func NewMMU(
	cart memory.Memory,
	gpu *gpu.GPU,
	io *ioports.IOPorts,
	hram *hram.HRAM,
	ram *wram.WRam,
	interrupt *interrupt.Manager,
	joypad *joypad.Joypad,
	unusableAddr *unusableaddr.UnusableAddr,
) *MMU {
	mmu := &MMU{
		cartridge:    cart,
		gpu:          gpu,
		io:           io,
		hram:         hram,
		wram:         ram,
		mirrorWram:   wram.NewTranlatedWram(ram),
		interrupt:    interrupt,
		joypad:       joypad,
		unusableAddr: unusableAddr,
	}
	mmu.oamDMA = &OamDmaManager{mmu: mmu}
	mmu.vramDMA = &VramDmaManager{mmu: mmu}
	return mmu
}
