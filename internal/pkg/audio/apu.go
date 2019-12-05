package audio

import (
	"github.com/jmontupet/gbcore/internal/pkg/constants"
	"github.com/jmontupet/gbcore/internal/pkg/ioports"
	"github.com/jmontupet/gbcore/pkg/coreio"
)

type APU struct {
	channel1 *SquareChannel
	channel2 *SquareChannel

	audioPlayer coreio.AudioPlayer

	nr11 *ioports.Ptr // FF11 - NR11 - Channel 1 Sound length/Wave pattern duty (R/W)
	nr13 *ioports.Ptr // FF13 - NR13 - Channel 1 Frequency lo (Write Only)
	nr14 *ioports.Ptr // FF14 - NR14 - Channel 1 Frequency hi (R/W)

	nr21 *ioports.Ptr // FF16 - NR21 - Channel 2 Sound Length/Wave Pattern Duty (R/W)
	nr23 *ioports.Ptr // FF18 - NR23 - Channel 2 Frequency lo data (W)
	nr24 *ioports.Ptr // FF19 - NR24 - Channel 2 Frequency hi data (R/W)

	nr50 *ioports.Ptr // FF24 - NR50 - Channel control / ON-OFF / Volume
	nr51 *ioports.Ptr // FF25 - NR51 - Selection of Sound output terminal
	nr52 *ioports.Ptr // FF26 - NR52 - Sound on/off
}

const resampleFactor = 70224 * constants.ScreenRefreshRate / constants.AudioFrequency
const audioSamplePerFrame = constants.AudioBufferSamples

var currentResamplerPosition = 0.0
var buffer = make([]uint8, audioSamplePerFrame*2)
var bufferPosition = 0

func (apu *APU) Tick(cycles uint8) {
	// Sound ON/OFF
	if apu.nr52.Get()>>7 == 0 {
		return
	}

	// Setup Channel 1
	apu.channel1.SelectPattern(apu.nr11.Get() >> 6)
	apu.channel1.SetFrequency(uint16(apu.nr13.Get()) | uint16(apu.nr14.Get()&0x07)<<8)
	apu.channel1.EnableSoundLength(apu.nr14.Get()&0x20 != 0)
	apu.channel1.SetSoundLength(apu.nr11.Get() & 0x3F)
	if apu.nr14.GetBit7() {
		apu.nr14.SetBit7(false)
		apu.channel1.Initial()
	}

	// Setup Channel 2
	apu.channel2.SelectPattern(apu.nr21.Get() >> 6)
	apu.channel2.SetFrequency(uint16(apu.nr23.Get()) | uint16(apu.nr24.Get()&0x07)<<8)
	apu.channel2.EnableSoundLength(apu.nr24.Get()&0x20 != 0)
	apu.channel2.SetSoundLength(apu.nr21.Get() & 0x3F)
	if apu.nr24.GetBit7() {
		apu.nr24.SetBit7(false)
		apu.channel2.Initial()
	}

	for i := uint8(0); i < cycles; i++ {
		chan1Sample := apu.channel1.Tick()
		chan2Sample := apu.channel2.Tick()
		// mixed := uint8((uint32(chan1Sample) + uint32(chan2Sample)) >> 2)

		currentResamplerPosition++
		if currentResamplerPosition >= resampleFactor {
			vS01 := uint32(apu.nr50.Get() & 0x07)
			vS02 := uint32(apu.nr50.Get() >> 4 & 0x07)
			vS01 = vS01 / 7
			vS02 = vS02 / 7

			out01 := make([]uint8, 4)
			out02 := make([]uint8, 4)
			// NR51 :
			// Bit 7 - Output sound 4 to SO2 terminal
			// Bit 6 - Output sound 3 to SO2 terminal
			// Bit 5 - Output sound 2 to SO2 terminal
			// Bit 4 - Output sound 1 to SO2 terminal
			// Bit 3 - Output sound 4 to SO1 terminal
			// Bit 2 - Output sound 3 to SO1 terminal
			// Bit 1 - Output sound 2 to SO1 terminal
			// Bit 0 - Output sound 1 to SO1 terminal
			if apu.nr51.GetBit0() {
				out01 = append(out01, chan1Sample)
			}
			if apu.nr51.GetBit1() {
				out01 = append(out01, chan2Sample)
			}
			if apu.nr51.GetBit4() {
				out02 = append(out02, chan1Sample)
			}
			if apu.nr51.GetBit5() {
				out02 = append(out02, chan2Sample)
			}

			if len(out01) > 0 {
				total := uint32(0)
				for _, v := range out01 {
					total += uint32(v)
				}
				buffer = append(buffer, uint8((total>>2)*vS01))
			} else {
				buffer = append(buffer, 0)
			}
			if len(out02) > 0 {
				total := uint32(0)
				for _, v := range out02 {
					total += uint32(v)
				}
				buffer = append(buffer, uint8((total>>2)*vS02))
			} else {
				buffer = append(buffer, 0)
			}

			// buffer = append(buffer, chan1Sample, chan2Sample)
			bufferPosition++
			if bufferPosition >= audioSamplePerFrame {
				bufferPosition = 0
				// SWAP
				buffer = apu.audioPlayer.SwapAudioBuffer(buffer)
			}
			currentResamplerPosition = 0.0
		}
	}
}

func NewAPU(io *ioports.IOPorts, audioPlayer coreio.AudioPlayer) *APU {
	return &APU{
		audioPlayer: audioPlayer,

		channel1: NewSquareChannel(),
		nr11:     io.NewPtr(0xFF11),
		nr13:     io.NewPtr(0xFF13),
		nr14:     io.NewPtr(0xFF14),

		channel2: NewSquareChannel(),
		nr21:     io.NewPtr(0xFF16),
		nr23:     io.NewPtr(0xFF18),
		nr24:     io.NewPtr(0xFF19),

		nr50: io.NewPtr(0xFF24),
		nr51: io.NewPtr(0xFF25),
		nr52: io.NewPtr(0xFF26),
	}
}

// 8-bit PCM 255 (FFh) 0
// 16-bit PCM 32767 (7FFFh) -32768 (-8000h)

// Duty   Waveform    Ratio
// -------------------------
// 0      00000001    12.5%
// 1      10000001    25%
// 2      10000111    50%
// 3      01111110    75%
