package audio

type SquareChannel struct {
	////// Square Generator Pattern //////
	pattern         [4][8]bool
	patternSelected uint8
	patternPosition uint8

	////// Square Generator Position //////
	timer   uint16
	counter uint16 // If counter >= Timer -> patternPosition++ % 8; counter=0

	////// Frame Sequencer //////
	sequencer        uint16
	sequencerCounter uint16

	////// Sample Counter //////
	soundLength        uint16
	soundLengthEnable  bool
	soundLengthCounter uint16

	Volume  uint8
	started bool
}

func (sc *SquareChannel) Initial() {
	// TODO : reinit counters
	sc.started = true
}

func (sc *SquareChannel) SetFrequency(freq uint16)    { sc.timer = (2048 - freq) * 4 }
func (sc *SquareChannel) SelectPattern(pattern uint8) { sc.patternSelected = pattern }

func (sc *SquareChannel) EnableSoundLength(enable bool) {
	sc.soundLengthEnable = enable
}
func (sc *SquareChannel) SetSoundLength(soundLength uint8) {
	sc.soundLength = (64 - uint16(soundLength)) * 1 / 256
}

func (sc *SquareChannel) Tick() uint8 {
	if !sc.started {
		return 0
	}
	var sampleValue uint8 = 0

	// SAMPLE VALUE
	if sc.pattern[sc.patternSelected][sc.patternPosition] {
		sampleValue = sc.Volume
	}

	// SOUND LENGTH
	if sc.soundLengthEnable {
		if sc.soundLengthCounter > sc.soundLength {
			sampleValue = 0
		} else {
			sc.soundLengthCounter++
		}
	}

	// INC
	sc.counter++
	if sc.counter >= sc.timer {
		sc.counter = 0
		sc.patternPosition = (sc.patternPosition + 1) % 8
	}
	return sampleValue
}

func NewSquareChannel() *SquareChannel {
	return &SquareChannel{
		pattern: [4][8]bool{
			0: {false, false, false, false, false, false, false, true},
			1: {true, false, false, false, false, false, false, true},
			2: {true, false, false, false, false, true, true, true},
			3: {false, true, true, true, true, true, true, false},
		},
		Volume: 255,
	}
}
