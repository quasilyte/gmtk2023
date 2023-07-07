package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gsignal"
)

type unit struct {
	stats *gamedata.UnitStats

	world *worldState

	pos       gmath.Vec
	spritePos gmath.Vec

	waypoint gmath.Vec

	sprite *ge.Sprite
	anim   *ge.Animation

	turret *turret

	leader *unit
	group  []*unit

	hp    float64
	maxHP float64

	disposed bool

	EventDestroyed gsignal.Event[*unit]
}

type unitConfig struct {
	Stats *gamedata.UnitStats
	Pos   gmath.Vec
}

func newUnit(world *worldState, config unitConfig) *unit {
	u := &unit{
		stats: config.Stats,
		world: world,
		pos:   config.Pos,
	}
	if config.Stats.Body != nil {
		u.maxHP = config.Stats.Body.HP + config.Stats.Turret.HP
	}
	u.maxHP += config.Stats.HP
	u.hp = u.maxHP
	u.hp *= world.Scene().Rand().FloatRange(0.2, 0.9)
	return u
}

func (u *unit) IsDisposed() bool {
	return u.disposed
}

func (u *unit) IsCommander() bool { return u.stats == gamedata.CommanderUnitStats }

func (u *unit) Dispose() {
	u.disposed = true
	u.sprite.Dispose()
}

func (u *unit) Init(scene *ge.Scene) {
	u.spritePos.X = math.Round(u.pos.X)
	u.spritePos.Y = math.Round(u.pos.Y)

	if u.stats.Image != assets.ImageNone {
		u.sprite = scene.NewSprite(u.stats.Image)
		u.sprite.Pos.Base = &u.spritePos
		if u.sprite.ImageWidth() != u.sprite.FrameWidth {
			u.anim = ge.NewRepeatedAnimation(u.sprite, -1)
			u.anim.SetAnimationSpan(0.5)
		}
		u.world.Stage().AddSpriteSlightlyAbove(u.sprite)
	} else {
		u.sprite = ge.NewSprite(scene.Context())
		u.sprite.SetImage(u.stats.Body.Texture)
		u.sprite.Pos.Base = &u.spritePos
		u.world.Stage().AddSprite(u.sprite)
		u.turret = newTurret(u.world, turretConfig{
			Image: u.stats.Turret.Texture,
			Pos:   &u.spritePos,
		})
		u.world.runner.AddObject(u.turret)
	}
}

func (u *unit) Update(delta float64) {
	if u.anim != nil {
		u.anim.Tick(delta)
	}

	u.spritePos.X = math.Round(u.pos.X)
	u.spritePos.Y = math.Round(u.pos.Y)

	if !u.waypoint.IsZero() {
		u.moveToWaypoint(delta)
	}
}

func (u *unit) SendTo(pos gmath.Vec) {
	u.waypoint = pos
}

func (u *unit) calcSpeed() float64 {
	return u.stats.Speed
}

func (u *unit) moveToWaypoint(delta float64) {
	travelled := u.calcSpeed() * delta
	switch u.stats.Movement {
	case gamedata.UnitMovementHover:
		u.pos = u.pos.MoveTowards(u.waypoint, travelled)
		if u.pos == u.waypoint {
			u.waypoint = gmath.Vec{}
		}
	}
}
