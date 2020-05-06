package nullio

import (
	"github.com/jmontupet/gbcore/pkg/coreio"
)

type nullInputsManager struct{}

func (*nullInputsManager) CurrentInput() coreio.KeyInputState {
	return 0x00
}

func NewNullInputsManager() coreio.InputsManager {
	return &nullInputsManager{}
}
