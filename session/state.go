package session

import (
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmtk2023/eui"
)

type State struct {
	UIResources *eui.Resources
	Input       *input.Handler
	Initialized bool
}

func NewState() *State {
	return &State{}
}
