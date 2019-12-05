package gpu

type tile [8 * 2]uint8

// appendPixelsLine Appends to buff a pixels line of the sprite.
//
// N (offset) first pixels can be skiped (max 7)
// N (nb) pixels only can be append (max 8)
// A colorPrefix can be added to each pixels (cf. palette manager)
func (t *tile) appendPixelsLine(buff []uint8, line uint8, offset uint8, nb uint8, colorPrefix uint8, hFlip bool, vFlip bool) []uint8 {
	if vFlip {
		line = uint8((int16(line) - 7) * -1)
	}
	loByte, hiByte := t[line*2], t[line*2+1]

	// if over := (len(buff) + int(nb) - int(offset)) - cap(buff); over > 0 {
	// 	println("OVER DRAW", over)
	// }

	// A loop is little bit slower than this :
	if hFlip {
		toAppend := [8]uint8{
			colorPrefix + (hiByte<<1)&2 + (loByte)&1,
			colorPrefix + (hiByte)&2 + (loByte>>1)&1,
			colorPrefix + (hiByte>>1)&2 + (loByte>>2)&1,
			colorPrefix + (hiByte>>2)&2 + (loByte>>3)&1,
			colorPrefix + (hiByte>>3)&2 + (loByte>>4)&1,
			colorPrefix + (hiByte>>4)&2 + (loByte>>5)&1,
			colorPrefix + (hiByte>>5)&2 + (loByte>>6)&1,
			colorPrefix + (hiByte>>6)&2 + loByte>>7,
		}
		return append(buff, toAppend[offset:nb]...)
	}
	toAppend := [8]uint8{
		colorPrefix + (hiByte>>6)&2 + loByte>>7,
		colorPrefix + (hiByte>>5)&2 + (loByte>>6)&1,
		colorPrefix + (hiByte>>4)&2 + (loByte>>5)&1,
		colorPrefix + (hiByte>>3)&2 + (loByte>>4)&1,
		colorPrefix + (hiByte>>2)&2 + (loByte>>3)&1,
		colorPrefix + (hiByte>>1)&2 + (loByte>>2)&1,
		colorPrefix + (hiByte)&2 + (loByte>>1)&1,
		colorPrefix + (hiByte<<1)&2 + (loByte)&1,
	}
	return append(buff, toAppend[offset:nb]...)
}
