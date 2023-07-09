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

	objects          []ge.SceneObject
	addedObjects     []ge.SceneObject
	projectiles      []*projectile
	addedProjectiles []*projectile

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

	r.world = newWorldState(scene)
	r.world.runner = r
	r.world.Camera = r.camera
	r.world.PlayerInput = r.config.PlayerInput
	r.world.pathgrid = pathing.NewGrid(gamedata.CellSize*gamedata.NumSegmentCells, r.camera.Rect.Height())
	r.world.bfs = pathing.NewGreedyBFS(r.world.pathgrid.Size())

	p := newHumanPlayer(r.world, r.config.PlayerDesigns)
	r.players = append(r.players, p)
	p.Init()

	cp := newComputerPlayer(r.world)
	r.players = append(r.players, cp)
	cp.Init()

	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: (40 * 4) - 20, Y: (40 * 9) - 20},
	// 	Stats: gamedata.ConstructorUnitStats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: (40 * 5) - 20, Y: (40 * 9) - 20},
	// 	Stats: gamedata.ConstructorUnitStats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: (40 * 3) - 20, Y: (40 * 9) - 20},
	// 	Stats: gamedata.ConstructorUnitStats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: (40 * 6) - 20, Y: (40 * 9) - 20},
	// 	Stats: gamedata.ConstructorUnitStats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: (40 * 6) - 20, Y: (40 * 10) - 20},
	// 	Stats: gamedata.ConstructorUnitStats,
	// }))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 5) - 20, Y: (40 * 10) - 20},
		Stats: gamedata.ConstructorUnitStats,
	}))

	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 96, Y: 96},
		Stats: gamedata.CommanderUnitStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 160, Y: 160},
		Stats: gamedata.CommanderUnitStats,
	}))

	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 3) - 20, Y: (40 * 10) - 20},
		Stats: gamedata.MCVUnitStats,
	}))

	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: (40 * 6) - 20, Y: (40 * 4) - 20},
	// 	Stats: gamedata.TankFactoryUnitStats,
	// }))

	tankStats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.FighterBodyStats,
		Turret:   gamedata.HurricaneStats,
	}
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: 240, Y: 240},
		Stats: tankStats,
	}))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 300, Y: 300},
	// 	Stats: tankStats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 100, Y: 200},
	// 	Stats: tankStats,
	// }))
	// tank2Stats := &gamedata.UnitStats{
	// 	Movement: gamedata.UnitMovementGround,
	// 	Body:     gamedata.ScoutBodyStats,
	// 	Turret:   gamedata.ScatterCannonStats,
	// }
	// tank3Stats := &gamedata.UnitStats{
	// 	Movement: gamedata.UnitMovementGround,
	// 	Body:     gamedata.WheelsBodyStats,
	// 	Turret:   gamedata.GatlingStats,
	// }
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 360, Y: 360},
	// 	Stats: tank3Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 400, Y: 360},
	// 	Stats: tank3Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 440, Y: 360},
	// 	Stats: tank3Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 360, Y: 400},
	// 	Stats: tank3Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 400, Y: 400},
	// 	Stats: tank2Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 440, Y: 400},
	// 	Stats: tank2Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 400, Y: 440},
	// 	Stats: tank2Stats,
	// }))
	// r.AddObject(r.world.NewUnit(unitConfig{
	// 	Pos:   gmath.Vec{X: 440, Y: 480},
	// 	Stats: tank2Stats,
	// }))

	enemyTankStats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.WheelsBodyStats,
		Turret:   gamedata.GatlingStats,
		Creep:    true,
	}
	enemyTank2Stats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.WheelsBodyStats,
		Turret:   gamedata.ScatterCannonStats,
		Creep:    true,
	}
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 22) - 20, Y: (40 * 5) - 20},
		Stats: enemyTankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 22) - 20, Y: (40 * 6) - 20},
		Stats: enemyTankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 22) - 20, Y: (40 * 7) - 20},
		Stats: enemyTankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 23) - 20, Y: (40 * 5) - 20},
		Stats: enemyTankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 23) - 20, Y: (40 * 6) - 20},
		Stats: enemyTankStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 23) - 20, Y: (40 * 7) - 20},
		Stats: enemyTankStats,
	}))

	enemyBunkerStats := &gamedata.UnitStats{
		Movement: gamedata.UnitMovementNone,
		Body:     gamedata.CreepBunkerBodyStats,
		Turret:   gamedata.LightCannonStats,
		Creep:    true,
	}
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 28) - 20, Y: (40 * 4) - 20},
		Stats: enemyBunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 26) - 20, Y: (40 * 6) - 20},
		Stats: enemyBunkerStats,
	}))
	r.AddObject(r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 28) - 20, Y: (40 * 8) - 20},
		Stats: enemyBunkerStats,
	}))

	creepFactory1 := r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 28) - 20, Y: (40 * 10) - 20},
		Stats: gamedata.CreepTankFactoryUnitStats,
		Extra: &tankFactoryExtra{
			tankDesign: enemyTankStats,
		},
	})
	r.AddObject(creepFactory1)

	creepFactory2 := r.world.NewUnit(unitConfig{
		Pos:   gmath.Vec{X: (40 * 30) - 20, Y: (40 * 9) - 20},
		Stats: gamedata.CreepTankFactoryUnitStats,
		Extra: &tankFactoryExtra{
			tankDesign: enemyTank2Stats,
		},
	})
	r.AddObject(creepFactory2)
}

func (r *Runner) Update(delta float64) {
	scaledDelta := delta * r.gameSpeedMultiplier

	r.camera.Stage.Update()

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

	liveProjectiles := r.projectiles[:0]
	for _, p := range r.projectiles {
		if p.IsDisposed() {
			p.world.FreeProjectileNode(p)
			continue
		}
		p.Update(scaledDelta)
		liveProjectiles = append(liveProjectiles, p)
	}
	r.projectiles = liveProjectiles
	r.projectiles = append(r.projectiles, r.addedProjectiles...)
	r.addedProjectiles = r.addedProjectiles[:0]

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

func (r *Runner) AddProjectile(p *projectile) {
	r.addedProjectiles = append(r.addedProjectiles, p)
	p.Init(r.scene)
}

func (r *Runner) AddObject(o ge.SceneObject) {
	r.addedObjects = append(r.addedObjects, o)
	o.Init(r.scene)
}
