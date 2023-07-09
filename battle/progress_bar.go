package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/styles"
	"github.com/quasilyte/gmtk2023/viewport"
)

type progressBar struct {
	bgRect    *ge.Rect
	valueRect *ge.Rect

	cached float64
}

func newProgressBar() *progressBar {
	return &progressBar{}
}

func (b *progressBar) Init(scene *ge.Scene, stage *viewport.Stage) {
	b.bgRect = ge.NewRect(scene.Context(), 32, 4)
	b.bgRect.FillColorScale.SetColor(styles.InterfaceBgColor)
	b.bgRect.Pos.Offset.Y = -24
	stage.AddGraphicsSlightlyAbove(b.bgRect)
	b.bgRect.Visible = false

	b.valueRect = ge.NewRect(scene.Context(), 0, 2)
	b.valueRect.FillColorScale.SetColor(styles.ProgressBarColor)
	b.valueRect.Pos.Offset.Y = -24
	stage.AddGraphicsSlightlyAbove(b.valueRect)
	b.valueRect.Visible = false
}

func (b *progressBar) SetPos(basePos *gmath.Vec) {
	b.bgRect.Pos.Base = basePos
	b.valueRect.Pos.Base = basePos
}

func (b *progressBar) SetVisibility(v bool) {
	b.bgRect.Visible = v
	b.valueRect.Visible = v
}

func (b *progressBar) SetValue(v float64) {
	if v == b.cached {
		return
	}
	b.cached = v
	percentage := v
	valueRectMaxWidth := b.bgRect.Width - 2
	b.valueRect.Width = math.Ceil(valueRectMaxWidth * percentage)
}
