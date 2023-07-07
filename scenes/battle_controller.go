package scenes

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/battle"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/viewport"
)

type BattleController struct {
	scene *ge.Scene

	config *gamedata.BattleConfig
	camera *viewport.Camera
	stage  *viewport.Stage

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

	worldRect := gmath.Rect{
		Max: gmath.Vec{
			X: 1920,
			Y: 1080.0 / 2,
		},
	}

	bg := ge.NewTiledBackground(scene.Context())
	bg.LoadTileset(scene.Context(), worldRect.Width(), 1080.0/2, assets.ImageBackgroundTiles, assets.RawBackgroundTileset)

	c.stage = viewport.NewStage()
	c.stage.SetBackground(bg)

	c.camera = viewport.NewCamera(c.stage, worldRect, 1920.0/2, 1080.0/2)
	scene.AddGraphics(c.camera)
}

func (c *BattleController) Update(delta float64) {
	c.camera.Offset.X += delta * 32
}
