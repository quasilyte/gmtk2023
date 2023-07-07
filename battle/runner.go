package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/session"
	"github.com/quasilyte/gmtk2023/viewport"
)

type Runner struct {
	scene *ge.Scene

	state *session.State

	world *worldState

	gameSpeedMultiplier float64

	camera *viewport.Camera

	players []player

	config *gamedata.BattleConfig
}

func NewRunner(state *session.State, config *gamedata.BattleConfig, cam *viewport.Camera) *Runner {
	return &Runner{
		state:  state,
		config: config,
		camera: cam,
	}
}

func (r *Runner) IsDisposed() bool { return false }

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene

	r.gameSpeedMultiplier = [...]float64{
		0.50,
		1.00,
		1.50,
		2.00,
	}[r.config.GameSpeed]

	r.world = newWorldState()
	r.world.Camera = r.camera
	r.world.PlayerInput = r.config.PlayerInput

	r.players = append(r.players, newHumanPlayer(r.world))
}

func (r *Runner) Update(delta float64) {
	scaledDelta := delta * r.gameSpeedMultiplier

	for _, p := range r.players {
		p.Update(scaledDelta, delta)
	}
}
