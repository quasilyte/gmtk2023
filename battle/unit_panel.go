package battle

import (
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/styles"
	"github.com/quasilyte/gmtk2023/viewport"
)

const (
	unitPanelOutline             = 2
	unitPanelPadding             = 2
	unitPanelHorizontalSeparator = 2
	unitPanelIconWidth           = 64
)

type unitPanel struct {
	bg *ge.Rect

	camera *viewport.Camera

	buttonRects []gmath.Rect

	buttonIcons []*ge.Sprite

	numButtons int

	input *input.Handler
}

func newUnitPanel(cam *viewport.Camera, h *input.Handler) *unitPanel {
	return &unitPanel{
		input:       h,
		camera:      cam,
		buttonRects: make([]gmath.Rect, 8),
		buttonIcons: make([]*ge.Sprite, 8),
	}
}

func (p *unitPanel) Init(scene *ge.Scene) {
	p.bg = ge.NewRect(scene.Context(), 0, 72+(2*unitPanelOutline)+(2*unitPanelPadding))
	p.bg.Centered = false
	p.bg.Pos.Offset = gmath.Vec{X: 4, Y: 4}
	p.bg.FillColorScale.SetColor(styles.UnitPanelBgColor)
	p.bg.OutlineColorScale.SetColor(styles.UnitPanelBgOutlineColor)
	p.bg.OutlineWidth = unitPanelOutline
	p.camera.UI.AddGraphics(p.bg)

	offsetX := float64(unitPanelOutline + unitPanelPadding)
	offsetY := float64(unitPanelOutline + unitPanelPadding)
	for i := range p.buttonIcons {
		s := ge.NewSprite(scene.Context())
		s.Centered = false
		s.Visible = false
		s.Pos.Offset = p.bg.Pos.Offset.Add(gmath.Vec{X: offsetX, Y: offsetY})
		p.buttonIcons[i] = s
		p.camera.UI.AddGraphics(s)
		offsetX += unitPanelIconWidth + unitPanelHorizontalSeparator
	}
}

func (p *unitPanel) SetButtons(icons []*ebiten.Image) {
	p.numButtons = len(icons)

	width := (2 * unitPanelOutline) + (2 * unitPanelPadding) +
		(unitPanelHorizontalSeparator * (len(icons) - 1)) +
		(unitPanelIconWidth * len(icons))
	p.bg.Width = float64(width)

	for i, s := range p.buttonIcons {
		s.Visible = i < p.numButtons
		if i >= p.numButtons {
			continue
		}
		s.SetImage(resource.Image{Data: icons[i]})
	}
}

func (p *unitPanel) HandleInput() int {
	return -1
}
