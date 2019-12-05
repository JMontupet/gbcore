package memory

type Memory interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}
