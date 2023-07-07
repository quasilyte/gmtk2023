package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/viewport"
)

type hpBar struct {
	bgRect    *ge.Rect
	valueRect *ge.Rect

	cached float64
	unit   *unit
}

func newHPBar() *hpBar {
	return &hpBar{}
}

func (b *hpBar) Init(scene *ge.Scene, stage *viewport.Stage) {
	b.bgRect = ge.NewRect(scene.Context(), 32, 4)
	b.bgRect.FillColorScale.SetColor(ge.RGB(0x1e3a27))
	b.bgRect.Pos.Offset.Y = 24
	stage.AddGraphicsSlightlyAbove(b.bgRect)
	b.bgRect.Visible = false

	b.valueRect = ge.NewRect(scene.Context(), 0, 2)
	b.valueRect.FillColorScale.SetColor(ge.RGB(0x26cd61))
	b.valueRect.Pos.Offset.Y = 24
	stage.AddGraphicsSlightlyAbove(b.valueRect)
	b.valueRect.Visible = false
}

func (b *hpBar) SetTarget(u *unit) {
	b.unit = u
	b.cached = -1
	b.bgRect.Visible = u != nil
	b.valueRect.Visible = u != nil
	if u != nil {
		b.bgRect.Pos.Base = &u.spritePos
		b.valueRect.Pos.Base = &u.spritePos
	}
}

func (b *hpBar) Update() {
	if b.unit != nil && b.unit.IsDisposed() {
		b.SetTarget(nil)
		return
	}
	if b.unit == nil {
		return
	}
	if b.unit.hp == b.cached {
		return
	}

	b.cached = b.unit.hp
	percentage := b.unit.hp / b.unit.maxHP
	valueRectMaxWidth := b.bgRect.Width - 2
	b.valueRect.Width = math.Ceil(valueRectMaxWidth * percentage)
}
