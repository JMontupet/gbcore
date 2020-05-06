package cpu

import (
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/memory"

	"github.com/jmontupet/gbcore/internal/pkg/cpu/registers"

	"github.com/jmontupet/gbcore/internal/pkg/interrupt"

	"github.com/jmontupet/gbcore/internal/pkg/mmu"
)

// CPU emulate GameBoy CPU
type CPU struct {
	regs        registers.Registers
	mmu         memory.Memory
	interrupts  *interrupt.Manager
	halt        bool
	DoubleSpeed bool
}

// readUint16 read next uint16 value from the mmu at ProgramCounter address and inc2 PC
func (c *CPU) readUint16() uint16 {
	pc := c.regs.GetPC()
	c.regs.SetPC(pc + 2)
	return uint16(c.mmu.Read(pc)) | (uint16(c.mmu.Read(pc+1)) << 8)
}

// readUint8 read next uint8 value from the mmu at ProgramCounter address and inc PC
func (c *CPU) readUint8() uint8 {
	pc := c.regs.GetPC()
	c.regs.SetPC(pc + 1)
	return c.mmu.Read(pc)
}

// stop CPU restart the cpu at new speed if required
func (c *CPU) stop() {
	println("STOP CPU")

	reg := c.mmu.Read(0xFF4D) // CPU speed / CGB mode
	prepareSwitchMode := reg & 0x01
	currentMode := reg >> 7

	if prepareSwitchMode == 1 {
		if currentMode == 1 {
			println("SWITCH TO NORMAL SPEED MODE")
			// Switch to Normal Mode
			c.DoubleSpeed = false
			c.mmu.Write(0xFF4D, 0x00)
		} else {
			println("SWITCH TO DOUBLE SPEED MODE")
			// Switch to Double Speed Mode
			c.DoubleSpeed = true
			c.mmu.Write(0xFF4D, 0x80)
		}
	}
}

// Tick read the next opcode at address PC and execute corresponding instruction
func (c *CPU) Tick() (clockUsed uint8) {
	if addrInterrupt := c.interrupts.GetNext(); addrInterrupt != 0x0000 {
		c.halt = false
		call(c, addrInterrupt)
		return 8
	}
	if c.halt {
		return 1
	}
	code := c.readUint8()
	if code == 0xCB {
		code = c.readUint8()
		if instructionCBList[code] != nil {
			return instructionCBList[code](c)
		}
		log.Fatalf("INVALID CB OPCODE : %X \n", code)
		return 0
	}
	if instructionList[code] != nil {
		return instructionList[code](c)
	}
	log.Fatalf("INVALID OPCODE : %X \n", code)
	return 0
}

// NewCPU return a new initialised GameBoy CPU
func NewCPU(mem *mmu.MMU, interrupts *interrupt.Manager) *CPU {
	var regs = registers.Registers{}

	// SHORTCUT TO INIT CPU & MEMORY WITHOUT BOOT SEQUENCE
	regs.SetPC(0x0100)
	regs.SetAF(0x11B0) // A = 0x11 -> CGB / A = 0x01 -> GB
	regs.SetBC(0x0012)
	regs.SetDE(0x00D8)
	regs.SetHL(0x014D)
	regs.SetSP(0xFFFE)

	mem.Write(0xFF50, 0x01)
	mem.Write(0xFF05, 0x00) // TIMA
	mem.Write(0xFF06, 0x00) // TMA
	mem.Write(0xFF07, 0x00) // TAC
	mem.Write(0xFF10, 0x80) // NR10
	mem.Write(0xFF11, 0xBF) // NR11
	mem.Write(0xFF12, 0xF3) // NR12
	mem.Write(0xFF14, 0xBF) // NR14
	mem.Write(0xFF16, 0x3F) // NR21
	mem.Write(0xFF17, 0x00) // NR22
	mem.Write(0xFF19, 0xBF) // NR24
	mem.Write(0xFF1A, 0x7F) // NR30
	mem.Write(0xFF1B, 0xFF) // NR31
	mem.Write(0xFF1C, 0x9F) // NR32
	mem.Write(0xFF1E, 0xBF) // NR33
	mem.Write(0xFF20, 0xFF) // NR41
	mem.Write(0xFF21, 0x00) // NR42
	mem.Write(0xFF22, 0x00) // NR43
	mem.Write(0xFF23, 0xBF) // NR30
	mem.Write(0xFF24, 0x77) // NR50
	mem.Write(0xFF25, 0xF3) // NR51
	mem.Write(0xFF26, 0xF1) // NR52
	mem.Write(0xFF40, 0x91) // LCDC
	mem.Write(0xFF42, 0x00) // SCY
	mem.Write(0xFF43, 0x00) // SCX
	mem.Write(0xFF45, 0x00) // LYC
	mem.Write(0xFF47, 0xFC) // BGP
	mem.Write(0xFF48, 0xFF) // OBP0
	mem.Write(0xFF49, 0xFF) // OBP1
	mem.Write(0xFF4A, 0x00) // WY
	mem.Write(0xFF4B, 0x00) // WX
	mem.Write(0xFFFF, 0x00) // IE

	return &CPU{
		mmu:        mem,
		interrupts: interrupts,
		regs:       regs,
	}
}
