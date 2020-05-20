package cartridge

import (
	"fmt"
	"log"
	"time"
)

type mbc3 struct {
	romBanks [][]uint8
	romBank  int

	ramBanks [][]uint8
	ramBank  int

	ramTimerEnable bool

	rtcEnable   bool
	rtcRegister uint8

	rtcTimer timer
}

func (c *mbc3) Read(addr uint16) uint8 {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF: // ROM CART FIXED
		return c.romBanks[0][addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // ROM CART BANK N
		return c.romBanks[c.romBank][addr-0x4000]
	case addr >= 0xA000 && addr <= 0xBFFF: // CART RAM
		if c.ramTimerEnable {
			if c.rtcEnable {
				if !c.rtcTimer.latched {
					fmt.Println("READ UNLATCHED TIMER")
				}
				dh, dl, h, m, s := c.rtcTimer.get()
				switch c.rtcRegister {
				case 0x08:
					return s
				case 0x09:
					return m
				case 0x0A:
					return h
				case 0x0B:
					return dl
				case 0x0C:
					return dh
				default:
					return 0x00
				}
			}
			return c.ramBanks[c.ramBank][addr-0xA000]
		}
		return 0x00
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
		return 0x00
	}
}

func (c *mbc3) Write(addr uint16, value uint8) {
	switch {
	// BANK CONTROLLER

	// 0000-1FFF - RAM and Timer Enable (Write Only)
	case addr >= 0x0000 && addr <= 0x1FFF:
		c.ramTimerEnable = value&0xF == 0x0A

	// 2000-3FFF - ROM Bank Number (Write Only)
	case addr >= 0x2000 && addr <= 0x3FFF:
		c.changeROMBank(value)

	// 4000-5FFF - RAM Bank Number - or - RTC Register Select (Write Only)
	case addr >= 0x4000 && addr <= 0x5FFF:
		if value < 0x08 { // 8 banks max
			c.ramBank = int(value)
			c.rtcEnable = false
		} else if value < 0x0D {
			c.rtcRegister = value
			c.rtcEnable = true
		}

	// 6000-7FFF - Latch Clock Data (Write Only)
	case addr >= 0x6000 && addr <= 0x7FFF:
		if value == 0x00 && c.rtcTimer.latched {
			c.rtcTimer.unlatch()
		} else if value == 0x01 && !c.rtcTimer.latched {
			c.rtcTimer.latch()
		}

	// CART RAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		if c.ramTimerEnable {
			if c.rtcEnable {
				dh, dl, h, m, s := c.rtcTimer.get()
				switch c.rtcRegister {
				case 0x08:
					c.rtcTimer.set(dh, dl, h, m, value)
				case 0x09:
					c.rtcTimer.set(dh, dl, h, value, s)
				case 0x0A:
					c.rtcTimer.set(dh, dl, value, m, s)
				case 0x0B:
					c.rtcTimer.set(dh, value, h, m, s)
				case 0x0C:
					c.rtcTimer.set(value, dl, h, m, s)
				default:
					return
				}
			}
			c.ramBanks[0][addr-0xA000] = value
		}
	// OFF RANGE
	default:
		log.Fatalf("MEMORY UNREACHABLE : 0x%04X", addr)
	}
}

func (c *mbc3) changeROMBank(v uint8) {
	value := int(v & 0x7F)
	if value > len(c.romBanks) {
		panic(fmt.Sprintf("ROM BANK 0x%02X DOES NOT EXIST. MAX : 0x%02X", value, len(c.romBanks)))
	}
	if value == 0 {
		c.romBank = 1
		return
	}
	c.romBank = value
}

func newMBC3(data []byte) (Cartridge, error) {
	cartridge := &mbc3{
		romBank:  1,
		rtcTimer: timer{start: time.Now().AddDate(1, 1, 1)},
	}

	romBanks, err := splitROMBanks(data[0x148], data)
	if err != nil {
		return nil, err
	}
	cartridge.romBanks = romBanks

	ramBanks, err := makeRAMBanks(ReadRAMSize(cartridge))
	if err != nil {
		return nil, err
	}
	cartridge.ramBanks = ramBanks

	return cartridge, nil
}

type timer struct {
	start time.Time
	stop  bool

	latched bool
	dh      uint8
	dl      uint8
	h       uint8
	m       uint8
	s       uint8
}

func (t *timer) set(dh uint8, dl uint8, h uint8, m uint8, s uint8) {
	var diff int64
	diff += int64(s)
	diff += int64(m * 60)
	diff += int64(h * 60 * 60)
	days := uint16(dl) & (uint16(dh&0x01) << 8)
	diff += int64(days * 60 * 60 * 24)

	t.start = time.Now().Add(-time.Duration(diff) * time.Second)
}

func (t *timer) latch() {
	t.latched = false
	t.dh, t.dl, t.h, t.m, t.s = t.get()
	t.latched = true
}

func (t *timer) unlatch() {
	t.latched = false
}

func (t *timer) get() (dh uint8, dl uint8, h uint8, m uint8, s uint8) {
	if t.latched {
		return t.dh, t.dl, t.h, t.m, t.s
	}

	now := time.Now()
	diff := int64(now.Sub(t.start).Seconds())
	s = uint8(diff % 60)
	diff /= 60
	m = uint8(diff % 60)
	diff /= 60
	h = uint8(diff % 24)
	diff /= 24
	days := uint16(diff)
	dl = uint8(days & 0x00ff)
	if days > 0xff {
		dh |= 1 << 0
	}
	if days > 0x1ff {
		dh |= 1 << 7
	}
	if t.stop {
		dh |= 1 << 6
	}
	return dh, dl, h, m, s
}
