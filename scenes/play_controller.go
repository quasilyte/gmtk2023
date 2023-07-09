package scenes

import (
	"fmt"
	"runtime"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/controls"
	"github.com/quasilyte/gmtk2023/eui"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/session"
	"github.com/quasilyte/gmtk2023/styles"
)

// FIXME: duplicated from battle package.
const (
	unitPanelOutline             = 2
	unitPanelPadding             = 2
	unitPanelHorizontalSeparator = 2
	unitPanelIconWidth           = 64
	unitPanelIconHeight          = 72
)

type PlayController struct {
	state *session.State

	designs *gamedata.PlayerDesigns

	scene *ge.Scene
}

func NewPlayController(state *session.State) *PlayController {
	return &PlayController{
		state:   state,
		designs: gamedata.NewPlayerDesigns(),
	}
}

func (c *PlayController) renderIcons() {
	designs := c.designs

	renderSimpleIcon(c.scene, designs.Icons[0], assets.ImageGenerator, "")

	renderTowerIcon(c.scene, designs.Icons[1], designs.Towers[0])
	renderTowerIcon(c.scene, designs.Icons[2], designs.Towers[1])

	renderSimpleIcon(c.scene, designs.Icons[3], assets.ImageRepairDepot, "")

	renderFactoryIcon(c.scene, designs.Icons[4], designs.Tanks[0])
	renderFactoryIcon(c.scene, designs.Icons[5], designs.Tanks[1])
	renderFactoryIcon(c.scene, designs.Icons[6], designs.Tanks[2])
	renderFactoryIcon(c.scene, designs.Icons[7], designs.Tanks[3])

	designs.IconConstructor = ebiten.NewImage(unitPanelIconWidth, unitPanelIconHeight)
	renderSimpleIcon(c.scene, designs.IconConstructor, assets.ImageDroneConstructor, fmt.Sprintf("%d ♦", gamedata.ConstructorEnergyCost))
	designs.IconCommander = ebiten.NewImage(unitPanelIconWidth, unitPanelIconHeight)
	renderSimpleIcon(c.scene, designs.IconCommander, assets.ImageDroneCommander, fmt.Sprintf("%d ♦", gamedata.CommanderEnergyCost))
}

func (c *PlayController) Init(scene *ge.Scene) {
	c.scene = scene

	// bigFont := assets.BitmapFont3
	// smallFont := assets.BitmapFont1

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(540, 10, nil)
	root.AddChild(rowContainer)

	grid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, true, true}, nil),
			widget.GridLayoutOpts.Spacing(4, 8),
		)))
	rowContainer.AddChild(grid)

	bodyDesigns := []*gamedata.BodyStats{
		gamedata.ScoutBodyStats,
		gamedata.FighterBodyStats,
		gamedata.HunterBodyStats,
		gamedata.DestroyerBodyStats,
	}
	turretDesigns := []*gamedata.TurretStats{
		gamedata.ScatterCannonStats,
		gamedata.LightCannonStats,
		gamedata.HurricaneStats,
		gamedata.AssaultLaserStats,
	}

	bodies := []int{
		0, 1, 2, 3,
	}
	turrets := []int{
		2, 1, 0, 3,
	}
	for i := range bodies {
		index := i
		grid.AddChild(eui.NewCenteredLabel(fmt.Sprintf("Design %d", i+1), assets.BitmapFont2))

		c.designs.Tanks[index] = &gamedata.UnitStats{
			Movement: gamedata.UnitMovementGround,
			Body:     bodyDesigns[bodies[index]],
			Turret:   turretDesigns[turrets[index]],
		}

		grid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources: c.state.UIResources,
			Input:     c.state.Input,
			Value:     &bodies[index],
			ValueNames: []string{
				"scout",
				"fighter",
				"hunter",
				"destroyer",
			},
			OnPressed: func() {
				c.designs.Tanks[index].Body = bodyDesigns[bodies[index]]
			},
		}))

		grid.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
			Resources: c.state.UIResources,
			Input:     c.state.Input,
			Value:     &turrets[index],
			ValueNames: []string{
				"scatter cannon",
				"light cannon",
				"hurricane",
				"assault laser",
			},
			OnPressed: func() {
				c.designs.Tanks[index].Turret = turretDesigns[turrets[index]]
			},
		}))
	}

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "START", func() {
		c.setDummyDesigns()
		c.renderIcons()
		config := &gamedata.BattleConfig{
			PlayerInput:   c.state.Input,
			GameSpeed:     1,
			PlayerDesigns: c.designs,
		}
		scene.Context().ChangeScene(NewBattleController(c.state, config))
	}))

	if runtime.GOARCH != "wasm" {
		rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))
		rowContainer.AddChild(eui.NewButton(c.state.UIResources, "BACK", func() {
			c.back()
		}))
	}

	initUI(scene, root)
}

func (c *PlayController) setDummyDesigns() {
	playerDesigns := c.designs
	// playerDesigns.Tanks[0] = &gamedata.UnitStats{
	// 	Movement: gamedata.UnitMovementGround,
	// 	Body:     gamedata.ScoutBodyStats,
	// 	Turret:   gamedata.HurricaneStats,
	// }
	// playerDesigns.Tanks[1] = &gamedata.UnitStats{
	// 	Movement: gamedata.UnitMovementGround,
	// 	Body:     gamedata.FighterBodyStats,
	// 	Turret:   gamedata.LightCannonStats,
	// }
	// playerDesigns.Tanks[2] = &gamedata.UnitStats{
	// 	Movement: gamedata.UnitMovementGround,
	// 	Body:     gamedata.HunterBodyStats,
	// 	Turret:   gamedata.ScatterCannonStats,
	// }
	// playerDesigns.Tanks[3] = &gamedata.UnitStats{
	// 	Movement: gamedata.UnitMovementGround,
	// 	Body:     gamedata.DestroyerBodyStats,
	// 	Turret:   gamedata.AssaultLaserStats,
	// }

	playerDesigns.Towers[0] = &gamedata.UnitStats{
		Movement:   gamedata.UnitMovementNone,
		Body:       gamedata.BunkerBodyStats,
		Turret:     gamedata.GatlingStats,
		Selectable: true,
	}
	playerDesigns.Towers[1] = &gamedata.UnitStats{
		Movement:   gamedata.UnitMovementNone,
		Body:       gamedata.BunkerBodyStats,
		Turret:     gamedata.LightCannonStats,
		Selectable: true,
	}
}

func (c *PlayController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.back()
	}
}

func (c *PlayController) back() {
	c.scene.Context().ChangeScene(NewMainMenuController(c.state))
}
