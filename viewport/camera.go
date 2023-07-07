package viewport

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/gmath"
)

type Camera struct {
	Stage *Stage

	Offset gmath.Vec

	width     float64
	height    float64
	WorldRect gmath.Rect

	Rect       gmath.Rect
	globalRect gmath.Rect

	screen *ebiten.Image
}

func NewCamera(stage *Stage, world gmath.Rect, width, height float64) *Camera {
	cam := &Camera{
		Stage: stage,

		WorldRect: world,
		width:     world.Width(),
		height:    world.Height(),

		Rect: gmath.Rect{
			Min: gmath.Vec{},
			Max: gmath.Vec{X: width, Y: height},
		},
		screen: ebiten.NewImage(int(width), int(height)),
	}
	return cam
}

func (c *Camera) IsDisposed() bool { return false }

func (c *Camera) Draw(screen *ebiten.Image) {
	c.globalRect = c.Rect
	c.globalRect.Min = c.Offset
	c.globalRect.Max = c.globalRect.Max.Add(c.Offset)

	c.screen.Clear()
	drawOffset := gmath.Vec{
		X: -c.Offset.X,
		Y: -c.Offset.Y,
	}
	c.Stage.bg.DrawPartialWithOffset(c.screen, c.globalRect, drawOffset)

	var options ebiten.DrawImageOptions
	screen.DrawImage(c.screen, &options)
}

func (c *Camera) ContainsPos(pos gmath.Vec) bool {
	globalRect := c.Rect
	globalRect.Min = c.Offset
	globalRect.Max = globalRect.Max.Add(c.Offset)
	return globalRect.Contains(pos)
}

func (c *Camera) checkBounds() {
	c.Offset.X = gmath.Clamp(c.Offset.X, 0, c.width-c.Rect.Width())
	c.Offset.Y = gmath.Clamp(c.Offset.Y, 0, c.height-c.Rect.Height())
}

func (c *Camera) Pan(delta gmath.Vec) {
	if delta.IsZero() {
		return
	}
	c.Offset = c.Offset.Add(delta)
	c.checkBounds()
}

func (c Camera) CenterPos() gmath.Vec {
	return c.Offset.Add(c.Rect.Center())
}

func (c *Camera) CenterOn(pos gmath.Vec) {
	c.Offset = pos.Sub(c.Rect.Center())
	c.checkBounds()
}

func (c *Camera) SetOffset(pos gmath.Vec) {
	c.Offset = pos
	c.checkBounds()
}
