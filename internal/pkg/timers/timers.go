package timers

import (
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
)

const cpuClock = 4194304
const divClock uint = cpuClock / 16384
const tacClock0 uint = cpuClock / 4096
const tacClock1 uint = cpuClock / 262144
const tacClock2 uint = cpuClock / 65536
const tacClock3 uint = cpuClock / 16384

type Timers struct {
	timerInt *ioports.BitPtr // Pointer interrupt when tima overflow

	div  *ioports.Ptr // Divider Register 16384Hz
	tima *ioports.Ptr // Timer counter
	tma  *ioports.Ptr // Timer Modulo
	tac  *ioports.Ptr // Timer Control   	Bit 2    - Timer Stop  (0=Stop, 1=Start)
	// 									Bits 1-0 - Input Clock Select
	// 									00:   4096 Hz    (~4194 Hz SGB)
	// 									01: 262144 Hz  (~268400 Hz SGB)
	// 									10:  65536 Hz   (~67110 Hz SGB)
	// 									11:  16384 Hz   (~16780 Hz SGB)

	// Cycles counters
	divCount  uint
	timaCount uint
}

func (t *Timers) Tick(cycles uint8) {
	// INC DIV register
	t.divCount += uint(cycles)
	if t.divCount >= divClock {
		div := t.div.Get()
		if div == 0xFF { // Overflow
			div = 0
		} else {
			div++
		}
		t.div.Set(div)
		t.divCount %= divClock
	}

	if t.tac.GetBit2() { // Check if timer is started
		// INC TIMA register
		t.timaCount += uint(cycles)
		clock := t.cyclesTIMAInc()
		if t.timaCount >= clock {
			tima := t.tima.Get()
			if tima == 0xFF { // Overflow
				tima = t.tma.Get()
				t.timerInt.Set(true)
			} else {
				tima++
			}
			t.tima.Set(tima)
			t.timaCount %= clock
		}
	}
}

func (t *Timers) cyclesTIMAInc() uint {
	switch t.tac.Get() & 0b00000011 {
	case 0:
		return tacClock0
	case 1:
		return tacClock1
	case 2:
		return tacClock2
	case 3:
		return tacClock3
	default:
		return 1 // Impossible. All 2bits combinaisons filtered
	}
}

func NewTimers(io *ioports.IOPorts) *Timers {
	return &Timers{
		div:      io.NewPtr(0xFF04),     // DIV
		tima:     io.NewPtr(0xFF05),     // TIMA
		tma:      io.NewPtr(0xFF06),     // TMA
		tac:      io.NewPtr(0xFF07),     // TAC
		timerInt: io.NewBit2Ptr(0xFF0F), // Interrupt
	}
}
