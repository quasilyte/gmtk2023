package battle

import (
	"fmt"
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/styles"
)

type resourcesPanel struct {
	player *humanPlayer

	bg *ge.Rect

	labelNormalRes  *ge.Label
	labelEnergyRes  *ge.Label
	labelGenerators *ge.Label
}

type resourcesPanelUpdate struct {
	numResources  float64
	numEnergy     float64
	numGenerators int
}

func newResourcesPanel(p *humanPlayer) *resourcesPanel {
	return &resourcesPanel{
		player: p,
	}
}

func (p *resourcesPanel) setVisibility(visible bool) {
	p.bg.Visible = visible
	p.labelNormalRes.Visible = visible
	p.labelEnergyRes.Visible = visible
}

func (p *resourcesPanel) Update(info resourcesPanelUpdate) {
	if !p.bg.Visible {
		return
	}
	p.labelNormalRes.Text = fmt.Sprintf("Resources: %d $", int(math.Trunc(info.numResources)))
	p.labelEnergyRes.Text = fmt.Sprintf("Energy: %d ♦", int(math.Trunc(info.numEnergy)))
	p.labelGenerators.Text = fmt.Sprintf("Generators: %d", info.numGenerators)
}

func (p *resourcesPanel) Init(scene *ge.Scene) {
	width := 164.0
	screenWidth := 1920.0 / 2

	p.bg = ge.NewRect(scene.Context(), width, unitPanelIconHeight+(2*unitPanelOutline)+(2*unitPanelPadding))
	p.bg.Centered = false
	p.bg.Visible = false
	p.bg.Pos.Offset = gmath.Vec{X: screenWidth - width - 4, Y: 4}
	p.bg.FillColorScale.SetColor(styles.UnitPanelBgColor)
	p.bg.OutlineColorScale.SetColor(styles.UnitPanelBgOutlineColor)
	p.bg.OutlineWidth = unitPanelOutline
	p.player.camera.UI.AddGraphics(p.bg)

	p.labelNormalRes = ge.NewLabel(assets.BitmapFont1)
	p.labelNormalRes.AlignHorizontal = ge.AlignHorizontalCenter
	p.labelNormalRes.AlignVertical = ge.AlignVerticalCenter
	p.labelNormalRes.Width = width - 4
	p.labelNormalRes.Height = 24
	p.labelNormalRes.Pos.Offset = p.bg.Pos.Offset.Add(gmath.Vec{X: 4, Y: 16})
	p.player.camera.UI.AddGraphics(p.labelNormalRes)

	p.labelEnergyRes = ge.NewLabel(assets.BitmapFont1)
	p.labelEnergyRes.AlignHorizontal = ge.AlignHorizontalCenter
	p.labelEnergyRes.AlignVertical = ge.AlignVerticalCenter
	p.labelEnergyRes.Width = width - 4
	p.labelEnergyRes.Height = 24
	p.labelEnergyRes.Pos.Offset = p.bg.Pos.Offset.Add(gmath.Vec{X: 4, Y: 32})
	p.player.camera.UI.AddGraphics(p.labelEnergyRes)

	p.labelGenerators = ge.NewLabel(assets.BitmapFont1)
	p.labelGenerators.AlignHorizontal = ge.AlignHorizontalCenter
	p.labelGenerators.AlignVertical = ge.AlignVerticalCenter
	p.labelGenerators.Width = width - 4
	p.labelGenerators.Height = 24
	p.labelGenerators.Pos.Offset = p.bg.Pos.Offset.Add(gmath.Vec{X: 4, Y: 48})
	p.player.camera.UI.AddGraphics(p.labelGenerators)

	p.setVisibility(false)
}
