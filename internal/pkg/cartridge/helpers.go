package cartridge

import "fmt"

func splitROMBanks(romSizeFlag uint8, rom []uint8) ([][]uint8, error) {
	var nbBanks uint

	switch romSizeFlag {
	case 0x00: //  32KByte (no ROM banking)
		nbBanks = 2
	case 0x01: //  64KByte (4 banks)
		nbBanks = 4
	case 0x02: // 128KByte (8 banks)
		nbBanks = 8
	case 0x03: // 256KByte (16 banks)
		nbBanks = 16
	case 0x04: // 512KByte (32 banks)
		nbBanks = 32
	case 0x05: //   1MByte (64 banks)  - only 63 banks used by MBC1
		nbBanks = 64
	case 0x06: //   2MByte (128 banks) - only 125 banks used by MBC1
		nbBanks = 128
	case 0x07: //   4MByte (256 banks)
		nbBanks = 256
	case 0x08: //   8MByte (512 banks)
		nbBanks = 512
	case 0x52: // 1.1MByte (72 banks)
		nbBanks = 72
	case 0x53: // 1.2MByte (80 banks)
		nbBanks = 80
	case 0x54: // 1.5MByte (96 banks)
		nbBanks = 96
	}

	banks := make([][]uint8, nbBanks)
	for i := range banks {
		banks[i] = rom[i*romBankSize : (i+1)*romBankSize]
	}
	return banks, nil
}

func makeRAMBanks(ramSizeFlag uint8) ([][]uint8, error) {
	switch ramSizeFlag {
	case 0x00: // 00h - None
		return [][]uint8{}, nil
	case 0x01: // 01h - 2 KBytes
		return [][]uint8{make([]uint8, 2*1024)}, nil
	case 0x02: // 02h - 8 Kbytes
		return [][]uint8{make([]uint8, 8*1024)}, nil
	case 0x03: // 03h - 32 KBytes (4 banks of 8KBytes each)
		banks := make([][]uint8, 4)
		for i := range banks {
			banks[i] = make([]uint8, 8*1024)
		}
		return banks, nil
	case 0x04: // 04h - 128 KBytes (16 banks of 8KBytes each)
		banks := make([][]uint8, 16)
		for i := range banks {
			banks[i] = make([]uint8, 8*1024)
		}
		return banks, nil
	case 0x05: // 05h - 64 KBytes (8 banks of 8KBytes each)
		banks := make([][]uint8, 8)
		for i := range banks {
			banks[i] = make([]uint8, 8*1024)
		}
		return banks, nil
	default:
		return nil, fmt.Errorf("INVALID MEMORY SIZE : 0x%02X", ramSizeFlag)
	}
}
