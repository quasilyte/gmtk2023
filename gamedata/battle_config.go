package gamedata

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/input"
)

type BattleConfig struct {
	PlayerInput *input.Handler

	PlayerDesigns *PlayerDesigns

	GameSpeed int
}

type PlayerDesigns struct {
	Towers []*UnitStats
	Tanks  []*UnitStats

	// [0] generator
	// [1] turret 1
	// [2] turret 2
	// [3] repair depot
	// [4] factory 1
	// [5] factory 2
	// [6] factory 3
	// [7] factory 4
	Icons []*ebiten.Image

	IconConstructor *ebiten.Image
	IconCommander   *ebiten.Image
}

func NewPlayerDesigns() *PlayerDesigns {
	icons := make([]*ebiten.Image, 8)
	for i := range icons {
		icons[i] = ebiten.NewImage(64, 72)
	}
	return &PlayerDesigns{
		Towers: make([]*UnitStats, 2),
		Tanks:  make([]*UnitStats, 4),
		Icons:  icons,
	}
}
