package cpu

import (
	"github.com/jmontupet/gbcore/internal/pkg/cpu/registers"
)

var instructionCBList = instructions{
	// RL n
	0x17: func(cpu *CPU) uint8 { cpu.regs.SetA(rotateLeft(cpu, cpu.regs.GetA(), true)); return 2 }, // RL A
	0x10: func(cpu *CPU) uint8 { cpu.regs.SetB(rotateLeft(cpu, cpu.regs.GetB(), true)); return 2 }, // RL B
	0x11: func(cpu *CPU) uint8 { cpu.regs.SetC(rotateLeft(cpu, cpu.regs.GetC(), true)); return 2 }, // RL C
	0x12: func(cpu *CPU) uint8 { cpu.regs.SetD(rotateLeft(cpu, cpu.regs.GetD(), true)); return 2 }, // RL D
	0x13: func(cpu *CPU) uint8 { cpu.regs.SetE(rotateLeft(cpu, cpu.regs.GetE(), true)); return 2 }, // RL E
	0x14: func(cpu *CPU) uint8 { cpu.regs.SetH(rotateLeft(cpu, cpu.regs.GetH(), true)); return 2 }, // RL H
	0x15: func(cpu *CPU) uint8 { cpu.regs.SetL(rotateLeft(cpu, cpu.regs.GetL(), true)); return 2 }, // RL L
	0x16: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, rotateLeft(cpu, cpu.mmu.Read(addr), true))
		return 4
	}, // RL (HL)

	// RLC n
	0x07: func(cpu *CPU) uint8 { cpu.regs.SetA(rotateLeft(cpu, cpu.regs.GetA(), false)); return 2 }, // RLC A
	0x00: func(cpu *CPU) uint8 { cpu.regs.SetB(rotateLeft(cpu, cpu.regs.GetB(), false)); return 2 }, // RLC B
	0x01: func(cpu *CPU) uint8 { cpu.regs.SetC(rotateLeft(cpu, cpu.regs.GetC(), false)); return 2 }, // RLC C
	0x02: func(cpu *CPU) uint8 { cpu.regs.SetD(rotateLeft(cpu, cpu.regs.GetD(), false)); return 2 }, // RLC D
	0x03: func(cpu *CPU) uint8 { cpu.regs.SetE(rotateLeft(cpu, cpu.regs.GetE(), false)); return 2 }, // RLC E
	0x04: func(cpu *CPU) uint8 { cpu.regs.SetH(rotateLeft(cpu, cpu.regs.GetH(), false)); return 2 }, // RLC H
	0x05: func(cpu *CPU) uint8 { cpu.regs.SetL(rotateLeft(cpu, cpu.regs.GetL(), false)); return 2 }, // RLC L
	0x06: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, rotateLeft(cpu, cpu.mmu.Read(addr), false))
		return 4
	}, // RLC (HL)

	// RRC n
	0x0F: func(cpu *CPU) uint8 { cpu.regs.SetA(rotateRight(cpu, cpu.regs.GetA(), false)); return 2 }, // RRC A
	0x08: func(cpu *CPU) uint8 { cpu.regs.SetB(rotateRight(cpu, cpu.regs.GetB(), false)); return 2 }, // RRC B
	0x09: func(cpu *CPU) uint8 { cpu.regs.SetC(rotateRight(cpu, cpu.regs.GetC(), false)); return 2 }, // RRC C
	0x0A: func(cpu *CPU) uint8 { cpu.regs.SetD(rotateRight(cpu, cpu.regs.GetD(), false)); return 2 }, // RRC D
	0x0B: func(cpu *CPU) uint8 { cpu.regs.SetE(rotateRight(cpu, cpu.regs.GetE(), false)); return 2 }, // RRC E
	0x0C: func(cpu *CPU) uint8 { cpu.regs.SetH(rotateRight(cpu, cpu.regs.GetH(), false)); return 2 }, // RRC H
	0x0D: func(cpu *CPU) uint8 { cpu.regs.SetL(rotateRight(cpu, cpu.regs.GetL(), false)); return 2 }, // RRC L
	0x0E: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, rotateRight(cpu, cpu.mmu.Read(addr), false))
		return 4
	}, // RRC (HL)

	// RR n
	0x1F: func(cpu *CPU) uint8 { cpu.regs.SetA(rotateRight(cpu, cpu.regs.GetA(), true)); return 2 }, // RR A
	0x18: func(cpu *CPU) uint8 { cpu.regs.SetB(rotateRight(cpu, cpu.regs.GetB(), true)); return 2 }, // RR B
	0x19: func(cpu *CPU) uint8 { cpu.regs.SetC(rotateRight(cpu, cpu.regs.GetC(), true)); return 2 }, // RR C
	0x1A: func(cpu *CPU) uint8 { cpu.regs.SetD(rotateRight(cpu, cpu.regs.GetD(), true)); return 2 }, // RR D
	0x1B: func(cpu *CPU) uint8 { cpu.regs.SetE(rotateRight(cpu, cpu.regs.GetE(), true)); return 2 }, // RR E
	0x1C: func(cpu *CPU) uint8 { cpu.regs.SetH(rotateRight(cpu, cpu.regs.GetH(), true)); return 2 }, // RR H
	0x1D: func(cpu *CPU) uint8 { cpu.regs.SetL(rotateRight(cpu, cpu.regs.GetL(), true)); return 2 }, // RR L
	0x1E: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, rotateRight(cpu, cpu.mmu.Read(addr), true))
		return 4
	}, // RR (HL)

	// BIT b,r
	0x47: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit0); return 2 },                // BIT 0,A
	0x40: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit0); return 2 },                // BIT 0,B
	0x41: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit0); return 2 },                // BIT 0,C
	0x42: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit0); return 2 },                // BIT 0,D
	0x43: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit0); return 2 },                // BIT 0,E
	0x44: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit0); return 2 },                // BIT 0,H
	0x45: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit0); return 2 },                // BIT 0,L
	0x46: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit0); return 4 }, // BIT 0,(HL)

	0x4F: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit1); return 2 },                // BIT 1,A
	0x48: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit1); return 2 },                // BIT 1,B
	0x49: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit1); return 2 },                // BIT 1,C
	0x4A: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit1); return 2 },                // BIT 1,D
	0x4B: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit1); return 2 },                // BIT 1,E
	0x4C: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit1); return 2 },                // BIT 1,H
	0x4D: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit1); return 2 },                // BIT 1,L
	0x4E: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit1); return 4 }, // BIT 1,(HL)

	0x57: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit2); return 2 },                // BIT 2,A
	0x50: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit2); return 2 },                // BIT 2,B
	0x51: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit2); return 2 },                // BIT 2,C
	0x52: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit2); return 2 },                // BIT 2,D
	0x53: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit2); return 2 },                // BIT 2,E
	0x54: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit2); return 2 },                // BIT 2,H
	0x55: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit2); return 2 },                // BIT 2,L
	0x56: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit2); return 4 }, // BIT 2,(HL)

	0x5F: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit3); return 2 },                // BIT 3,A
	0x58: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit3); return 2 },                // BIT 3,B
	0x59: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit3); return 2 },                // BIT 3,C
	0x5A: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit3); return 2 },                // BIT 3,D
	0x5B: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit3); return 2 },                // BIT 3,E
	0x5C: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit3); return 2 },                // BIT 3,H
	0x5D: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit3); return 2 },                // BIT 3,L
	0x5E: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit3); return 4 }, // BIT 3,(HL)

	0x67: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit4); return 2 },                // BIT 4,A
	0x60: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit4); return 2 },                // BIT 4,B
	0x61: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit4); return 2 },                // BIT 4,C
	0x62: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit4); return 2 },                // BIT 4,D
	0x63: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit4); return 2 },                // BIT 4,E
	0x64: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit4); return 2 },                // BIT 4,H
	0x65: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit4); return 2 },                // BIT 4,L
	0x66: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit4); return 4 }, // BIT 4,(HL)

	0x6F: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit5); return 2 },                // BIT 5,A
	0x68: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit5); return 2 },                // BIT 5,B
	0x69: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit5); return 2 },                // BIT 5,C
	0x6A: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit5); return 2 },                // BIT 5,D
	0x6B: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit5); return 2 },                // BIT 5,E
	0x6C: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit5); return 2 },                // BIT 5,H
	0x6D: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit5); return 2 },                // BIT 5,L
	0x6E: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit5); return 4 }, // BIT 5,(HL)

	0x77: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit6); return 2 },                // BIT 6,A
	0x70: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit6); return 2 },                // BIT 6,B
	0x71: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit6); return 2 },                // BIT 6,C
	0x72: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit6); return 2 },                // BIT 6,D
	0x73: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit6); return 2 },                // BIT 6,E
	0x74: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit6); return 2 },                // BIT 6,H
	0x75: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit6); return 2 },                // BIT 6,L
	0x76: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit6); return 4 }, // BIT 6,(HL)

	0x7F: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetA(), Bit7); return 2 },                // BIT 7,A
	0x78: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetB(), Bit7); return 2 },                // BIT 7,B
	0x79: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetC(), Bit7); return 2 },                // BIT 7,C
	0x7A: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetD(), Bit7); return 2 },                // BIT 7,D
	0x7B: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetE(), Bit7); return 2 },                // BIT 7,E
	0x7C: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetH(), Bit7); return 2 },                // BIT 7,H
	0x7D: func(cpu *CPU) uint8 { bit(cpu, cpu.regs.GetL(), Bit7); return 2 },                // BIT 7,L
	0x7E: func(cpu *CPU) uint8 { bit(cpu, cpu.mmu.Read(cpu.regs.GetHL()), Bit7); return 4 }, // BIT 7,(HL)

	// RES b,r
	0x87: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit0)); return 2 }, // RES 0,A
	0x80: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit0)); return 2 }, // RES 0,B
	0x81: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit0)); return 2 }, // RES 0,C
	0x82: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit0)); return 2 }, // RES 0,D
	0x83: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit0)); return 2 }, // RES 0,E
	0x84: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit0)); return 2 }, // RES 0,H
	0x85: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit0)); return 2 }, // RES 0,L
	0x86: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit0))
		return 4
	}, // RES 0,(HL)

	0x8F: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit1)); return 2 }, // RES 1,A
	0x88: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit1)); return 2 }, // RES 1,B
	0x89: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit1)); return 2 }, // RES 1,C
	0x8A: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit1)); return 2 }, // RES 1,D
	0x8B: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit1)); return 2 }, // RES 1,E
	0x8C: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit1)); return 2 }, // RES 1,H
	0x8D: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit1)); return 2 }, // RES 1,L
	0x8E: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit1))
		return 4
	}, // RES 1,(HL)

	0x97: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit2)); return 2 }, // RES 2,A
	0x90: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit2)); return 2 }, // RES 2,B
	0x91: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit2)); return 2 }, // RES 2,C
	0x92: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit2)); return 2 }, // RES 2,D
	0x93: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit2)); return 2 }, // RES 2,E
	0x94: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit2)); return 2 }, // RES 2,H
	0x95: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit2)); return 2 }, // RES 2,L
	0x96: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit2))
		return 4
	}, // RES 2,(HL)

	0x9F: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit3)); return 2 }, // RES 3,A
	0x98: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit3)); return 2 }, // RES 3,B
	0x99: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit3)); return 2 }, // RES 3,C
	0x9A: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit3)); return 2 }, // RES 3,D
	0x9B: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit3)); return 2 }, // RES 3,E
	0x9C: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit3)); return 2 }, // RES 3,H
	0x9D: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit3)); return 2 }, // RES 3,L
	0x9E: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit3))
		return 4
	}, // RES 3,(HL)

	0xA7: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit4)); return 2 }, // RES 4,A
	0xA0: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit4)); return 2 }, // RES 4,B
	0xA1: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit4)); return 2 }, // RES 4,C
	0xA2: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit4)); return 2 }, // RES 4,D
	0xA3: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit4)); return 2 }, // RES 4,E
	0xA4: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit4)); return 2 }, // RES 4,H
	0xA5: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit4)); return 2 }, // RES 4,L
	0xA6: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit4))
		return 4
	}, // RES 4,(HL)

	0xAF: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit5)); return 2 }, // RES 5,A
	0xA8: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit5)); return 2 }, // RES 5,B
	0xA9: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit5)); return 2 }, // RES 5,C
	0xAA: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit5)); return 2 }, // RES 5,D
	0xAB: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit5)); return 2 }, // RES 5,E
	0xAC: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit5)); return 2 }, // RES 5,H
	0xAD: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit5)); return 2 }, // RES 5,L
	0xAE: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit5))
		return 4
	}, // RES 5,(HL)

	0xB7: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit6)); return 2 }, // RES 6,A
	0xB0: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit6)); return 2 }, // RES 6,B
	0xB1: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit6)); return 2 }, // RES 6,C
	0xB2: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit6)); return 2 }, // RES 6,D
	0xB3: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit6)); return 2 }, // RES 6,E
	0xB4: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit6)); return 2 }, // RES 6,H
	0xB5: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit6)); return 2 }, // RES 6,L
	0xB6: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit6))
		return 4
	}, // RES 6,(HL)

	0xBF: func(cpu *CPU) uint8 { cpu.regs.SetA(res(cpu, cpu.regs.GetA(), Bit7)); return 2 }, // RES 7,A
	0xB8: func(cpu *CPU) uint8 { cpu.regs.SetB(res(cpu, cpu.regs.GetB(), Bit7)); return 2 }, // RES 7,B
	0xB9: func(cpu *CPU) uint8 { cpu.regs.SetC(res(cpu, cpu.regs.GetC(), Bit7)); return 2 }, // RES 7,C
	0xBA: func(cpu *CPU) uint8 { cpu.regs.SetD(res(cpu, cpu.regs.GetD(), Bit7)); return 2 }, // RES 7,D
	0xBB: func(cpu *CPU) uint8 { cpu.regs.SetE(res(cpu, cpu.regs.GetE(), Bit7)); return 2 }, // RES 7,E
	0xBC: func(cpu *CPU) uint8 { cpu.regs.SetH(res(cpu, cpu.regs.GetH(), Bit7)); return 2 }, // RES 7,H
	0xBD: func(cpu *CPU) uint8 { cpu.regs.SetL(res(cpu, cpu.regs.GetL(), Bit7)); return 2 }, // RES 7,L
	0xBE: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, res(cpu, cpu.mmu.Read(addr), Bit7))
		return 4
	}, // RES 7,(HL)

	// SET b,r
	0xC7: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit0)); return 2 }, // SET 0,A
	0xC0: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit0)); return 2 }, // SET 0,B
	0xC1: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit0)); return 2 }, // SET 0,C
	0xC2: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit0)); return 2 }, // SET 0,D
	0xC3: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit0)); return 2 }, // SET 0,E
	0xC4: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit0)); return 2 }, // SET 0,H
	0xC5: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit0)); return 2 }, // SET 0,L
	0xC6: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit0))
		return 4
	}, // SET 0,(HL)

	0xCF: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit1)); return 2 }, // RES 1,A
	0xC8: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit1)); return 2 }, // RES 1,B
	0xC9: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit1)); return 2 }, // RES 1,C
	0xCA: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit1)); return 2 }, // RES 1,D
	0xCB: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit1)); return 2 }, // RES 1,E
	0xCC: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit1)); return 2 }, // RES 1,H
	0xCD: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit1)); return 2 }, // RES 1,L
	0xCE: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit1))
		return 4
	}, // SET 1,(HL)

	0xD7: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit2)); return 2 }, // SET 2,A
	0xD0: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit2)); return 2 }, // SET 2,B
	0xD1: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit2)); return 2 }, // SET 2,C
	0xD2: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit2)); return 2 }, // SET 2,D
	0xD3: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit2)); return 2 }, // SET 2,E
	0xD4: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit2)); return 2 }, // SET 2,H
	0xD5: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit2)); return 2 }, // SET 2,L
	0xD6: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit2))
		return 4
	}, // SET 2,(HL)

	0xDF: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit3)); return 2 }, // SET 3,A
	0xD8: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit3)); return 2 }, // SET 3,B
	0xD9: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit3)); return 2 }, // SET 3,C
	0xDA: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit3)); return 2 }, // SET 3,D
	0xDB: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit3)); return 2 }, // SET 3,E
	0xDC: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit3)); return 2 }, // SET 3,H
	0xDD: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit3)); return 2 }, // SET 3,L
	0xDE: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit3))
		return 4
	}, // SET 3,(HL)

	0xE7: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit4)); return 2 }, // SET 4,A
	0xE0: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit4)); return 2 }, // SET 4,B
	0xE1: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit4)); return 2 }, // SET 4,C
	0xE2: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit4)); return 2 }, // SET 4,D
	0xE3: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit4)); return 2 }, // SET 4,E
	0xE4: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit4)); return 2 }, // SET 4,H
	0xE5: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit4)); return 2 }, // SET 4,L
	0xE6: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit4))
		return 4
	}, // SET 4,(HL)

	0xEF: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit5)); return 2 }, // SET 5,A
	0xE8: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit5)); return 2 }, // SET 5,B
	0xE9: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit5)); return 2 }, // SET 5,C
	0xEA: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit5)); return 2 }, // SET 5,D
	0xEB: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit5)); return 2 }, // SET 5,E
	0xEC: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit5)); return 2 }, // SET 5,H
	0xED: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit5)); return 2 }, // SET 5,L
	0xEE: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit5))
		return 4
	}, // SET 5,(HL)

	0xF7: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit6)); return 2 }, // SET 6,A
	0xF0: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit6)); return 2 }, // SET 6,B
	0xF1: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit6)); return 2 }, // SET 6,C
	0xF2: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit6)); return 2 }, // SET 6,D
	0xF3: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit6)); return 2 }, // SET 6,E
	0xF4: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit6)); return 2 }, // SET 6,H
	0xF5: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit6)); return 2 }, // SET 6,L
	0xF6: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit6))
		return 4
	}, // SET 6,(HL)

	0xFF: func(cpu *CPU) uint8 { cpu.regs.SetA(set(cpu, cpu.regs.GetA(), Bit7)); return 2 }, // SET 7,A
	0xF8: func(cpu *CPU) uint8 { cpu.regs.SetB(set(cpu, cpu.regs.GetB(), Bit7)); return 2 }, // SET 7,B
	0xF9: func(cpu *CPU) uint8 { cpu.regs.SetC(set(cpu, cpu.regs.GetC(), Bit7)); return 2 }, // SET 7,C
	0xFA: func(cpu *CPU) uint8 { cpu.regs.SetD(set(cpu, cpu.regs.GetD(), Bit7)); return 2 }, // SET 7,D
	0xFB: func(cpu *CPU) uint8 { cpu.regs.SetE(set(cpu, cpu.regs.GetE(), Bit7)); return 2 }, // SET 7,E
	0xFC: func(cpu *CPU) uint8 { cpu.regs.SetH(set(cpu, cpu.regs.GetH(), Bit7)); return 2 }, // SET 7,H
	0xFD: func(cpu *CPU) uint8 { cpu.regs.SetL(set(cpu, cpu.regs.GetL(), Bit7)); return 2 }, // SET 7,L
	0xFE: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, set(cpu, cpu.mmu.Read(addr), Bit7))
		return 4
	}, // SET 7,(HL)

	// SWAP n
	0x37: func(cpu *CPU) uint8 { cpu.regs.SetA(swap(cpu, cpu.regs.GetA())); return 2 }, // SWAP A
	0x30: func(cpu *CPU) uint8 { cpu.regs.SetB(swap(cpu, cpu.regs.GetB())); return 2 }, // SWAP B
	0x31: func(cpu *CPU) uint8 { cpu.regs.SetC(swap(cpu, cpu.regs.GetC())); return 2 }, // SWAP C
	0x32: func(cpu *CPU) uint8 { cpu.regs.SetD(swap(cpu, cpu.regs.GetD())); return 2 }, // SWAP D
	0x33: func(cpu *CPU) uint8 { cpu.regs.SetE(swap(cpu, cpu.regs.GetE())); return 2 }, // SWAP E
	0x34: func(cpu *CPU) uint8 { cpu.regs.SetH(swap(cpu, cpu.regs.GetH())); return 2 }, // SWAP H
	0x35: func(cpu *CPU) uint8 { cpu.regs.SetL(swap(cpu, cpu.regs.GetL())); return 2 }, // SWAP L
	0x36: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, swap(cpu, cpu.mmu.Read(addr)))
		return 4
	}, // SWAP (HL)

	// SLA n
	0x27: func(cpu *CPU) uint8 { cpu.regs.SetA(shiftLeft(cpu, cpu.regs.GetA())); return 2 }, // SLA A
	0x20: func(cpu *CPU) uint8 { cpu.regs.SetB(shiftLeft(cpu, cpu.regs.GetB())); return 2 }, // SLA B
	0x21: func(cpu *CPU) uint8 { cpu.regs.SetC(shiftLeft(cpu, cpu.regs.GetC())); return 2 }, // SLA C
	0x22: func(cpu *CPU) uint8 { cpu.regs.SetD(shiftLeft(cpu, cpu.regs.GetD())); return 2 }, // SLA D
	0x23: func(cpu *CPU) uint8 { cpu.regs.SetE(shiftLeft(cpu, cpu.regs.GetE())); return 2 }, // SLA E
	0x24: func(cpu *CPU) uint8 { cpu.regs.SetH(shiftLeft(cpu, cpu.regs.GetH())); return 2 }, // SLA H
	0x25: func(cpu *CPU) uint8 { cpu.regs.SetL(shiftLeft(cpu, cpu.regs.GetL())); return 2 }, // SLA L
	0x26: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, shiftLeft(cpu, cpu.mmu.Read(addr)))
		return 4
	}, // SLA (HL)

	// SRL n
	0x3F: func(cpu *CPU) uint8 { cpu.regs.SetA(shiftRight(cpu, cpu.regs.GetA(), false)); return 2 }, // SRL A
	0x38: func(cpu *CPU) uint8 { cpu.regs.SetB(shiftRight(cpu, cpu.regs.GetB(), false)); return 2 }, // SRL B
	0x39: func(cpu *CPU) uint8 { cpu.regs.SetC(shiftRight(cpu, cpu.regs.GetC(), false)); return 2 }, // SRL C
	0x3A: func(cpu *CPU) uint8 { cpu.regs.SetD(shiftRight(cpu, cpu.regs.GetD(), false)); return 2 }, // SRL D
	0x3B: func(cpu *CPU) uint8 { cpu.regs.SetE(shiftRight(cpu, cpu.regs.GetE(), false)); return 2 }, // SRL E
	0x3C: func(cpu *CPU) uint8 { cpu.regs.SetH(shiftRight(cpu, cpu.regs.GetH(), false)); return 2 }, // SRL H
	0x3D: func(cpu *CPU) uint8 { cpu.regs.SetL(shiftRight(cpu, cpu.regs.GetL(), false)); return 2 }, // SRL L
	0x3E: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, shiftRight(cpu, cpu.mmu.Read(addr), false))
		return 4
	}, // SRL (HL)

	// SRA n
	0x2F: func(cpu *CPU) uint8 { cpu.regs.SetA(shiftRight(cpu, cpu.regs.GetA(), true)); return 2 }, // SRA A
	0x28: func(cpu *CPU) uint8 { cpu.regs.SetB(shiftRight(cpu, cpu.regs.GetB(), true)); return 2 }, // SRA B
	0x29: func(cpu *CPU) uint8 { cpu.regs.SetC(shiftRight(cpu, cpu.regs.GetC(), true)); return 2 }, // SRA C
	0x2A: func(cpu *CPU) uint8 { cpu.regs.SetD(shiftRight(cpu, cpu.regs.GetD(), true)); return 2 }, // SRA D
	0x2B: func(cpu *CPU) uint8 { cpu.regs.SetE(shiftRight(cpu, cpu.regs.GetE(), true)); return 2 }, // SRA E
	0x2C: func(cpu *CPU) uint8 { cpu.regs.SetH(shiftRight(cpu, cpu.regs.GetH(), true)); return 2 }, // SRA H
	0x2D: func(cpu *CPU) uint8 { cpu.regs.SetL(shiftRight(cpu, cpu.regs.GetL(), true)); return 2 }, // SRA L
	0x2E: func(cpu *CPU) uint8 {
		addr := cpu.regs.GetHL()
		cpu.mmu.Write(addr, shiftRight(cpu, cpu.mmu.Read(addr), true))
		return 4
	}, // SRA (HL)
}

const (
	Bit0 uint8 = 1 << iota
	Bit1 uint8 = 1 << iota
	Bit2 uint8 = 1 << iota
	Bit3 uint8 = 1 << iota
	Bit4 uint8 = 1 << iota
	Bit5 uint8 = 1 << iota
	Bit6 uint8 = 1 << iota
	Bit7 uint8 = 1 << iota
)

// Shift n left into Carry. LSB of n set to 0.
func shiftLeft(cpu *CPU, n uint8) uint8 {
	result := n << 1
	bit7 := n >> 7
	cpu.regs.SetFlag(registers.FlagZERO, result == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, bit7 == 1)
	return result
}

// Shift Right + bit 0 to carry
func shiftRight(cpu *CPU, n uint8, keepBit7 bool) uint8 {
	result := n >> 1
	if keepBit7 {
		result |= (n & 0x80)
	}
	cpu.regs.SetFlag(registers.FlagZERO, result == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, n&1 == 1)
	return result
}

func swap(cpu *CPU, value uint8) uint8 {
	result := value<<4 | value>>4
	cpu.regs.SetFlag(registers.FlagZERO, result == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagCARRY, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	return result
}

func res(cpu *CPU, value uint8, mask uint8) uint8 {
	return value & ^mask
}

func set(cpu *CPU, value uint8, mask uint8) uint8 {
	return value | mask
}

func bit(cpu *CPU, value uint8, mask uint8) {
	cpu.regs.SetFlag(registers.FlagZERO, value&mask == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, true)
}

func rotateLeft(cpu *CPU, value uint8, withCarry bool) uint8 {
	oldBit7 := (value >> 7) != 0
	res := value << 1
	if withCarry {
		if cpu.regs.GetFlag(registers.FlagCARRY) {
			res |= 0x01
		}
	} else if oldBit7 {
		res |= 0x01
	}

	cpu.regs.SetFlag(registers.FlagZERO, res == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, oldBit7)

	return res
}

func rotateRight(cpu *CPU, value uint8, withCarry bool) uint8 {
	oldBit0 := (value & 1) != 0
	res := value >> 1
	if withCarry {
		if cpu.regs.GetFlag(registers.FlagCARRY) {
			res |= 0x80
		}
	} else if oldBit0 {
		res |= 0x80
	}

	cpu.regs.SetFlag(registers.FlagZERO, res == 0)
	cpu.regs.SetFlag(registers.FlagSUB, false)
	cpu.regs.SetFlag(registers.FlagHCARRY, false)
	cpu.regs.SetFlag(registers.FlagCARRY, oldBit0)

	return res
}
