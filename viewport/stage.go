package viewport

import (
	"github.com/quasilyte/ge"
)

type Stage struct {
	bg *ge.TiledBackground
}

func NewStage() *Stage {
	return &Stage{}
}

func (stage *Stage) SetBackground(bg *ge.TiledBackground) {
	stage.bg = bg
}

func (stage *Stage) Update() {}
