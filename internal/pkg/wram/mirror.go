package wram

const mirrorOffset = -0x2000

type TranslatedWram struct {
	memory *WRam
}

func (m *TranslatedWram) translateAddr(addr uint16) uint16 {
	return uint16(int32(addr) + mirrorOffset)
}

func (m *TranslatedWram) Read(addr uint16) uint8 {
	return m.memory.Read(m.translateAddr(addr))
}
func (m *TranslatedWram) Write(addr uint16, value uint8) {
	m.memory.Write(m.translateAddr(addr), value)
}

func NewTranlatedWram(ram *WRam) *TranslatedWram {
	return &TranslatedWram{
		memory: ram,
	}
}
