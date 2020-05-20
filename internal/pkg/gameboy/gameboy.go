package gameboy

import (
	"math"
	"time"

	"github.com/jmontupet/gbcore/internal/pkg/constants"

	"github.com/jmontupet/gbcore/pkg/coreio"

	"github.com/jmontupet/gbcore/internal/pkg/audio"
	"github.com/jmontupet/gbcore/internal/pkg/joypad"
	"github.com/jmontupet/gbcore/internal/pkg/unusableaddr"

	"github.com/jmontupet/gbcore/internal/pkg/interrupt"
	"github.com/jmontupet/gbcore/internal/pkg/timers"

	"github.com/jmontupet/gbcore/internal/pkg/wram"

	"github.com/jmontupet/gbcore/internal/pkg/cartridge"
	"github.com/jmontupet/gbcore/internal/pkg/cpu"
	"github.com/jmontupet/gbcore/internal/pkg/gpu"
	"github.com/jmontupet/gbcore/internal/pkg/hram"
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/mmu"
)

type GameBoy interface {
	Run()
}

type gameboy struct {
	cpu    *cpu.CPU
	gpu    *gpu.GPU
	apu    *audio.APU
	mmu    *mmu.MMU
	joypad *joypad.Joypad
	timers *timers.Timers

	inputsManager coreio.InputsManager
}

func (gb *gameboy) Run() {
	const nbRefreshPerFrame = constants.InputRefreshPerFrame
	const frameDiv = 154 / nbRefreshPerFrame

	ticker := time.NewTicker(time.Duration(math.Round(1000000000 / constants.ScreenRefreshRate / nbRefreshPerFrame)))
	defer ticker.Stop()

	prevLine := uint8(0)
	for {
		nbClockUsed := gb.cpu.Tick()
		var clockMul uint8 = 4
		if gb.cpu.DoubleSpeed {
			clockMul = 2
		}

		line := gb.gpu.Tick(nbClockUsed * 4)
		gb.timers.Tick(nbClockUsed * clockMul)
		gb.mmu.GetOamDMA().Tick(nbClockUsed * clockMul)
		gb.mmu.GetVramDMA().Tick(nbClockUsed * 4)

		// gb.apu.Tick(nbClockUsed)

		if line%frameDiv == 0 && prevLine%frameDiv != 0 { // 0 - 153
			gb.joypad.UpdateInput(uint8(gb.inputsManager.CurrentInput()))
			// <-ticker.C
		}
		prevLine = line
	}
}

func NewGameBoy(
	cart cartridge.Cartridge,
	renderer coreio.FrameDrawer,
	inputsManager coreio.InputsManager,
	audioPlayer coreio.AudioPlayer,

) GameBoy {
	cgb := cartridge.ReadCGBCompatible(cart)

	gbIO := ioports.NewGBIOPorts()
	gbTImers := timers.NewTimers(gbIO)
	gbHRAM := hram.NewGBHRAM()
	gbWRAM := wram.NewWram(gbIO)
	gbInterrupt := interrupt.NewInterrupt(gbIO)
	gbJoypad := joypad.NewJoypad(gbIO)

	unusableAddr := unusableaddr.NewUnusableAddr()
	gbGPU := gpu.NewGBGPU(gbIO, renderer, cgb)
	gbMMU := mmu.NewMMU(cart, gbGPU, gbIO, gbHRAM, gbWRAM, gbInterrupt, gbJoypad, unusableAddr)
	proc := cpu.NewCPU(gbMMU, gbInterrupt)
	apu := audio.NewAPU(gbIO, audioPlayer)

	return &gameboy{
		cpu:           proc,
		gpu:           gbGPU,
		apu:           apu,
		mmu:           gbMMU,
		timers:        gbTImers,
		joypad:        gbJoypad,
		inputsManager: inputsManager,
	}
}
