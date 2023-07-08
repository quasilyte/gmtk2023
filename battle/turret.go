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

	rotation    gmath.Rad
	dstRotation gmath.Rad
	needRotate  bool

	config turretConfig
}

type turretConfig struct {
	Image resource.Image

	Pos *gmath.Vec

	Stats *gamedata.UnitStats
}

func newTurret(world *worldState, config turretConfig) *turret {
	return &turret{
		world:  world,
		config: config,
	}
}

func (t *turret) Init(scene *ge.Scene) {
	t.sprite = ge.NewSprite(scene.Context())
	t.sprite.SetImage(t.config.Image)
	t.sprite.Pos.Base = t.config.Pos
	t.world.Stage().AddSprite(t.sprite)

	t.setRotation(scene.Rand().Rad())
}

func (t *turret) IsDisposed() bool {
	return t.sprite.IsDisposed()
}

func (t *turret) Dispose() {
	t.sprite.Dispose()
}

func (t *turret) Update(delta float64) {
	if t.needRotate {
		t.setRotation(t.rotation.RotatedTowards(t.dstRotation, t.config.Stats.Turret.RotationSpeed*gmath.Rad(delta)))
		if t.rotation == t.dstRotation {
			t.needRotate = false
		}
	}
}

func (t *turret) setRotation(v gmath.Rad) {
	t.rotation = v
	spriteAngle := t.rotation.Normalized() - (gamedata.TankFrameAngleStep / 2)
	t.sprite.FrameOffset.X = 48 * math.Trunc(float64(spriteAngle/gamedata.TankFrameAngleStep))
}

func (t *turret) AlignRequest(rotation gmath.Rad) {
	if rotation == t.rotation {
		return
	}

	t.needRotate = true
	t.dstRotation = rotation
}
