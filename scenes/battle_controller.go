package scenes

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/battle"
	"github.com/quasilyte/gmtk2023/controls"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/session"
	"github.com/quasilyte/gmtk2023/viewport"
)

type BattleController struct {
	scene *ge.Scene

	state *session.State

	config *gamedata.BattleConfig
	camera *viewport.Camera
	stage  *viewport.Stage

	runner *battle.Runner
}

func NewBattleController(state *session.State, config *gamedata.BattleConfig) *BattleController {
	return &BattleController{
		config: config,
		state:  state,
	}
}

func (c *BattleController) Init(scene *ge.Scene) {
	c.scene = scene

	worldRect := gmath.Rect{
		Max: gmath.Vec{
			X: gamedata.CellSize * gamedata.NumSegmentCells,
			Y: 1080.0 / 2,
		},
	}

	scene.Audio().PlayMusic(assets.AudioMusicTrack1)

	bg := ge.NewTiledBackground(scene.Context())
	bg.LoadTileset(scene.Context(), worldRect.Width(), 1080.0/2, assets.ImageBackgroundTiles, assets.RawBackgroundTileset)

	c.stage = viewport.NewStage()
	c.stage.SetBackground(bg)

	c.camera = viewport.NewCamera(c.stage, worldRect, 1920.0/2, 1080.0/2)
	scene.AddGraphics(c.camera)

	c.runner = battle.NewRunner(c.state, c.config, c.camera)
	scene.AddObject(c.runner)
}

func (c *BattleController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *BattleController) back() {
	c.scene.Audio().PauseCurrentMusic()
	c.scene.Context().ChangeScene(NewPlayController(c.state))
}
