package scenes

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/battle"
	"github.com/quasilyte/gmtk2023/gamedata"
)

type BattleController struct {
	scene *ge.Scene

	config *gamedata.BattleConfig

	runner *battle.Runner
}

func NewBattleController(config *gamedata.BattleConfig) *BattleController {
	return &BattleController{
		config: config,
	}
}

func (c *BattleController) Init(scene *ge.Scene) {
	c.scene = scene

	c.runner = battle.NewRunner(c.config)
}

func (c *BattleController) Update(delta float64) {}
