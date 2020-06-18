package cpu

import (
	"testing"

	"github.com/jmontupet/gbcore/internal/pkg/cpu/registers"
	"github.com/stretchr/testify/assert"
)

type testMem struct {
	main  []uint8
	stack [][]uint8
	t     *testing.T
}

func (m *testMem) Push(newMem []uint8) {
	m.stack = append(m.stack, m.main)
	m.main = newMem
}

func (m *testMem) Pop() (old []uint8) {
	if len(m.stack) == 0 {
		m.t.Fatalf("TRY TO POP FROM EMPTY MEM STACK")
	}
	old = m.main
	m.main, m.stack = m.stack[len(m.stack)-1], m.stack[:len(m.stack)-1]
	return old
}

func newTestMem(t *testing.T) *testMem {
	main := [0x10000]uint8{}
	return &testMem{
		main:  main[0:],
		stack: [][]uint8{},
		t:     t,
	}
}

func (m *testMem) Read(addr uint16) uint8 {
	if int(addr) >= len(m.main) {
		m.t.Fatalf("TRY TO READ ADDRESS %X (max %X)", addr, len(m.main)-1)
	}
	return m.main[addr]
}
func (m *testMem) Write(addr uint16, value uint8) {
	if int(addr) >= len(m.main) {
		m.t.Errorf("TRY TO WRITE ADDRESS %X (max %X)", addr, len(m.main)-1)
	}
	m.main[addr] = value
}

func (m *testMem) Clear(_ *testing.T, value uint8) {
	for i := range m.main {
		m.main[i] = value
	}
}

type testInterrupts struct{}

func (*testInterrupts) Read(uint16) uint8   { return 0x00 }
func (*testInterrupts) Write(uint16, uint8) {}
func (*testInterrupts) GetNext() uint16     { return 0x00 }
func (*testInterrupts) EnableMaster()       {}
func (*testInterrupts) DisableMaster()      {}

func clearRegs(_ *testing.T, c *CPU) {
	c.regs.SetAF(0)
	c.regs.SetBC(0)
	c.regs.SetDE(0)
	c.regs.SetHL(0)
	c.regs.SetSP(0)
	c.regs.SetPC(0)
}

func initTestRegs(src *testRegsValue, dest *registers.Registers) {
	dest.SetA(src.a)
	dest.SetF(src.f)
	dest.SetB(src.b)
	dest.SetC(src.c)
	dest.SetD(src.d)
	dest.SetE(src.e)
	dest.SetH(src.h)
	dest.SetL(src.l)
	dest.SetSP(src.sp)
	dest.SetPC(src.pc)
}

func assertRegs(t *testing.T, expected *testRegsValue, regs *registers.Registers) {
	assert.Equal(t, expected.a, regs.GetA())
	assert.Equal(t, expected.f, regs.GetF())
	assert.Equal(t, expected.b, regs.GetB())
	assert.Equal(t, expected.c, regs.GetC())
	assert.Equal(t, expected.d, regs.GetD())
	assert.Equal(t, expected.e, regs.GetE())
	assert.Equal(t, expected.h, regs.GetH())
	assert.Equal(t, expected.l, regs.GetL())
	assert.Equal(t, expected.sp, regs.GetSP())
	assert.Equal(t, expected.pc, regs.GetPC())
}

type testData struct {
	name           string
	inMemory       []uint8
	inRegs         testRegsValue
	expectedMemory []uint8
	expectedRegs   testRegsValue
	cycles         uint8
}
type testRegsValue struct {
	a  uint8
	f  uint8
	b  uint8
	c  uint8
	d  uint8
	e  uint8
	h  uint8
	l  uint8
	sp uint16
	pc uint16
}

func TestCPU(t *testing.T) {
	testMMU := newTestMem(t)
	testCPU := NewCPU(testMMU, &testInterrupts{})

	// Clear memory & regs
	testMMU.Clear(t, 0x00)
	clearRegs(t, testCPU)

	testData := []testData{}
	testData = append(testData, testDataLDR8...)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			testMMU.Push(tt.inMemory)
			initTestRegs(&tt.inRegs, &testCPU.regs)
			cycles := testCPU.Tick()
			assert.Equal(t, tt.cycles, cycles)
			if tt.expectedMemory == nil {
				assert.Equal(t, tt.inMemory, testMMU.Pop())
			} else {
				assert.Equal(t, tt.expectedMemory, testMMU.Pop())
			}
			assertRegs(t, &tt.expectedRegs, &testCPU.regs)
		})
	}
}
