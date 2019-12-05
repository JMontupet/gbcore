package coreio

import (
	"github.com/jmontupet/gbcore/internal/pkg/constants"
)

const FrameBufferSize = constants.ScreenWidth * constants.ScreenHeight       // X * Y pixels
const FrameColorsSize = constants.ScreenColors * constants.ScreenColorsDepth // N Colors * B bytes

type FrameBuffer [FrameBufferSize]uint8
type FrameColors [FrameColorsSize]uint8
type AudioBuffer [constants.AudioBufferSamples]uint8
type KeyInputState uint8

type FrameDrawer interface {
	SwapFrameBuffer(frameBuffer *FrameBuffer, colors *FrameColors) (*FrameBuffer, *FrameColors)
}
type AudioPlayer interface {
	SwapAudioBuffer(frameBuffer []uint8) []uint8
}
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
