package registers

const (
	FlagZERO   uint8 = 1 << 7
	FlagSUB    uint8 = 1 << 6
	FlagHCARRY uint8 = 1 << 5
	FlagCARRY  uint8 = 1 << 4
)

// Registers manage CPU registers values
//
// 8 bit registers can be combined to 16 bit registers
type Registers struct {
	a uint8
	f uint8
	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8

	sp uint16
	pc uint16
}

func (r *Registers) GetFlag(mask uint8) bool { return r.f&mask != 0 }
func (r *Registers) SetFlag(mask uint8, set bool) {
	if set {
		r.f |= mask
	} else {
		r.f &= ^mask
	}
}

func (r *Registers) GetAF() uint16 { return uint16(r.a)<<8 | uint16(r.f) }
func (r *Registers) GetBC() uint16 { return uint16(r.b)<<8 | uint16(r.c) }
func (r *Registers) GetDE() uint16 { return uint16(r.d)<<8 | uint16(r.e) }
func (r *Registers) GetHL() uint16 { return uint16(r.h)<<8 | uint16(r.l) }
func (r *Registers) GetSP() uint16 { return r.sp }
func (r *Registers) GetPC() uint16 { return r.pc }

func (r *Registers) SetAF(v uint16) { r.a = uint8(v >> 8); r.f = uint8(v & 0x00FF) }
func (r *Registers) SetBC(v uint16) { r.b = uint8(v >> 8); r.c = uint8(v & 0x00FF) }
func (r *Registers) SetDE(v uint16) { r.d = uint8(v >> 8); r.e = uint8(v & 0x00FF) }
func (r *Registers) SetHL(v uint16) { r.h = uint8(v >> 8); r.l = uint8(v & 0x00FF) }
func (r *Registers) SetSP(v uint16) { r.sp = v }
func (r *Registers) SetPC(v uint16) { r.pc = v }

func (r *Registers) GetA() uint8 { return r.a }
func (r *Registers) GetF() uint8 { return r.f }
func (r *Registers) GetB() uint8 { return r.b }
func (r *Registers) GetC() uint8 { return r.c }
func (r *Registers) GetD() uint8 { return r.d }
func (r *Registers) GetE() uint8 { return r.e }
func (r *Registers) GetH() uint8 { return r.h }
func (r *Registers) GetL() uint8 { return r.l }

func (r *Registers) SetA(v uint8) { r.a = v }
func (r *Registers) SetF(v uint8) { r.f = v }
func (r *Registers) SetB(v uint8) { r.b = v }
func (r *Registers) SetC(v uint8) { r.c = v }
func (r *Registers) SetD(v uint8) { r.d = v }
func (r *Registers) SetE(v uint8) { r.e = v }
func (r *Registers) SetH(v uint8) { r.h = v }
func (r *Registers) SetL(v uint8) { r.l = v }
