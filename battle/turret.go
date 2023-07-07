package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type turret struct {
	world *worldState

	sprite *ge.Sprite

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

func (t *turret) Update(delta float64) {}
