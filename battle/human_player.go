package battle

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/controls"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/viewport"
)

type humanPlayer struct {
	world *worldState

	input *input.Handler

	camera *viewport.Camera

	droneSelector    *ge.Sprite
	selectedUnit     *unit
	selectedUnitPath *ge.Line

	inactiveTankSelectors []*tankSelector
	activeTankSelectors   []*tankSelector

	unitPanel *unitPanel

	designs *gamedata.PlayerDesigns

	cameraPanSpeed    float64
	cameraPanBoundary float64
}

func newHumanPlayer(world *worldState, designs *gamedata.PlayerDesigns) *humanPlayer {
	return &humanPlayer{
		world:                 world,
		designs:               designs,
		input:                 world.PlayerInput,
		camera:                world.Camera,
		cameraPanSpeed:        8,
		cameraPanBoundary:     8,
		inactiveTankSelectors: make([]*tankSelector, 0, 32),
		activeTankSelectors:   make([]*tankSelector, 0, 32),
	}
}

func (p *humanPlayer) Init() {
	p.droneSelector = p.world.Scene().NewSprite(assets.ImageUIDroneSelector)
	p.droneSelector.Visible = false
	p.camera.Stage.AddSpriteSlightlyAbove(p.droneSelector)

	p.selectedUnitPath = ge.NewLine(ge.Pos{}, ge.Pos{})
	p.selectedUnitPath.SetColorScaleRGBA(0x4b, 0xc2, 0x75, 200)
	p.selectedUnitPath.Width = 2
	p.selectedUnitPath.Visible = false
	p.camera.Stage.AddGraphicsSlightlyAbove(p.selectedUnitPath)

	p.renderIcons()
	p.unitPanel = newUnitPanel(p.camera, p.input)
	p.unitPanel.Init(p.world.scene)
}

func (p *humanPlayer) renderIcons() {
	// TODO: this should be done somewhere else, before the battle starts.

	renderGeneratorIcon(p.world.scene, p.designs.Icons[0])

	renderTowerIcon(p.world.scene, p.designs.Icons[1], p.designs.Towers[0])
	renderTowerIcon(p.world.scene, p.designs.Icons[2], p.designs.Towers[1])

	renderFactoryIcon(p.world.scene, p.designs.Icons[3], p.designs.Tanks[0])
	renderFactoryIcon(p.world.scene, p.designs.Icons[4], p.designs.Tanks[1])
	renderFactoryIcon(p.world.scene, p.designs.Icons[5], p.designs.Tanks[2])
	renderFactoryIcon(p.world.scene, p.designs.Icons[6], p.designs.Tanks[3])
}

func (p *humanPlayer) Update(scaledDelta, delta float64) {
	p.panCamera(delta)
	p.handleInput()

	if len(p.activeTankSelectors) != 0 {
		stillActive := p.activeTankSelectors[:0]
		for _, sel := range p.activeTankSelectors {
			sel.Update()
			if sel.IsActive() {
				stillActive = append(stillActive, sel)
			} else {
				sel.SetUnit(nil)
				p.inactiveTankSelectors = append(p.inactiveTankSelectors, sel)
			}
		}
		p.activeTankSelectors = stillActive
	}
}

func (p *humanPlayer) executeDeconstructableAction(actionIndex int) bool {
	if actionIndex != 0 {
		panic("unreachable")
	}
	p.selectedUnit.Deconstruct()
	return true
}

func (p *humanPlayer) executeConstructorAction(actionIndex int) bool {
	pos := p.world.FindConstructionSitePos(p.selectedUnit.pos)
	if pos.IsZero() {
		return false
	}

	var stats *gamedata.UnitStats
	var newUnitExtra any
	time := 0.0
	switch actionIndex {
	case 0:
		stats = gamedata.GeneratorUnitStats
		time = stats.ConstructionTime
	case 1, 2:
		index := actionIndex - 1
		stats = gamedata.TowerConstruction
		newUnitExtra = p.designs.Towers[index]
		time = p.designs.Towers[index].Turret.ProductionTime * 3
	case 3, 4, 5, 6:
		index := actionIndex - 3
		stats = gamedata.TankFactoryUnitStats
		newUnitExtra = &tankFactoryExtra{tankDesign: p.designs.Tanks[index]}
		time = stats.ConstructionTime
	}
	if stats == gamedata.TankFactoryUnitStats && newUnitExtra.(*tankFactoryExtra).tankDesign.Body.Heavy {
		stats = gamedata.HeavyTankFactoryUnitStats
	}
	p.selectedUnit.extra = &constructionOrder{
		siteStats: stats,
		siteExtra: &constructionSiteExtra{
			newUnitExtra: newUnitExtra,
			goalProgress: time,
		},
	}
	p.selectedUnit.SendTo(pos)
	return true
}

func (p *humanPlayer) executeUnitAction(actionIndex int) bool {
	if p.selectedUnit.IsConstructor() {
		return p.executeConstructorAction(actionIndex)
	}
	if p.selectedUnit.IsSimpleDeconstructible() {
		return p.executeDeconstructableAction(actionIndex)
	}
	return true
}

func (p *humanPlayer) handleInput() {
	if p.selectedUnit != nil && p.unitPanel.bg.Visible {
		actionIndex := p.unitPanel.HandleInput()
		if actionIndex != -1 {
			if p.executeUnitAction(actionIndex) {
				playGlobalSound(p.world, assets.AudioUnitAck1)
			} else {
				playGlobalSound(p.world, assets.AudioError)
			}
			return
		}
	}

	if p.selectedUnit != nil {
		if p.selectedUnit.stats.Movement != gamedata.UnitMovementNone {
			if info, ok := p.input.JustPressedActionInfo(controls.ActionSendUnit); ok {
				worldPos := p.camera.AbsPos(info.Pos)
				p.selectedUnit.SendTo(worldPos)
				p.updateUnitPath(p.selectedUnit)
				playGlobalSound(p.world, assets.AudioUnitAck1)
			}
		}

		if p.selectedUnit.NeedsMoreConstructors() {
			if info, ok := p.input.JustPressedActionInfo(controls.ActionAddToGroup); ok {
				worldPos := p.camera.AbsPos(info.Pos)
				u := p.world.FindConstructor(worldPos)
				if u != nil {
					u.SendTo(p.selectedUnit.pos)
					u.extra = &constructorEntryTarget{site: p.selectedUnit}
					return
				}
			}
		}

		if p.selectedUnit.IsCommander() {
			if info, ok := p.input.JustPressedActionInfo(controls.ActionAddToGroup); ok {
				worldPos := p.camera.AbsPos(info.Pos)
				u := p.world.FindAssignable(worldPos)
				if u != nil {
					if u.leader == p.selectedUnit {
						u.leader = nil
						for _, sel := range p.activeTankSelectors {
							if sel.GetUnit() == u {
								sel.SetUnit(nil)
								break
							}
						}
					} else {
						if len(p.selectedUnit.group) < gamedata.MaxGroupSize {
							u.leader = p.selectedUnit
							p.activeTankSelectors = append(p.activeTankSelectors, p.createTankSelector(u))
						}
					}
				}
			}
		}
	}

	if info, ok := p.input.JustPressedActionInfo(controls.ActionSelectUnit); ok {
		worldPos := p.camera.AbsPos(info.Pos)
		u := p.world.FindSelectable(worldPos)
		if u != nil && p.selectedUnit != u {
			p.setSelectedUnit(u)
		}
	}
}

func (p *humanPlayer) updateUnitPath(u *unit) {
	if u != nil {
		p.selectedUnitPath.BeginPos.Base = &p.selectedUnit.spritePos
		p.selectedUnitPath.EndPos.Offset = p.selectedUnit.waypoint
		p.selectedUnitPath.Visible = !p.selectedUnit.waypoint.IsZero()
	} else {
		p.selectedUnitPath.Visible = false
	}
}

func (p *humanPlayer) IsDisposed() bool { return false }

func (p *humanPlayer) setSelectedUnit(u *unit) {
	if p.selectedUnit != nil {
		p.selectedUnit.EventDestroyed.Disconnect(p)
		p.unitPanel.SetButtons(nil)
	}

	p.selectedUnit = u

	p.droneSelector.Visible = u != nil
	p.updateUnitPath(u)

	if u != nil {
		img := assets.ImageUIDroneSelector
		if u.stats.Large {
			img = assets.ImageUILargeSelector
		} else if u.IsBuilding() {
			img = assets.ImageUITowerSelector
		}
		if p.droneSelector.ImageID() != img {
			p.droneSelector.SetImage(p.world.scene.LoadImage(img))
		}

		switch {
		case u.IsConstructor():
			p.unitPanel.SetButtons(p.designs.Icons[:7])
		case u.IsSimpleDeconstructible():
			p.unitPanel.SetButtons([]*ebiten.Image{p.world.Scene().LoadImage(assets.ImageUIDeconstuctIcon).Data})
		}
	}

	if len(p.activeTankSelectors) != 0 {
		for _, sel := range p.activeTankSelectors {
			sel.SetUnit(nil)
			p.inactiveTankSelectors = append(p.inactiveTankSelectors, sel)
		}
		p.activeTankSelectors = p.activeTankSelectors[:0]
	}

	if u != nil {
		p.selectedUnit.EventDestroyed.Connect(p, func(u *unit) {
			if p.selectedUnit == u {
				p.setSelectedUnit(nil)
			}
		})

		p.droneSelector.Pos.Base = &p.selectedUnit.spritePos

		p.activeTankSelectors = append(p.activeTankSelectors, p.createTankSelector(u))
		for _, gu := range u.group {
			p.activeTankSelectors = append(p.activeTankSelectors, p.createTankSelector(gu))
		}
	}
}

func (p *humanPlayer) createTankSelector(u *unit) *tankSelector {
	if len(p.inactiveTankSelectors) != 0 {
		sel := p.inactiveTankSelectors[len(p.inactiveTankSelectors)-1]
		p.inactiveTankSelectors = p.inactiveTankSelectors[:len(p.inactiveTankSelectors)-1]
		sel.SetUnit(u)
		return sel
	}

	sel := newTankSelector()
	sel.Init(p.world.Scene(), p.camera.Stage)
	sel.SetUnit(u)
	return sel
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
