package coreio

import (
	"github.com/jmontupet/gbcore/internal/pkg/constants"
)

// FrameBufferSize is the size of a FrameBuffer
const FrameBufferSize = constants.ScreenWidth * constants.ScreenHeight // X * Y pixels
// FrameColorsSize is the size of a FrameColors array
const FrameColorsSize = constants.ScreenColors * constants.ScreenColorsDepth // N Colors * B bytes

// FrameBuffer contains one frame colors indexes
type FrameBuffer [FrameBufferSize]uint8

// FrameColors contains one frame color palette
type FrameColors [FrameColorsSize]uint8

// AudioBuffer contains n samples of audio data
type AudioBuffer [constants.AudioBufferSamples]uint8

// KeyInputState store inputs states on 8 bits
type KeyInputState uint8

// FrameDrawer receives frame data and draw it in is own way.
// Old buffers should be returned to prevent new memory allocation.
type FrameDrawer interface {
	SwapFrameBuffer(frameBuffer *FrameBuffer, colors *FrameColors) (*FrameBuffer, *FrameColors)
}

// AudioPlayer plays audio
// Old buffers should be returned to prevent new memory allocation.
type AudioPlayer interface {
	SwapAudioBuffer(frameBuffer []uint8) []uint8
}

// InputsManager return current inputs states
type InputsManager interface {
	CurrentInput() KeyInputState
}

const (
	GBKeyA      KeyInputState = 1 << iota
	GBKeyB      KeyInputState = 1 << iota
	GBKeySELECT KeyInputState = 1 << iota
	GBKeySTART  KeyInputState = 1 << iota
	GBKeyRIGHT  KeyInputState = 1 << iota
	GBKeyLEFT   KeyInputState = 1 << iota
	GBKeyUP     KeyInputState = 1 << iota
	GBKeyDOWN   KeyInputState = 1 << iota
)
