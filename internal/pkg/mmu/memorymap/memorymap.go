package memorymap

const (
	// BootROM (disabled after boot) - Removed from this emulator
	BootRomStart uint16 = 0x0
	BootRomEnd   uint16 = 0xFF

	// Cartridge ROM (16KB)
	FixedRomStart uint16 = 0x0
	FixedRomEnd   uint16 = 0x3FFF

	// Cartridge Switchable ROM (16KB)
	SwitchableRomStart uint16 = 0x4000
	SwitchableRomEnd   uint16 = 0x7FFF

	// Internal VRAM (8KB)
	VRamStart uint16 = 0x8000
	VRamEnd   uint16 = 0x9FFF

	// Cartridge RAM (switchable) (8KB)
	ExternalRAMStart uint16 = 0xA000
	ExternalRAMEnd   uint16 = 0xBFFF

	// Internal WorkRAM bank 0 (4KB)
	WRam0Start uint16 = 0xC000
	WRam0End   uint16 = 0xCFFF

	// Internal WorkRAM bank 1~7 (4KB) (7 only for CGB)
	WRam1Start uint16 = 0xD000
	WRam1End   uint16 = 0xDFFF

	// Mirror of C000~DDFF
	MirrorRAMStart uint16 = 0xE000
	MirrorRAMEnd   uint16 = 0xFDFF

	// OAM
	OAMStart uint16 = 0xFE00
	OAMEnd   uint16 = 0xFE9F

	// Not Usable (but used !)
	NotUsableStart uint16 = 0xFEA0
	NotUsableEnd   uint16 = 0xFEFF

	// I/O Registers
	IORegistersStart uint16 = 0xFF00
	IORegistersEnd   uint16 = 0xFF7F

	// HRAM
	HRamStart uint16 = 0xFF80
	HRamEnd   uint16 = 0xFFFE

	// Interrupts
	Interrupts uint16 = 0xFFFF
)
