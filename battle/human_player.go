package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/controls"
	"github.com/quasilyte/gmtk2023/viewport"
)

type humanPlayer struct {
	world *worldState

	input *input.Handler

	camera *viewport.Camera

	droneSelector *ge.Sprite
	selectedUnit  *unit

	cameraPanSpeed    float64
	cameraPanBoundary float64
}

func newHumanPlayer(world *worldState) *humanPlayer {
	return &humanPlayer{
		world:             world,
		input:             world.PlayerInput,
		camera:            world.Camera,
		cameraPanSpeed:    8,
		cameraPanBoundary: 8,
	}
}

func (p *humanPlayer) Init() {
	p.droneSelector = p.world.Scene().NewSprite(assets.ImageUIDroneSelector)
	p.droneSelector.Visible = false
}

func (p *humanPlayer) Update(scaledDelta, delta float64) {
	p.panCamera(delta)
	p.handleInput()
}

func (p *humanPlayer) handleInput() {
	if info, ok := p.input.JustPressedActionInfo(controls.ActionSelectUnit); ok {
		worldPos := p.camera.AbsPos(info.Pos)
		u := p.world.FindSelectable(worldPos)
		if u != nil && p.selectedUnit != u {
			p.setSelectedUnit(u)
		}
	}
}

func (p *humanPlayer) setSelectedUnit(u *unit) {
	p.selectedUnit = u

	p.droneSelector.Visible = true
	p.droneSelector.Pos.Base = &p.selectedUnit.spritePos
	p.camera.Stage.AddSpriteSlightlyAbove(p.droneSelector)
}

func (p *humanPlayer) panCamera(delta float64) {
	var cameraPan gmath.Vec
	if p.input.ActionIsPressed(controls.ActionPanRight) {
		cameraPan.X += p.cameraPanSpeed
	}
	if p.input.ActionIsPressed(controls.ActionPanDown) {
		cameraPan.Y += p.cameraPanSpeed
	}
	if p.input.ActionIsPressed(controls.ActionPanLeft) {
		cameraPan.X -= p.cameraPanSpeed
	}
	if p.input.ActionIsPressed(controls.ActionPanUp) {
		cameraPan.Y -= p.cameraPanSpeed
	}
	if cameraPan.IsZero() && p.cameraPanBoundary != 0 {
		// Mouse cursor can pan the camera too.
		cursor := p.input.CursorPos()
		if cursor.X > p.camera.Rect.Width()-p.cameraPanBoundary {
			cameraPan.X += p.cameraPanSpeed
		}
		if cursor.Y > p.camera.Rect.Height()-p.cameraPanBoundary {
			cameraPan.Y += p.cameraPanSpeed
		}
		if cursor.X < p.cameraPanBoundary {
			cameraPan.X -= p.cameraPanSpeed
		}
		if cursor.Y < p.cameraPanBoundary {
			cameraPan.Y -= p.cameraPanSpeed
		}
	}
	if !cameraPan.IsZero() {
		p.camera.Pan(cameraPan)
	}
}
