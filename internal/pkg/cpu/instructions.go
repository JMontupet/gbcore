package cpu

import (
	"math/bits"

	"github.com/jmontupet/gbcore/internal/pkg/cpu/registers"
)

type instruction func(*CPU) (clockUsed uint8)

type instructions [256]instruction

// type instructions map[uint8]instruction

var instructionList = instructions{
	0x00: func(cpu *CPU) uint8 { return 1 },             // NOP
	0x10: func(cpu *CPU) uint8 { cpu.stop(); return 1 }, // STOP

	///// 8-Bit Loads //////
	// LD r,n
	0x3E: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.readUint8()); return 2 }, // LD A,d8
	0x06: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.readUint8()); return 2 }, // LD B,d8
	0x0E: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.readUint8()); return 2 }, // LD C,d8
	0x16: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.readUint8()); return 2 }, // LD D,d8
	0x1E: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.readUint8()); return 2 }, // LD E,d8
	0x26: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.readUint8()); return 2 }, // LD H,d8
	0x2E: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.readUint8()); return 2 }, // LD L,d8

	// LD r1,r2
	0x7F: func(cpu *CPU) uint8 { return 1 },                                 // LD A,A
	0x78: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.regs.GetB()); return 1 }, // LD A,B
	0x79: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.regs.GetC()); return 1 }, // LD A,C
	0x7A: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.regs.GetD()); return 1 }, // LD A,D
	0x7B: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.regs.GetE()); return 1 }, // LD A,E
	0x7C: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.regs.GetH()); return 1 }, // LD A,H
	0x7D: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.regs.GetL()); return 1 }, // LD A,L

	0x47: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.regs.GetA()); return 1 }, // LD B,A
	0x40: func(cpu *CPU) uint8 { return 1 },                                 // LD B,B
	0x41: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.regs.GetC()); return 1 }, // LD B,C
	0x42: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.regs.GetD()); return 1 }, // LD B,D
	0x43: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.regs.GetE()); return 1 }, // LD B,E
	0x44: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.regs.GetH()); return 1 }, // LD B,H
	0x45: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.regs.GetL()); return 1 }, // LD B,L

	0x4F: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.regs.GetA()); return 1 }, // LD C,A
	0x48: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.regs.GetB()); return 1 }, // LD C,B
	0x49: func(cpu *CPU) uint8 { return 1 },                                 // LD C,C
	0x4A: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.regs.GetD()); return 1 }, // LD C,D
	0x4B: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.regs.GetE()); return 1 }, // LD C,E
	0x4C: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.regs.GetH()); return 1 }, // LD C,H
	0x4D: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.regs.GetL()); return 1 }, // LD C,L

	0x57: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.regs.GetA()); return 1 }, // LD D,A
	0x50: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.regs.GetB()); return 1 }, // LD D,B
	0x51: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.regs.GetC()); return 1 }, // LD D,C
	0x52: func(cpu *CPU) uint8 { return 1 },                                 // LD D,D
	0x53: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.regs.GetE()); return 1 }, // LD D,E
	0x54: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.regs.GetH()); return 1 }, // LD D,H
	0x55: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.regs.GetL()); return 1 }, // LD D,L

	0x5F: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.regs.GetA()); return 1 }, // LD E,A
	0x58: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.regs.GetB()); return 1 }, // LD E,B
	0x59: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.regs.GetC()); return 1 }, // LD E,C
	0x5A: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.regs.GetD()); return 1 }, // LD E,D
	0x5B: func(cpu *CPU) uint8 { return 1 },                                 // LD E,E
	0x5C: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.regs.GetH()); return 1 }, // LD E,H
	0x5D: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.regs.GetL()); return 1 }, // LD E,L

	0x67: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.regs.GetA()); return 1 }, // LD H,A
	0x60: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.regs.GetB()); return 1 }, // LD H,B
	0x61: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.regs.GetC()); return 1 }, // LD H,C
	0x62: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.regs.GetD()); return 1 }, // LD H,D
	0x63: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.regs.GetE()); return 1 }, // LD H,E
	0x64: func(cpu *CPU) uint8 { return 1 },                                 // LD H,H
	0x65: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.regs.GetL()); return 1 }, // LD H,L

	0x6F: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.regs.GetA()); return 1 }, // LD L,A
	0x68: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.regs.GetB()); return 1 }, // LD L,B
	0x69: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.regs.GetC()); return 1 }, // LD L,C
	0x6A: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.regs.GetD()); return 1 }, // LD L,D
	0x6B: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.regs.GetE()); return 1 }, // LD L,E
	0x6C: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.regs.GetH()); return 1 }, // LD L,H
	0x6D: func(cpu *CPU) uint8 { return 1 },                                 // LD L,L

	// LD r, (HL)
	0x7E: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD A, (HL)
	0x46: func(cpu *CPU) uint8 { cpu.regs.SetB(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD B, (HL)
	0x4E: func(cpu *CPU) uint8 { cpu.regs.SetC(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD C, (HL)
	0x56: func(cpu *CPU) uint8 { cpu.regs.SetD(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD D, (HL)
	0x5E: func(cpu *CPU) uint8 { cpu.regs.SetE(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD E, (HL)
	0x66: func(cpu *CPU) uint8 { cpu.regs.SetH(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD H, (HL)
	0x6E: func(cpu *CPU) uint8 { cpu.regs.SetL(cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // LD L, (HL)

	// LD (HL), r
	0x70: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetB()); return 2 }, // LD (HL), B
	0x71: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetC()); return 2 }, // LD (HL), C
	0x72: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetD()); return 2 }, // LD (HL), D
	0x73: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetE()); return 2 }, // LD (HL), E
	0x74: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetH()); return 2 }, // LD (HL), H
	0x75: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetL()); return 2 }, // LD (HL), L

	// LD (HL), n
	0x36: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.readUint8()); return 3 }, // LD (HL), n

	// LD A, (nn)
	0xFA: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.mmu.Read(cpu.readUint16())); return 4 }, // LD A, (nn)

	// LD (nn), A
	0xEA: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.readUint16(), cpu.regs.GetA()); return 4 }, // LD (nn), A

	// LD A, (r)
	0x0A: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.mmu.Read(cpu.regs.GetBC())); return 2 }, // LD A, (BC)
	0x1A: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.mmu.Read(cpu.regs.GetDE())); return 2 }, // LD A, (DE)

	// LD A, (C)
	0xF2: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.mmu.Read(uint16(cpu.regs.GetC()) | 0xFF00)); return 2 }, // LD A, (C)

	// LD (C), A
	0xE2: func(cpu *CPU) uint8 { cpu.mmu.Write(uint16(cpu.regs.GetC())|0xFF00, cpu.regs.GetA()); return 2 }, // LD (C),A

	// LD A, (n)
	0xF0: func(cpu *CPU) uint8 { cpu.regs.SetA(cpu.mmu.Read(uint16(cpu.readUint8()) | 0xFF00)); return 3 }, // LD A, (n)

	// LD (n), A
	0xE0: func(cpu *CPU) uint8 { cpu.mmu.Write(uint16(cpu.readUint8())|0xFF00, cpu.regs.GetA()); return 3 }, // LD (n),A

	// LD (r), A
	0x02: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetBC(), cpu.regs.GetA()); return 2 }, // LD (BC), A
	0x12: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetDE(), cpu.regs.GetA()); return 2 }, // LD (DE), A
	0x77: func(cpu *CPU) uint8 { cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetA()); return 2 }, // LD (HL), A

	// LD A, (HL-)
	0x3A: func(cpu *CPU) uint8 {
		cpu.regs.SetA(cpu.mmu.Read(cpu.regs.GetHL()))
		cpu.regs.SetHL(cpu.regs.GetHL() - 1)
		return 2
	}, // LD A, (HL-)
	// LD A, (HL+)
	0x2A: func(cpu *CPU) uint8 {
		cpu.regs.SetA(cpu.mmu.Read(cpu.regs.GetHL()))
		cpu.regs.SetHL(cpu.regs.GetHL() + 1)
		return 2
	}, // LD A, (HL+)

	// LD (HL-), A
	0x32: func(cpu *CPU) uint8 {
		cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetA())
		cpu.regs.SetHL(cpu.regs.GetHL() - 1)
		return 2
	}, // LD (HL-), A
	// LD (HL+), A
	0x22: func(cpu *CPU) uint8 {
		cpu.mmu.Write(cpu.regs.GetHL(), cpu.regs.GetA())
		cpu.regs.SetHL(cpu.regs.GetHL() + 1)
		return 2
	}, // LD (HL+), A

	///// 16-Bit Loads //////
	// LD r,nn
	0x01: func(cpu *CPU) uint8 { cpu.regs.SetBC(cpu.readUint16()); return 3 }, // LD BC, nn
	0x11: func(cpu *CPU) uint8 { cpu.regs.SetDE(cpu.readUint16()); return 3 }, // LD DE, nn
	0x21: func(cpu *CPU) uint8 { cpu.regs.SetHL(cpu.readUint16()); return 3 }, // LD HL, nn
	0x31: func(cpu *CPU) uint8 { cpu.regs.SetSP(cpu.readUint16()); return 3 }, // LD SP, nn

	// LD SP,HL
	0xF9: func(cpu *CPU) uint8 { cpu.regs.SetSP(cpu.regs.GetHL()); return 2 }, // LD SP,HL

	// LD HL, SP+n
	0xF8: func(cpu *CPU) uint8 {
		n := int32(int8(cpu.readUint8()))
		sp := int32(cpu.regs.GetSP())
		result := sp + n
		cpu.regs.SetHL(uint16(result))
		cpu.regs.SetFlag(registers.FlagZERO, false)
		cpu.regs.SetFlag(registers.FlagSUB, false)
		cpu.regs.SetFlag(registers.FlagHCARRY, (result&0xF) < (sp&0xF))
		cpu.regs.SetFlag(registers.FlagCARRY, (result&0xFF) < (sp&0xFF))

		return 3
	}, // LD HL, SP+n

	// LD (nn),SP
	0x08: func(cpu *CPU) uint8 {
		sp := cpu.regs.GetSP()
		addr := cpu.readUint16()
		cpu.mmu.Write(addr, uint8(sp&0x00FF))
		cpu.mmu.Write(addr+1, uint8(sp>>8))
		return 5
	}, // LD (nn),SP

	// PUSH nn
	0xF5: func(cpu *CPU) uint8 { push(cpu, cpu.regs.GetAF()); return 4 }, // PUSH AF
	0xC5: func(cpu *CPU) uint8 { push(cpu, cpu.regs.GetBC()); return 4 }, // PUSH BC
	0xD5: func(cpu *CPU) uint8 { push(cpu, cpu.regs.GetDE()); return 4 }, // PUSH DE
	0xE5: func(cpu *CPU) uint8 { push(cpu, cpu.regs.GetHL()); return 4 }, // PUSH HL

	// POP nn
	0xF1: func(cpu *CPU) uint8 { cpu.regs.SetAF(pop(cpu) & 0xFFF0); return 3 }, // POP AF 4bits of unused flags set to 0
	0xC1: func(cpu *CPU) uint8 { cpu.regs.SetBC(pop(cpu)); return 3 },          // POP BC
	0xD1: func(cpu *CPU) uint8 { cpu.regs.SetDE(pop(cpu)); return 3 },          // POP DE
	0xE1: func(cpu *CPU) uint8 { cpu.regs.SetHL(pop(cpu)); return 3 },          // POP HL

	///// 8-Bit ALU //////
	// ADD A,r
	0x87: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetA(), false); return 1 }, // ADD A, A
	0x80: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetB(), false); return 1 }, // ADD A, B
	0x81: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetC(), false); return 1 }, // ADD A, C
	0x82: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetD(), false); return 1 }, // ADD A, D
	0x83: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetE(), false); return 1 }, // ADD A, E
	0x84: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetH(), false); return 1 }, // ADD A, H
	0x85: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetL(), false); return 1 }, // ADD A, L
	// ADD A,(HL)
	0x86: func(cpu *CPU) uint8 { add8(cpu, cpu.mmu.Read(cpu.regs.GetHL()), false); return 2 }, // ADD A, (HL)
	// ADD A,n
	0xC6: func(cpu *CPU) uint8 { add8(cpu, cpu.readUint8(), false); return 2 }, // ADD A,n

	// ADC A,r
	0x8F: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetA(), true); return 1 }, // ADC A, A
	0x88: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetB(), true); return 1 }, // ADC A, B
	0x89: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetC(), true); return 1 }, // ADC A, C
	0x8A: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetD(), true); return 1 }, // ADC A, D
	0x8B: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetE(), true); return 1 }, // ADC A, E
	0x8C: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetH(), true); return 1 }, // ADC A, H
	0x8D: func(cpu *CPU) uint8 { add8(cpu, cpu.regs.GetL(), true); return 1 }, // ADC A, L
	// ADC A,(HL)
	0x8E: func(cpu *CPU) uint8 { add8(cpu, cpu.mmu.Read(cpu.regs.GetHL()), true); return 2 }, // ADC A, (HL)
	// ADC A,n
	0xCE: func(cpu *CPU) uint8 { add8(cpu, cpu.readUint8(), true); return 2 }, // ADC A,n

	// SUB A,r
	0x97: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetA(), false); return 1 }, // SUB A, A
	0x90: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetB(), false); return 1 }, // SUB A, B
	0x91: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetC(), false); return 1 }, // SUB A, C
	0x92: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetD(), false); return 1 }, // SUB A, D
	0x93: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetE(), false); return 1 }, // SUB A, E
	0x94: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetH(), false); return 1 }, // SUB A, H
	0x95: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetL(), false); return 1 }, // SUB A, L
	// SUB A,(HL)
	0x96: func(cpu *CPU) uint8 { sub8(cpu, cpu.mmu.Read(cpu.regs.GetHL()), false); return 2 }, // SUB A, (HL)
	// SUB A,n
	0xD6: func(cpu *CPU) uint8 { sub8(cpu, cpu.readUint8(), false); return 2 }, // SUB A,n

	// SBC A,r
	0x9F: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetA(), true); return 1 }, // SBC A, A
	0x98: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetB(), true); return 1 }, // SBC A, B
	0x99: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetC(), true); return 1 }, // SBC A, C
	0x9A: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetD(), true); return 1 }, // SBC A, D
	0x9B: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetE(), true); return 1 }, // SBC A, E
	0x9C: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetH(), true); return 1 }, // SBC A, H
	0x9D: func(cpu *CPU) uint8 { sub8(cpu, cpu.regs.GetL(), true); return 1 }, // SBC A, L
	// SBC A,(HL)
	0x9E: func(cpu *CPU) uint8 { sub8(cpu, cpu.mmu.Read(cpu.regs.GetHL()), true); return 2 }, // SBC A, (HL)
	// SBC A,n
	0xDE: func(cpu *CPU) uint8 { sub8(cpu, cpu.readUint8(), true); return 2 }, // SBC A,n

	// AND A,r
	0xA7: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetA()); return 1 }, // AND A, A
	0xA0: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetB()); return 1 }, // AND A, B
	0xA1: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetC()); return 1 }, // AND A, C
	0xA2: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetD()); return 1 }, // AND A, D
	0xA3: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetE()); return 1 }, // AND A, E
	0xA4: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetH()); return 1 }, // AND A, H
	0xA5: func(cpu *CPU) uint8 { and8(cpu, cpu.regs.GetL()); return 1 }, // AND A, L
	// AND A,(HL)
	0xA6: func(cpu *CPU) uint8 { and8(cpu, cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // AND A, (HL)
	// AND A,n
	0xE6: func(cpu *CPU) uint8 { and8(cpu, cpu.readUint8()); return 2 }, // AND A,n

	// OR A,r
	0xB7: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetA()); return 1 }, // OR A, A
	0xB0: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetB()); return 1 }, // OR A, B
	0xB1: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetC()); return 1 }, // OR A, C
	0xB2: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetD()); return 1 }, // OR A, D
	0xB3: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetE()); return 1 }, // OR A, E
	0xB4: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetH()); return 1 }, // OR A, H
	0xB5: func(cpu *CPU) uint8 { or8(cpu, cpu.regs.GetL()); return 1 }, // OR A, L
	// OR A,(HL)
	0xB6: func(cpu *CPU) uint8 { or8(cpu, cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // OR A, (HL)
	// OR A,n
	0xF6: func(cpu *CPU) uint8 { or8(cpu, cpu.readUint8()); return 2 }, // OR A,n

	// XOR A,r
	0xAF: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetA()); return 1 }, // XOR A, A
	0xA8: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetB()); return 1 }, // XOR A, B
	0xA9: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetC()); return 1 }, // XOR A, C
	0xAA: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetD()); return 1 }, // XOR A, D
	0xAB: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetE()); return 1 }, // XOR A, E
	0xAC: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetH()); return 1 }, // XOR A, H
	0xAD: func(cpu *CPU) uint8 { xor8(cpu, cpu.regs.GetL()); return 1 }, // XOR A, L
	// XOR A,(HL)
	0xAE: func(cpu *CPU) uint8 { xor8(cpu, cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // XOR A, (HL)
	// XOR A,n
	0xEE: func(cpu *CPU) uint8 { xor8(cpu, cpu.readUint8()); return 2 }, // XOR A,n

	// CP A,r
	0xBF: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetA()); return 1 }, // CP A, A
	0xB8: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetB()); return 1 }, // CP A, B
	0xB9: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetC()); return 1 }, // CP A, C
	0xBA: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetD()); return 1 }, // CP A, D
	0xBB: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetE()); return 1 }, // CP A, E
	0xBC: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetH()); return 1 }, // CP A, H
	0xBD: func(cpu *CPU) uint8 { cp8(cpu, cpu.regs.GetL()); return 1 }, // CP A, L
	// CP A,(HL)
	0xBE: func(cpu *CPU) uint8 { cp8(cpu, cpu.mmu.Read(cpu.regs.GetHL())); return 2 }, // CP A, (HL)
	// CP A,n
	0xFE: func(cpu *CPU) uint8 { cp8(cpu, cpu.readUint8()); return 2 }, // CP A,n

	// INC r
	0x3C: func(cpu *CPU) uint8 { cpu.regs.SetA(inc8(cpu, cpu.regs.GetA())); return 1 }, // INC A
	0x04: func(cpu *CPU) uint8 { cpu.regs.SetB(inc8(cpu, cpu.regs.GetB())); return 1 }, // INC B
	0x0C: func(cpu *CPU) uint8 { cpu.regs.SetC(inc8(cpu, cpu.regs.GetC())); return 1 }, // INC C
	0x14: func(cpu *CPU) uint8 { cpu.regs.SetD(inc8(cpu, cpu.regs.GetD())); return 1 }, // INC D
	0x1C: func(cpu *CPU) uint8 { cpu.regs.SetE(inc8(cpu, cpu.regs.GetE())); return 1 }, // INC E
	0x24: func(cpu *CPU) uint8 { cpu.regs.SetH(inc8(cpu, cpu.regs.GetH())); return 1 }, // INC H
	0x2C: func(cpu *CPU) uint8 { cpu.regs.SetL(inc8(cpu, cpu.regs.GetL())); return 1 }, // INC L

	// INC (HL)
	0x34: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, inc8(cpu, cpu.mmu.Read(addr)))
		return 3
	}, // INC (HL)

	// DEC r
	0x3D: func(cpu *CPU) uint8 { cpu.regs.SetA(dec8(cpu, cpu.regs.GetA())); return 1 }, // DEC A
	0x05: func(cpu *CPU) uint8 { cpu.regs.SetB(dec8(cpu, cpu.regs.GetB())); return 1 }, // DEC B
	0x0D: func(cpu *CPU) uint8 { cpu.regs.SetC(dec8(cpu, cpu.regs.GetC())); return 1 }, // DEC C
	0x15: func(cpu *CPU) uint8 { cpu.regs.SetD(dec8(cpu, cpu.regs.GetD())); return 1 }, // DEC D
	0x1D: func(cpu *CPU) uint8 { cpu.regs.SetE(dec8(cpu, cpu.regs.GetE())); return 1 }, // DEC E
	0x25: func(cpu *CPU) uint8 { cpu.regs.SetH(dec8(cpu, cpu.regs.GetH())); return 1 }, // DEC H
	0x2D: func(cpu *CPU) uint8 { cpu.regs.SetL(dec8(cpu, cpu.regs.GetL())); return 1 }, // DEC L

	// DEC (HL)
	0x35: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, dec8(cpu, cpu.mmu.Read(addr)))
		return 3
	}, // DEC (HL)

	///// 16-Bit ALU //////
	// ADD HL, rr
	0x09: func(cpu *CPU) uint8 { addHL(cpu, cpu.regs.GetBC()); return 2 }, // ADD HL, BC
	0x19: func(cpu *CPU) uint8 { addHL(cpu, cpu.regs.GetDE()); return 2 }, // ADD HL, DE
	0x29: func(cpu *CPU) uint8 { addHL(cpu, cpu.regs.GetHL()); return 2 }, // ADD HL, HL
	0x39: func(cpu *CPU) uint8 { addHL(cpu, cpu.regs.GetSP()); return 2 }, // ADD HL, SP

	// ADD SP, n
	0xE8: func(cpu *CPU) uint8 { slideSP(cpu, cpu.readUint8()); return 4 }, // ADD SP, n (TODO: CHECK CARRY)

	// INC rr
	0x03: func(cpu *CPU) uint8 { cpu.regs.SetBC(cpu.regs.GetBC() + 1); return 2 }, // INC BC
	0x13: func(cpu *CPU) uint8 { cpu.regs.SetDE(cpu.regs.GetDE() + 1); return 2 }, // INC DE
	0x23: func(cpu *CPU) uint8 { cpu.regs.SetHL(cpu.regs.GetHL() + 1); return 2 }, // INC HL
	0x33: func(cpu *CPU) uint8 { cpu.regs.SetSP(cpu.regs.GetSP() + 1); return 2 }, // INC SP

	// DEC rr
	0x0B: func(cpu *CPU) uint8 { cpu.regs.SetBC(cpu.regs.GetBC() - 1); return 2 }, // DEC BC
	0x1B: func(cpu *CPU) uint8 { cpu.regs.SetDE(cpu.regs.GetDE() - 1); return 2 }, // DEC DE
	0x2B: func(cpu *CPU) uint8 { cpu.regs.SetHL(cpu.regs.GetHL() - 1); return 2 }, // DEC HL
	0x3B: func(cpu *CPU) uint8 { cpu.regs.SetSP(cpu.regs.GetSP() - 1); return 2 }, // DEC SP

	///// Jumps //////
	// JP nn
	0xC3: func(cpu *CPU) uint8 { cpu.regs.SetPC(cpu.readUint16()); return 3 }, // JP nn

	// JP (HL)
	0xE9: func(cpu *CPU) uint8 { cpu.regs.SetPC(cpu.regs.GetHL()); return 1 }, // JP (HL)

	// JP cc, nn
	0xC2: func(cpu *CPU) uint8 {
		return jumpIF(cpu, !cpu.regs.GetFlag(registers.FlagZERO), cpu.readUint16())
	}, // JP NZ, nn
	0xCA: func(cpu *CPU) uint8 {
		return jumpIF(cpu, cpu.regs.GetFlag(registers.FlagZERO), cpu.readUint16())
	}, // JP Z, nn
	0xD2: func(cpu *CPU) uint8 {
		return jumpIF(cpu, !cpu.regs.GetFlag(registers.FlagCARRY), cpu.readUint16())
	}, // JP NC, nn
	0xDA: func(cpu *CPU) uint8 {
		return jumpIF(cpu, cpu.regs.GetFlag(registers.FlagCARRY), cpu.readUint16())
	}, // JP C, nn

	// JR n
	0x18: func(cpu *CPU) uint8 { jumpOffset(cpu, cpu.readUint8()); return 2 }, // JR n

	// JR cc, n
	0x20: func(cpu *CPU) uint8 {
		return jumpOffsetIF(cpu, !cpu.regs.GetFlag(registers.FlagZERO), cpu.readUint8())
	}, // JR NZ, n
	0x28: func(cpu *CPU) uint8 {
		return jumpOffsetIF(cpu, cpu.regs.GetFlag(registers.FlagZERO), cpu.readUint8())
	}, // JR Z, n
	0x30: func(cpu *CPU) uint8 {
		return jumpOffsetIF(cpu, !cpu.regs.GetFlag(registers.FlagCARRY), cpu.readUint8())
	}, // JR NC, n
	0x38: func(cpu *CPU) uint8 {
		return jumpOffsetIF(cpu, cpu.regs.GetFlag(registers.FlagCARRY), cpu.readUint8())
	}, // JR C, n

	///// Calls //////
	// CALL nn
	0xCD: func(cpu *CPU) uint8 { call(cpu, cpu.readUint16()); return 3 }, // CALL nn

	// CALL cc,nn
	0xC4: func(cpu *CPU) uint8 {
		return callIF(cpu, !cpu.regs.GetFlag(registers.FlagZERO), cpu.readUint16())
	}, // CALL NZ, nn
	0xCC: func(cpu *CPU) uint8 {
		return callIF(cpu, cpu.regs.GetFlag(registers.FlagZERO), cpu.readUint16())
	}, // CALL Z, nn
	0xD4: func(cpu *CPU) uint8 {
		return callIF(cpu, !cpu.regs.GetFlag(registers.FlagCARRY), cpu.readUint16())
	}, // CALL NC nn
	0xDC: func(cpu *CPU) uint8 {
		return callIF(cpu, cpu.regs.GetFlag(registers.FlagCARRY), cpu.readUint16())
	}, // CALL C, nn

	///// Restarts //////
	// RST n
	0xC7: func(cpu *CPU) uint8 { call(cpu, 0x0000); return 8 }, // RST 0x00
	0xCF: func(cpu *CPU) uint8 { call(cpu, 0x0008); return 8 }, // RST 0x08
	0xD7: func(cpu *CPU) uint8 { call(cpu, 0x0010); return 8 }, // RST 0x10
	0xDF: func(cpu *CPU) uint8 { call(cpu, 0x0018); return 8 }, // RST 0x18
	0xE7: func(cpu *CPU) uint8 { call(cpu, 0x0020); return 8 }, // RST 0x20
	0xEF: func(cpu *CPU) uint8 { call(cpu, 0x0028); return 8 }, // RST 0x28
	0xF7: func(cpu *CPU) uint8 { call(cpu, 0x0030); return 8 }, // RST 0x30
	0xFF: func(cpu *CPU) uint8 { call(cpu, 0x0038); return 8 }, // RST 0x38

	///// Returns //////
	// RET
	0xC9: func(cpu *CPU) uint8 { cpu.regs.SetPC(pop(cpu)); return 2 }, // RET

	// RET cc
	0xC0: func(cpu *CPU) uint8 { return retIF(cpu, !cpu.regs.GetFlag(registers.FlagZERO)) },  // RET NZ
	0xC8: func(cpu *CPU) uint8 { return retIF(cpu, cpu.regs.GetFlag(registers.FlagZERO)) },   // RET Z
	0xD0: func(cpu *CPU) uint8 { return retIF(cpu, !cpu.regs.GetFlag(registers.FlagCARRY)) }, // RET NC
	0xD8: func(cpu *CPU) uint8 { return retIF(cpu, cpu.regs.GetFlag(registers.FlagCARRY)) },  // RET C

	// RETI
	0xD9: func(cpu *CPU) uint8 { cpu.regs.SetPC(pop(cpu)); cpu.interrupts.EnableMaster(); return 2 }, // RETI

	///// Miscellaneous //////
	// DI
	0xF3: func(cpu *CPU) uint8 { cpu.interrupts.DisableMaster(); return 1 }, // DI
	// EI
	0xFB: func(cpu *CPU) uint8 { cpu.interrupts.EnableMaster(); return 1 }, // EI
	// DAA
	0x27: func(cpu *CPU) uint8 { daa(cpu); return 1 }, // DAA
	// CPL
	0x2F: func(cpu *CPU) uint8 { cpl(cpu); return 1 }, // CPL
	// CCF
	0x3F: func(cpu *CPU) uint8 { ccf(cpu); return 1 }, // CCF
	// SCF
	0x37: func(cpu *CPU) uint8 { scf(cpu); return 1 }, // SCF

	///// Rotates & Shifts //////
	// RLCA
	0x07: func(cpu *CPU) uint8 { rotateLeftA(cpu, false); return 1 }, // RLCA
	// RLA
	0x17: func(cpu *CPU) uint8 { rotateLeftA(cpu, true); return 1 }, // RLA
	// RRCA
	0x0F: func(cpu *CPU) uint8 { rotateRightA(cpu, false); return 1 }, // RRCA
	// RRA
	0x1F: func(cpu *CPU) uint8 { rotateRightA(cpu, true); return 1 }, // RRA

	// HALT
	0x76: func(cpu *CPU) uint8 { cpu.halt = true; return 1 }, // HALT

}

func scf(cpu *CPU) {
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, true)
}

func ccf(cpu *CPU) {
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, !cpu.regs.GetFlag(registers.FlagCARRY))
}

func cpl(cpu *CPU) {
	cpu.regs.SetA(^(cpu.regs.GetA()))
	cpu.regs.SetFlag(registers.FlagSUB, true)
	cpu.regs.SetFlag(registers.FlagHCARRY, true)
}

func daa(cpu *CPU) {
	a := cpu.regs.GetA()
	var adjust uint8
	if cpu.regs.GetFlag(registers.FlagCARRY) {
		adjust = 0x60
	}
	if cpu.regs.GetFlag(registers.FlagHCARRY) {
		adjust |= 0x06
	}
	if !cpu.regs.GetFlag(registers.FlagSUB) {
		if a&0x0F > 0x09 {
			adjust |= 0x06
		}
		if a > 0x99 {
			adjust |= 0x60
		}
		a += adjust
	} else {
		a -= adjust
	}

	cpu.regs.SetFlag(registers.FlagCARRY, adjust >= 0x60)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagZERO, a == 0)
	cpu.regs.SetA(a)
}

func rotateLeftA(cpu *CPU, withCarry bool) {
	a := cpu.regs.GetA()
	oldBit7 := (a >> 7) != 0
	if withCarry {
		a <<= 1
		if cpu.regs.GetFlag(registers.FlagCARRY) {
			a |= 0x01
		}
	} else {
		a = bits.RotateLeft8(a, 1)
	}
	cpu.regs.SetFlag(0xF0, false) // ALL FLAGS RESET
	cpu.regs.SetFlag(registers.FlagCARRY, oldBit7)

	cpu.regs.SetA(a)
}

func rotateRightA(cpu *CPU, withCarry bool) {
	a := cpu.regs.GetA()
	oldBit0 := (a & 1) != 0

	if withCarry {
		a >>= 1
		if cpu.regs.GetFlag(registers.FlagCARRY) {
			a |= 0x80
		}
	} else {
		a = bits.RotateLeft8(a, -1)
	}

	cpu.regs.SetFlag(0xF0, false) // ALL FLAGS RESET
	cpu.regs.SetFlag(registers.FlagCARRY, oldBit0)

	cpu.regs.SetA(a)
}

func retIF(cpu *CPU, cond bool) uint8 {
	if cond {
		cpu.regs.SetPC(pop(cpu))
		return 5
	}
	return 2
}

func jumpOffset(cpu *CPU, offset uint8) {
	pc := cpu.regs.GetPC()
	cpu.regs.SetPC(uint16(int32(pc) + int32(int8(offset))))
}

func jumpOffsetIF(cpu *CPU, cond bool, offset uint8) uint8 {
	if cond {
		jumpOffset(cpu, offset)
		return 3
	}
	return 2
}

func jumpIF(cpu *CPU, cond bool, dest uint16) uint8 {
	if cond {
		cpu.regs.SetPC(dest)
		return 4
	}
	return 3
}

func call(cpu *CPU, dest uint16) {
	push(cpu, cpu.regs.GetPC())
	cpu.regs.SetPC(dest)
}

func callIF(cpu *CPU, cond bool, dest uint16) uint8 {
	if cond {
		call(cpu, dest)
		return 6
	}
	return 3
}

func push(cpu *CPU, val uint16) {
	sp := cpu.regs.GetSP()
	cpu.mmu.Write(sp-1, uint8(val>>8))
	cpu.mmu.Write(sp-2, uint8(val&0x00FF))
	cpu.regs.SetSP(sp - 2)
}
func pop(cpu *CPU) uint16 {
	sp := cpu.regs.GetSP()

	res := uint16(cpu.mmu.Read(sp)) + // LO
		uint16(cpu.mmu.Read(sp+1))<<8 // HI
	cpu.regs.SetSP(sp + 2)
	return res
}

func add8(cpu *CPU, n uint8, withCarry bool) {
	var carry uint8
	if withCarry && cpu.regs.GetFlag(registers.FlagCARRY) {
		carry = 1
	}
	oldA := cpu.regs.GetA()
	a := oldA + n + carry
	cpu.regs.SetFlag(registers.FlagZERO, a == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, uint32(oldA&0xF)+uint32(n&0xF)+uint32(carry) > 0xF)
	cpu.regs.SetFlag(registers.FlagCARRY, uint32(oldA)+uint32(n)+uint32(carry) > 0xFF)
	cpu.regs.SetA(a)
}

func inc8(cpu *CPU, value uint8) uint8 {
	cpu.regs.SetFlag(registers.FlagZERO, value == 0xFF)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, value&0x0F == 0x0F)
	return value + 1
}

func sub8(cpu *CPU, n uint8, withCarry bool) {
	var carry uint8
	if withCarry && cpu.regs.GetFlag(registers.FlagCARRY) {
		carry = 1
	}
	oldA := cpu.regs.GetA()
	a := oldA - n - carry
	cpu.regs.SetFlag(registers.FlagZERO, a == 0)
	cpu.regs.SetFlag(registers.FlagSUB, true)
	cpu.regs.SetFlag(registers.FlagHCARRY, int32(oldA&0xF)-int32(n&0xF)-int32(carry) < 0)
	cpu.regs.SetFlag(registers.FlagCARRY, int32(oldA)-int32(n)-int32(carry) < 0)
	cpu.regs.SetA(a)
}

func dec8(cpu *CPU, value uint8) uint8 {
	cpu.regs.SetFlag(registers.FlagZERO, value == 0x01)
	cpu.regs.SetFlag(registers.FlagSUB, true)
	cpu.regs.SetFlag(registers.FlagHCARRY, value&0x0F == 0x00)
	return value - 1
}

func cp8(cpu *CPU, n uint8) {
	a := cpu.regs.GetA()
	sub8(cpu, n, false)
	cpu.regs.SetA(a)
}

func and8(cpu *CPU, n uint8) {
	a := cpu.regs.GetA() & n
	cpu.regs.SetFlag(registers.FlagZERO, a == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, true)
	cpu.regs.SetFlag(registers.FlagCARRY, false)
	cpu.regs.SetA(a)
}

func or8(cpu *CPU, n uint8) {
	a := cpu.regs.GetA() | n
	cpu.regs.SetFlag(registers.FlagZERO, a == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, false)
	cpu.regs.SetA(a)
}

func xor8(cpu *CPU, n uint8) {
	a := cpu.regs.GetA() ^ n
	cpu.regs.SetFlag(registers.FlagZERO, a == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, false)
	cpu.regs.SetA(a)
}

func addHL(cpu *CPU, n uint16) {
	oldHL := cpu.regs.GetHL()
	res := oldHL + n
	cpu.regs.SetHL(res)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	// cpu.regs.SetFlag(registers.FlagHCARRY, (uint32(oldHL&0x0FFF)+uint32(n&0x0FFF))&0x01000 != 0)
	cpu.regs.SetFlag(registers.FlagHCARRY, ((oldHL&n)|((oldHL|n)&^res))>>11&0x01 == 1)
	// cpu.regs.SetFlag(registers.FlagCARRY, oldHL > 0xFFFF-n)
	cpu.regs.SetFlag(registers.FlagCARRY, ((oldHL&n)|((oldHL|n)&^res))>>15 == 1)
}

func slideSP(cpu *CPU, n uint8) {
	signedN := int32(int8(n))
	sp := cpu.regs.GetSP()
	result := uint16(int32(sp) + signedN)
	cpu.regs.SetFlag(registers.FlagZERO, false)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, (result&0xF) < (sp&0xF))
	cpu.regs.SetFlag(registers.FlagCARRY, (result&0xFF) < (sp&0xFF))
	cpu.regs.SetSP(result)
}
