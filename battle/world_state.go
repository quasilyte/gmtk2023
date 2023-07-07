package battle

import (
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmtk2023/viewport"
)

type worldState struct {
	Camera *viewport.Camera

	PlayerInput *input.Handler
}

func newWorldState() *worldState {
	return &worldState{}
}
