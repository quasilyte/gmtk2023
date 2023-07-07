package scenes

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
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

	bg := ge.NewTiledBackground(scene.Context())
	bg.LoadTileset(scene.Context(), 1920/2, 1080/2, assets.ImageBackgroundTiles, assets.RawBackgroundTileset)
	scene.AddGraphics(bg)
}

func (c *BattleController) Update(delta float64) {}
