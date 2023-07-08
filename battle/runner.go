package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/pathing"
	"github.com/quasilyte/gmtk2023/session"
	"github.com/quasilyte/gmtk2023/viewport"
)

type Runner struct {
	scene *ge.Scene

	state *session.State

	world *worldState

	objects      []ge.SceneObject
	addedObjects []ge.SceneObject

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

		objects:      make([]ge.SceneObject, 0, 512),
		addedObjects: make([]ge.SceneObject, 0, 32),
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
	r.world.runner = r
	r.world.Camera = r.camera
	r.world.PlayerInput = r.config.PlayerInput
	r.world.pathgrid = pathing.NewGrid(gamedata.CellSize*gamedata.NumSegmentCells, r.camera.Rect.Height())
	r.world.bfs = pathing.NewGreedyBFS(r.world.pathgrid.Size())

	p := newHumanPlayer(r.world)
	r.players = append(r.players, p)
	p.Init()

	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 96, Y: 96},
		Stats: gamedata.CommanderUnitStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 160, Y: 160},
		Stats: gamedata.CommanderUnitStats,
	}))

	tankStats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.DestroyerBodyStats,
		Turret:   gamedata.LightCannonStats,
	}
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 240, Y: 240},
		Stats: tankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 300, Y: 300},
		Stats: tankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 100, Y: 200},
		Stats: tankStats,
	}))
	tank2Stats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.ScoutBodyStats,
		Turret:   gamedata.LightCannonStats,
	}
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 360, Y: 360},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 400, Y: 360},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 440, Y: 360},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 360, Y: 400},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 400, Y: 400},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 440, Y: 400},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 400, Y: 440},
		Stats: tank2Stats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 440, Y: 480},
		Stats: tank2Stats,
	}))

	bunkerStats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementNone,
		Body:     gamedata.BunkerBodyStats,
		Turret:   gamedata.LightCannonStats,
	}
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 8) - 20, Y: (40 * 5) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 9) - 20, Y: (40 * 5) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 8) - 20, Y: (40 * 4) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 9) - 20, Y: (40 * 4) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 10) - 20, Y: (40 * 5) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 11) - 20, Y: (40 * 5) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 10) - 20, Y: (40 * 4) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 11) - 20, Y: (40 * 4) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 8) - 20, Y: (40 * 7) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 9) - 20, Y: (40 * 7) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 10) - 20, Y: (40 * 7) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 11) - 20, Y: (40 * 7) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 8) - 20, Y: (40 * 9) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 9) - 20, Y: (40 * 9) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 10) - 20, Y: (40 * 9) - 20},
		Stats: bunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 11) - 20, Y: (40 * 9) - 20},
		Stats: bunkerStats,
	}))
}

func (r *Runner) Update(delta float64) {
	scaledDelta := delta * r.gameSpeedMultiplier

	for _, u := range r.world.playerUnits.selectable {
		if !u.IsCommander() {
			continue
		}
		u.group = u.group[:0]
	}
	for _, u := range r.world.playerUnits.nonSelectable {
		if u.leader != nil {
			if u.leader.IsDisposed() {
				u.leader = nil
			} else {
				u.leader.group = append(u.leader.group, u)
			}
		}
	}

	for _, p := range r.players {
		p.Update(scaledDelta, delta)
	}

	liveObjects := r.objects[:0]
	for _, o := range r.objects {
		if o.IsDisposed() {
			continue
		}
		o.Update(scaledDelta)
		liveObjects = append(liveObjects, o)
	}
	r.objects = liveObjects
	r.objects = append(r.objects, r.addedObjects...)
	r.addedObjects = r.addedObjects[:0]
}

func (r *Runner) AddObject(o ge.SceneObject) {
	r.addedObjects = append(r.addedObjects, o)
	o.Init(r.scene)
}
