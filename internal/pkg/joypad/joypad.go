package joypad

import (
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
)

// GO-GB MAP
// 0b00000001 A
// 0b00000010 B
// 0b00000100 Select
// 0b00001000 Start
// 0b00010000 Right
// 0b00100000 Left
// 0b01000000 Up
// 0b10000000 Down

// 0xFF00
// Bit 7 - Not used
// Bit 6 - Not used
// Bit 5 - P15 Select Button Keys      (0=Select)
// Bit 4 - P14 Select Direction Keys   (0=Select)
// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)

const buttonMask uint8 = 0b00100000
const arrowMask uint8 = 0b00010000
const selectMask uint8 = buttonMask | arrowMask

type Joypad struct {
	_memory uint8

	joypadInterrupt *ioports.BitPtr

	hwButton uint8 // /!\ bit down = key pressed. Complement of the received value from chan
	hwArrow  uint8 // /!\ bit down = key pressed. Complement of the received value from chan
}

func (j *Joypad) UpdateInput(input uint8) {
	hwButton := ^(input & 0x0F) // Complement
	hwArrow := ^(input >> 4)    // Complement

	// if any bit go dow
	if j.hwButton&hwButton != j.hwButton || j.hwArrow&hwArrow != j.hwArrow {
		j.joypadInterrupt.Set(true)
	}

	j.hwButton = hwButton
	j.hwArrow = hwArrow
}

func (j *Joypad) Read(_ uint16) uint8 { // addr = 0xFF00
	result := j._memory&selectMask | 0x0F
	if j._memory&buttonMask == 0 { // Button selected
		result &= j.hwButton
	}
	if j._memory&arrowMask == 0 { // Arrow selected
		result &= j.hwArrow
	}
	return result
}
func (j *Joypad) Write(_ uint16, value uint8) { // addr = 0xFF00
	// only bits 4-5 writable
	j._memory = value & selectMask
}

func NewJoypad(io *ioports.IOPorts) *Joypad {
	return &Joypad{
		_memory:  0xFF,
		hwButton: 0xFF,
		hwArrow:  0xFF,

		joypadInterrupt: io.NewBit4Ptr(0xFF0F),
	}
}
