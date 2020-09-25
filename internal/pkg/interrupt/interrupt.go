package interrupt

import (
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/internal/pkg/memory"
)

type Manager interface {
	memory.Memory
	GetNext() uint16
	EnableMaster()
	DisableMaster()
}

func (interrupt *manager) EnableMaster() {
	// fmt.Println("INTERRUPTS ENABLED")
	interrupt.masterFlag = true
}
func (interrupt *manager) DisableMaster() {
	// fmt.Println("INTERRUPTS DISABLED")
	interrupt.masterFlag = false
}

type manager struct {
	// ME - Interrupt Master Enable Flag
	masterFlag bool

	// 	iEnable -> (0xFFFF)  IE - Interrupt Enable
	// 		Bit 0: V-Blank  Interrupt Enable  (INT 40h)  (1=Enable)
	// 		Bit 1: LCD STAT Interrupt Enable  (INT 48h)  (1=Enable)
	// 		Bit 2: Timer    Interrupt Enable  (INT 50h)  (1=Enable)
	// 		Bit 3: Serial   Interrupt Enable  (INT 58h)  (1=Enable)
	// 		Bit 4: Joypad   Interrupt Enable  (INT 60h)  (1=Enable)
	iEnable uint8

	// 	iFlag -> (0xFF0F) IF - Interrupt Flag
	// 		Bit 0: V-Blank  Interrupt Request (INT 40h)  (1=Request)
	// 		Bit 1: LCD STAT Interrupt Request (INT 48h)  (1=Request)
	// 		Bit 2: Timer    Interrupt Request (INT 50h)  (1=Request)
	// 		Bit 3: Serial   Interrupt Request (INT 58h)  (1=Request)
	// 		Bit 4: Joypad   Interrupt Request (INT 60h)  (1=Request)
	iFlag *ioports.Ptr
}

func (interrupt *manager) GetNext() uint16 {
	if interrupt.masterFlag {
		validInterrupts := interrupt.iFlag.Get() & interrupt.iEnable
		switch {
		case validInterrupts == 0:
			return 0x0000
		case validInterrupts&0x01 > 0: // V-Blank
			interrupt.iFlag.SetBit0(false)
			interrupt.masterFlag = false
			return 0x0040
		case validInterrupts&0x02 > 0: // LCD STAT
			interrupt.iFlag.SetBit1(false)
			interrupt.masterFlag = false
			return 0x0048
		case validInterrupts&0x04 > 0: // Timer
			interrupt.iFlag.SetBit2(false)
			interrupt.masterFlag = false
			return 0x0050
		case validInterrupts&0x08 > 0: // Serial
			interrupt.iFlag.SetBit3(false)
			interrupt.masterFlag = false
			return 0x0058
		case validInterrupts&0x10 > 0: // Joypad
			interrupt.iFlag.SetBit4(false)
			interrupt.masterFlag = false
			return 0x0060
		default:
			return 0x000
		}
	}
	return 0x0000 // 0x0000 is not a valid interrupt jump address, used as 'NO INTERRUPT' for this function ONLY !!!
}

func (interrupt *manager) Read(_ uint16) uint8 {
	return interrupt.iEnable
}

func (interrupt *manager) Write(_ uint16, value uint8) {
	interrupt.iEnable = value
}

func NewInterrupt(io *ioports.IOPorts) Manager {
	interrupt := &manager{
		iFlag: io.NewPtr(0xFF0F),
	}
	return interrupt
}
