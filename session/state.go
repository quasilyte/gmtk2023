package session

import (
	"github.com/quasilyte/gmtk2023/eui"
)

type State struct {
	UIResources *eui.Resources
}

func NewState() *State {
	return &State{}
}
