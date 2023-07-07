package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/viewport"
)

type tankSelector struct {
	sprite *ge.Sprite

	hp *hpBar
}

func newTankSelector() *tankSelector {
	return &tankSelector{
		hp: newHPBar(),
	}
}

func (sel *tankSelector) Init(scene *ge.Scene, stage *viewport.Stage) {
	sel.hp.Init(scene, stage)

	sel.sprite = scene.NewSprite(assets.ImageUITankSelector)
	sel.sprite.Visible = false
	stage.AddGraphicsSlightlyAbove(sel.sprite)
}

func (sel *tankSelector) Update() {
	sel.hp.Update()
}

func (sel *tankSelector) IsActive() bool {
	return sel.hp.unit != nil
}

func (sel *tankSelector) GetUnit() *unit { return sel.hp.unit }

func (sel *tankSelector) SetUnit(u *unit) {
	sel.hp.SetTarget(u)

	isTank := u != nil && u.stats.Movement == gamedata.UnitMovementGround
	sel.sprite.Visible = isTank
	if isTank {
		sel.sprite.Pos.Base = &u.spritePos
	}
}
