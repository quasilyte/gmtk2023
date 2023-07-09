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

	unstuckDelay float64

	harvestDelay   float64
	normalResource float64
	energyResource float64
	numGenerators  int

	progressBar         *progressBar
	constructorsCounter *ge.Sprite
	droneSelector       *ge.Sprite
	selectedUnit        *unit
	selectedUnitPath    *ge.Line

	inactiveTankSelectors []*tankSelector
	activeTankSelectors   []*tankSelector

	unitPanel      *unitPanel
	resourcesPanel *resourcesPanel

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

	p.normalResource = 100
	p.energyResource = 160

	p.selectedUnitPath = ge.NewLine(ge.Pos{}, ge.Pos{})
	p.selectedUnitPath.SetColorScaleRGBA(0x4b, 0xc2, 0x75, 200)
	p.selectedUnitPath.Width = 2
	p.selectedUnitPath.Visible = false
	p.camera.Stage.AddGraphicsAbove(p.selectedUnitPath)

	p.progressBar = newProgressBar()
	p.progressBar.Init(p.world.scene, p.camera.Stage)
	p.progressBar.SetVisibility(false)

	p.constructorsCounter = ge.NewSprite(p.world.scene.Context())
	p.constructorsCounter.Visible = false
	p.constructorsCounter.Pos.Offset.Y = -16
	p.camera.Stage.AddSpriteAbove(p.constructorsCounter)

	p.unitPanel = newUnitPanel(p.camera, p.input)
	p.unitPanel.Init(p.world.scene)

	p.resourcesPanel = newResourcesPanel(p)
	p.resourcesPanel.Init(p.world.scene)

	p.world.EventUnitCreated.Connect(p, func(u *unit) {
		if !u.IsGenerator() {
			return
		}
		if u.IsConstructionSite() {
			return
		}
		p.numGenerators++
		u.EventDestroyed.Connect(p, func(u *unit) {
			p.numGenerators--
		})
	})
}

func (p *humanPlayer) Update(scaledDelta, delta float64) {
	p.panCamera(delta)
	p.handleInput()
	p.maybeHarvest(scaledDelta)

	p.unstuckDelay -= scaledDelta
	if p.unstuckDelay <= 0 {
		p.unstuckDelay = p.world.Rand().FloatRange(0.5, 1.2)
		p.doUnstuck()
	}

	// FIXME: this is a hack.
	if p.selectedUnit != nil {
		if p.selectedUnit.finalWaypoint.IsZero() {
			p.updateUnitPath(nil)
		}
	}

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

func (p *humanPlayer) maybeHarvest(delta float64) {
	p.harvestDelay = gmath.ClampMin(p.harvestDelay-delta, 0)
	if p.harvestDelay > 0 {
		return
	}
	p.harvestDelay = 1

	p.normalResource += 5

	// 0.7, 0.55, 0.40, 0.25, 0.10, 0.10, ...
	nextGeneratorYield := 0.7
	generated := 1.0
	for i := 0; i < p.numGenerators; i++ {
		if nextGeneratorYield <= gmath.Epsilon {
			break
		}
		generated += nextGeneratorYield
		nextGeneratorYield = gmath.ClampMin(nextGeneratorYield-0.15, 0.1)
	}
	p.energyResource += generated

	p.updateResourcesPanel()
}

func (p *humanPlayer) updateResourcesPanel() {
	p.resourcesPanel.Update(resourcesPanelUpdate{
		numResources:  p.normalResource,
		numEnergy:     p.energyResource,
		numGenerators: p.numGenerators,
	})
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
	if p.world.MayBlockFactory(pos) {
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
	case 3:
		if !p.world.IsInnerPos(pos) {
			return false
		}
		stats = gamedata.RepairDepotUnitStats
		time = stats.ConstructionTime
	case 4, 5, 6, 7:
		if !p.world.IsInnerPos(pos) {
			return false
		}
		for _, offset := range factoryCheckOffsets {
			if !p.world.pathgrid.CellIsFree(p.world.pathgrid.PosToCoord(pos.Add(offset))) {
				return false
			}
		}

		index := actionIndex - 4
		stats = gamedata.TankFactoryUnitStats
		extra := &tankFactoryExtra{tankDesign: p.designs.Tanks[index]}
		newUnitExtra = extra
		if stats == gamedata.TankFactoryUnitStats && extra.tankDesign.Body.Heavy {
			stats = gamedata.HeavyTankFactoryUnitStats
		}
		time = stats.ConstructionTime
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

func (p *humanPlayer) doUnstuck() {
	constructor := randIterate(p.world.Rand(), p.world.playerUnits.selectable, func(u *unit) bool {
		return u.IsConstructor() && u.waypoint.IsZero()
	})
	if constructor == nil {
		return
	}

	for _, u := range p.world.playerUnits.selectable {
		if !u.waypoint.IsZero() || u == constructor {
			continue
		}
		if u.pos != constructor.pos {
			continue
		}

		offset := randIterate(p.world.Rand(), groupOffsets, func(offset gmath.Vec) bool {
			probe := u.pos.Add(offset)
			if p.world.pathgrid.CellIsFree(p.world.pathgrid.PosToCoord(probe)) {
				return true
			}
			return false
		})
		if offset.IsZero() {
			return
		}
		constructor.SendTo(constructor.pos.Add(offset))
		if constructor == p.selectedUnit {
			p.updateUnitPath(p.selectedUnit)
		}
		break
	}
}

func (p *humanPlayer) executeUnitAction(actionIndex int) bool {
	if p.selectedUnit.IsConstructor() {
		return p.executeConstructorAction(actionIndex)
	}
	if p.selectedUnit.IsSimpleDeconstructible() {
		return p.executeDeconstructableAction(actionIndex)
	}
	if p.selectedUnit.IsMCV() {
		return p.executeMCVAction(actionIndex)
	}
	return true
}

func (p *humanPlayer) executeMCVAction(actionIndex int) bool {
	if !p.world.IsInnerPos(p.selectedUnit.pos) {
		return false
	}

	var stats *gamedata.UnitStats

	switch actionIndex {
	case 0:
		if p.energyResource < gamedata.ConstructorEnergyCost {
			return false
		}
		p.energyResource -= gamedata.ConstructorEnergyCost
		stats = gamedata.ConstructorUnitStats

	case 1:
		if p.energyResource < gamedata.CommanderEnergyCost {
			return false
		}
		p.energyResource -= gamedata.CommanderEnergyCost
		stats = gamedata.CommanderUnitStats

	default:
		return false
	}

	if stats == nil {
		return false
	}

	spawnPos := p.selectedUnit.pos.Add(gmath.RandElem(p.world.Rand(), groupOffsets))
	u := p.world.NewUnit(unitConfig{
		Stats: stats,
		Pos:   p.selectedUnit.pos,
	})
	p.world.runner.AddObject(u)
	u.SendTo(spawnPos)

	effect := newEffectNode(u.world, u.pos, false, assets.ImageConstructorMerge)
	effect.rotates = true
	u.world.runner.AddObject(effect)

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
				playGlobalSound(p.world, assets.AudioUnitAck1)
				p.selectedUnit.SendTo(worldPos)
				p.updateUnitPath(p.selectedUnit)
				if p.selectedUnit.IsConstructor() {
					constructionSite := p.world.FindConstructionSiteAt(worldPos)
					if constructionSite != nil {
						p.selectedUnit.extra = &constructorEntryTarget{site: constructionSite}
					}
				}
				return
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
							u.SendTo(p.selectedUnit.pos.Add(gmath.RandElem(p.world.Rand(), groupOffsets)))
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
		p.selectedUnitPath.EndPos.Offset = p.selectedUnit.finalWaypoint
		p.selectedUnitPath.Visible = !p.selectedUnit.finalWaypoint.IsZero()
	} else {
		p.selectedUnitPath.Visible = false
	}
}

func (p *humanPlayer) IsDisposed() bool { return false }

func (p *humanPlayer) setSelectedUnit(u *unit) {
	if p.selectedUnit != nil {
		p.selectedUnit.EventDestroyed.Disconnect(p)
		p.selectedUnit.EventConstructorEntered.Disconnect(p)
		p.selectedUnit.EventReselectRequest.Disconnect(p)
		p.selectedUnit.EventProductionProgress.Disconnect(p)
		p.unitPanel.SetButtons(nil)
	}

	p.selectedUnit = u

	p.constructorsCounter.Visible = false
	p.progressBar.SetVisibility(false)
	p.resourcesPanel.setVisibility(false)
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
		case u.IsMCV():
			p.unitPanel.SetButtons([]*ebiten.Image{
				p.designs.IconConstructor,
				p.designs.IconCommander,
			})
			p.resourcesPanel.setVisibility(true)
			p.updateResourcesPanel()
		case u.IsConstructor():
			p.unitPanel.SetButtons(p.designs.Icons)
		case u.IsTankFactory():
			p.progressBar.SetVisibility(true)
			p.progressBar.SetPos(&u.spritePos)
			p.progressBar.SetValue(u.extra.(*tankFactoryExtra).percentage)
			p.unitPanel.SetButtons([]*ebiten.Image{p.world.Scene().LoadImage(assets.ImageUIDeconstuctIcon).Data})
		case u.IsConstructionSite():
			p.constructorsCounter.Visible = true
			p.constructorsCounter.Pos.Base = &u.spritePos
			p.updateConsnstructorsCounter(u)
			p.unitPanel.SetButtons([]*ebiten.Image{p.world.Scene().LoadImage(assets.ImageUIDeconstuctIcon).Data})
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
		p.selectedUnit.EventReselectRequest.Connect(p, func(u *unit) {
			p.setSelectedUnit(u)
		})
		if p.selectedUnit.IsTankFactory() {
			p.selectedUnit.EventProductionProgress.Connect(p, func(percentage float64) {
				p.progressBar.SetValue(percentage)
			})
		}
		if p.selectedUnit.IsConstructionSite() {
			p.selectedUnit.EventConstructorEntered.Connect(p, func(u *unit) {
				p.updateConsnstructorsCounter(u)
			})
		}

		p.droneSelector.Pos.Base = &p.selectedUnit.spritePos

		p.activeTankSelectors = append(p.activeTankSelectors, p.createTankSelector(u))
		for _, gu := range u.group {
			p.activeTankSelectors = append(p.activeTankSelectors, p.createTankSelector(gu))
		}
	}
}

func (p *humanPlayer) updateConsnstructorsCounter(u *unit) {
	imageID := assets.ImageUIConstructors1outof1
	extra := u.extra.(*constructionSiteExtra)
	switch u.stats.ConstructorsNeeded {
	case 1:
		// 1outof1
	case 2:
		switch extra.constructors {
		case 1:
			imageID = assets.ImageUIConstructors1outof2
		case 2:
			imageID = assets.ImageUIConstructors2outof2
		}
	case 3:
		switch extra.constructors {
		case 1:
			imageID = assets.ImageUIConstructors1outof3
		case 2:
			imageID = assets.ImageUIConstructors2outof3
		case 3:
			imageID = assets.ImageUIConstructors3outof3
		}
	}
	p.constructorsCounter.SetImage(p.world.Scene().LoadImage(imageID))
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
