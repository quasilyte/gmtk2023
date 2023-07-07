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

	sprite *ge.Sprite
	anim   *ge.Animation

	disposed bool

	EventDestroyed gsignal.Event[*unit]
}

type unitConfig struct {
	Stats *gamedata.UnitStats
	Pos   gmath.Vec
}

func newUnit(world *worldState, config unitConfig) *unit {
	return &unit{
		stats: config.Stats,
		world: world,
		pos:   config.Pos,
	}
}

func (u *unit) IsDisposed() bool {
	return u.disposed
}

func (u *unit) Dispose() {
	u.disposed = true
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
		u.world.Stage().AddSprite(u.sprite)
	} else {
		// TODO: a tank texture?
	}
}

func (u *unit) Update(delta float64) {
	if u.anim != nil {
		u.anim.Tick(delta)
	}

	u.spritePos.X = math.Round(u.pos.X)
	u.spritePos.Y = math.Round(u.pos.Y)
}
