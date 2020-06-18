package emulator

import (
	"log"

	"github.com/jmontupet/gbcore/pkg/nullio"

	"github.com/jmontupet/gbcore/internal/pkg/cartridge"
	"github.com/jmontupet/gbcore/internal/pkg/gameboy"
	"github.com/jmontupet/gbcore/pkg/coreio"
)

type Emulator interface {
	gameboy.GameBoy
	GetGameTitle() string
}

type gbcEmulator struct {
	gbc      gameboy.GameBoy
	cartidge cartridge.Cartridge
}

func (e *gbcEmulator) Start() { e.gbc.Start() }
func (e *gbcEmulator) Stop()  { e.gbc.Stop() }
func (e *gbcEmulator) Pause() { e.gbc.Pause() }
func (e *gbcEmulator) GetGameTitle() string {
	return cartridge.ReadTitle(e.cartidge)
}

func NewEmulator(
	gameData []byte,
	renderer coreio.FrameDrawer,
	inputsManager coreio.InputsManager,
	audioPlayer coreio.AudioPlayer,
) (Emulator, error) {
	cart, err := cartridge.NewCartridge(gameData)
	if err != nil {
		return nil, err
	}
	if audioPlayer == nil {
		log.Println("No AudioPlayer. Null Audio player used.")
		audioPlayer = nullio.NewNullAudioPlayer()
	}
	if renderer == nil {
		log.Println("No Renderer. Null Frame Drawer used.")
		renderer = nullio.NewNullFrameDrawer()
	}
	if inputsManager == nil {
		log.Println("No Inputs Manager. Null Inputs Manager used.")
		inputsManager = nullio.NewNullInputsManager()
	}

	return &gbcEmulator{
		cartidge: cart,
		gbc: gameboy.NewGameBoy(
			cart,
			renderer,
			inputsManager,
			audioPlayer,
		),
	}, nil
}
