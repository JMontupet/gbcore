package nullio

import (
	"github.com/jmontupet/gbcore/pkg/coreio"
)

type nullAudioPlayer struct{}

func (p *nullAudioPlayer) SwapAudioBuffer(data []uint8) []uint8 {
	return data
}

func NewNullAudioPlayer() coreio.AudioPlayer {
	return &nullAudioPlayer{}
}
