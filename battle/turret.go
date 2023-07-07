package battle

import (
	"math"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/gamedata"
)

type turret struct {
	world *worldState

	sprite *ge.Sprite

	rotation gmath.Rad

	config turretConfig
}

type turretConfig struct {
	Image resource.Image

	Pos *gmath.Vec
}

func newTurret(world *worldState, config turretConfig) *turret {
	return &turret{
		world:  world,
		config: config,
	}
}

func (t *turret) Init(scene *ge.Scene) {
	t.rotation = scene.Rand().Rad()

	t.sprite = ge.NewSprite(scene.Context())
	t.sprite.SetImage(t.config.Image)
	t.sprite.Pos.Base = t.config.Pos
	t.world.Stage().AddSprite(t.sprite)
}

func (t *turret) IsDisposed() bool {
	return t.sprite.IsDisposed()
}

func (t *turret) Dispose() {
	t.sprite.Dispose()
}

func (t *turret) Update(delta float64) {
	spriteAngle := t.rotation.Normalized() - (gamedata.TankFrameAngleStep / 2)
	t.sprite.FrameOffset.X = 48 * math.Trunc(float64(spriteAngle/gamedata.TankFrameAngleStep))
}
