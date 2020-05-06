package gpu

import (
	"reflect"
	"testing"
)

func TestTile(t *testing.T) {
	tl := tile{
		0b00000000,
		0b11111111,
	}

	var buff [8]uint8
	pixels := tl.appendPixelsLine(buff[:0], 0, 0, 8, 0b100000, false, false)

	if l := len(pixels); l != 8 {
		t.Fatalf("Pixels number : expected: %d, got: %d", 8, l)
	}
	if c := cap(pixels); c != 8 {
		t.Fatalf("Buff capacity : expected: %d, got: %d", 8, c)
	}
	if !reflect.DeepEqual(buff[:8], pixels) {
		t.Fatal("New buff allocated")
	}
}

func BenchmarkTile(b *testing.B) {
	tl := tile{
		0b00000000,
		0b11111111,
	}
	var buff [8]uint8
	for i := 0; i < b.N; i++ {
		tl.appendPixelsLine(buff[:0], 0, 0, 8, 0b100000, false, false)
	}
}
