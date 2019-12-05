package nullio

import (
	"github.com/jmontupet/gbcore/pkg/coreio"
)

type nullFrameDrawer struct{}

func (p *nullFrameDrawer) SwapFrameBuffer(
	frameBuffer *coreio.FrameBuffer, colors *coreio.FrameColors,
) (*coreio.FrameBuffer, *coreio.FrameColors) {
	return frameBuffer, colors
}

func NewNullFrameDrawer() coreio.FrameDrawer {
	return &nullFrameDrawer{}
}
