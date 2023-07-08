package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/pathing"
	"github.com/quasilyte/gsignal"
)

type unit struct {
	stats *gamedata.UnitStats

	world *worldState

	pos        gmath.Vec
	spritePos  gmath.Vec
	turretPos  gmath.Vec
	frameWidth float64

	path             pathing.GridPath
	partialPathSteps int

	finalWaypoint gmath.Vec
	waypoint      gmath.Vec
	rotation      gmath.Rad
	dstRotation   gmath.Rad

	sprite *ge.Sprite
	anim   *ge.Animation

	turret *turret

	leader *unit
	group  []*unit

	extra any

	hp    float64
	maxHP float64

	needRotate bool
	disposed   bool

	EventDestroyed gsignal.Event[*unit]
}

type tankFactoryExtra struct {
	tankDesign *gamedata.UnitStats
}

type constructorEntryTarget struct {
	site *unit
}

type constructionOrder struct {
	siteStats *gamedata.UnitStats
	siteExtra *constructionSiteExtra
}

type constructionSiteExtra struct {
	constructors int

	newUnitExtra any

	percentage   float64
	progress     float64
	maxProgress  float64
	goalProgress float64
}

type unitConfig struct {
	Stats *gamedata.UnitStats
	Pos   gmath.Vec
}

func newUnit(world *worldState, config unitConfig) *unit {
	u := &unit{
		stats:      config.Stats,
		world:      world,
		pos:        config.Pos,
		frameWidth: config.Stats.Body.Texture.DefaultFrameWidth,
	}
	u.maxHP = config.Stats.Body.HP
	if u.stats.Turret != nil {
		u.maxHP += config.Stats.Turret.HP
	}
	u.hp = u.maxHP
	return u
}

func (u *unit) IsDisposed() bool {
	return u.disposed
}

func (u *unit) IsCommander() bool { return u.stats == gamedata.CommanderUnitStats }

func (u *unit) IsConstructor() bool { return u.stats == gamedata.ConstructorUnitStats }

func (u *unit) IsConstructionSite() bool {
	_, ok := u.extra.(*constructionSiteExtra)
	return ok
}

func (u *unit) NeedsMoreConstructors() bool {
	extra, ok := u.extra.(*constructionSiteExtra)
	if !ok {
		return false
	}
	return extra.constructors < u.stats.ConstructorsNeeded
}

func (u *unit) IsSimpleDeconstructible() bool {
	if u.IsBuilding() && u.turret != nil {
		return true
	}
	if u.stats == gamedata.GeneratorUnitStats {
		return true
	}

	switch u.extra.(type) {
	case *tankFactoryExtra, *constructionSiteExtra:
		return true
	default:
		return false
	}
}

func (u *unit) IsBuilding() bool { return u.stats.Movement == gamedata.UnitMovementNone }

func (u *unit) Dispose() {
	u.disposed = true
	u.sprite.Dispose()

	if u.turret != nil {
		u.turret.Dispose()
	}

	if u.IsBuilding() {
		u.world.UnmarkPos(u.pos)
	}

	u.EventDestroyed.Emit(u)
}

func (u *unit) updatePos() {
	u.spritePos.X = math.Round(u.pos.X)
	u.spritePos.Y = math.Round(u.pos.Y)
	u.turretPos = u.spritePos.Add(gmath.Vec{Y: u.stats.Body.TurretOffset})
}

func (u *unit) Init(scene *ge.Scene) {
	u.updatePos()

	if u.IsBuilding() {
		u.world.MarkPos(u.pos)
	}

	if u.stats.Body.Image != assets.ImageNone {
		u.sprite = scene.NewSprite(u.stats.Body.Image)
		u.sprite.Pos.Base = &u.spritePos
		if u.sprite.ImageWidth() != u.sprite.FrameWidth {
			u.anim = ge.NewRepeatedAnimation(u.sprite, -1)
			u.anim.SetAnimationSpan(0.5)
		}
		if u.stats.Movement == gamedata.UnitMovementHover {
			u.world.Stage().AddSpriteSlightlyAbove(u.sprite)
		} else {
			u.world.Stage().AddSprite(u.sprite)
		}
	} else {
		u.sprite = ge.NewSprite(scene.Context())
		u.sprite.SetImage(u.stats.Body.Texture)
		u.sprite.Pos.Base = &u.spritePos
		u.world.Stage().AddSprite(u.sprite)
	}

	if u.stats.Creep {
		u.sprite.SetHue(gmath.DegToRad(80))
	}
	if u.IsConstructionSite() {
		u.sprite.Shader = scene.NewShader(assets.ShaderConstructionLarge)
		u.sprite.Shader.SetFloatValue("Time", 0.15)
	}

	if u.stats.Turret != nil {
		u.turret = newTurret(u.world, turretConfig{
			Image: u.stats.Turret.Texture,
			Pos:   &u.turretPos,
			Owner: u,
		})
		u.world.runner.AddObject(u.turret)
	}
}

func (u *unit) Update(delta float64) {
	if u.anim != nil {
		u.anim.Tick(delta)
	}

	if u.IsConstructionSite() {
		u.updateConstructionSite(delta)
	}

	u.updatePos()

	if !u.waypoint.IsZero() {
		u.moveToWaypoint(delta)
	}
}

func (u *unit) Deconstruct() {
	hpPercentage := u.hp / u.maxHP
	for i, gu := range u.group {
		released := u.world.NewUnit(unitConfig{
			Stats: gu.stats,
			Pos:   u.pos.Add(deconstructSpawnOffsets[i]),
		})
		u.world.runner.AddObject(released)
		if i != 0 {
			released.SendTo(released.pos.Add(deconstructWaypointOffsets[i]))
		}
		released.hp = released.maxHP * hpPercentage
		effect := newEffectNode(u.world, released.pos, true, assets.ImageConstructorMerge)
		effect.rotates = true
		u.world.runner.AddObject(effect)
	}
	u.Dispose()
}

func (u *unit) AddConstructorToSite(constructor *unit) bool {
	if u.IsDisposed() {
		return false
	}
	extra := u.extra.(*constructionSiteExtra)
	if extra.constructors >= u.stats.ConstructorsNeeded {
		return false
	}
	extra.constructors++

	effect := newEffectNode(u.world, constructor.pos, false, assets.ImageConstructorMerge)
	effect.rotates = true
	u.world.runner.AddObject(effect)

	constructor.Dispose()
	u.group = append(u.group, constructor)
	if extra.constructors >= u.stats.ConstructorsNeeded {
		extra.maxProgress = extra.goalProgress
	} else {
		extra.maxProgress = extra.goalProgress * (float64(extra.constructors) / float64(u.stats.ConstructorsNeeded))
	}
	return true
}

func (u *unit) updateConstructionSite(delta float64) {
	extra := u.extra.(*constructionSiteExtra)
	if extra.progress >= extra.maxProgress {
		return // Not enough constructors to continue
	}

	extra.progress += delta
	extra.percentage = extra.progress / extra.goalProgress
	u.sprite.Shader.SetFloatValue("Time", extra.percentage+0.15)
	if extra.progress >= extra.goalProgress {
		stats := u.stats
		if statsOverride, ok := extra.newUnitExtra.(*gamedata.UnitStats); ok {
			stats = statsOverride
			extra.newUnitExtra = nil
		}

		totalPercentage := 0.0
		for _, gu := range u.group {
			totalPercentage += gu.hp / gu.maxHP
		}
		totalPercentage /= float64(len(u.group))

		effect := newEffectNode(u.world, u.pos, true, assets.ImageConstructorMerge)
		effect.rotates = true
		u.world.runner.AddObject(effect)
		building := u.world.NewUnit(unitConfig{
			Stats: stats,
			Pos:   u.pos,
		})
		building.hp = building.maxHP * totalPercentage
		u.world.runner.AddObject(building)
		building.group = u.group
		building.extra = extra.newUnitExtra
		u.Dispose()
		return
	}
}

func (u *unit) setRotation(v gmath.Rad) {
	u.rotation = v
	spriteAngle := u.rotation.Normalized() - (gamedata.TankFrameAngleStep / 2)
	u.sprite.FrameOffset.X = u.frameWidth * math.Trunc(float64(spriteAngle/gamedata.TankFrameAngleStep))
}

func (u *unit) OnDamage(d gamedata.DamageValue, attacker *unit) {
	u.hp -= d.Health
	if u.hp <= 0 {
		u.Destroy()
		return
	}

	if u.turret != nil {
		u.turret.OnAttacked(attacker)
	}
}

func (u *unit) Destroy() {
	createAreaExplosion(u.world, spriteRect(u.pos, u.stats.Body.Size.X, u.stats.Body.Size.Y), true)

	u.Dispose()
}

func (u *unit) SendTo(pos gmath.Vec) {
	if u.IsCommander() {
		u.sendCommanderTo(pos)
		return
	}
	if u.IsConstructor() {
		u.sendConstructorTo(pos)
		return
	}

	switch u.stats.Movement {
	case gamedata.UnitMovementGround:
		if !u.world.pathgrid.CellIsFree(u.world.pathgrid.PosToCoord(pos)) {
			var slider gmath.Slider
			slider.SetBounds(0, len(groupOffsets)-1)
			slider.TrySetValue(u.world.Rand().IntRange(0, len(groupOffsets)-1))
			for i := 0; i < len(groupOffsets); i++ {
				offset := groupOffsets[slider.Value()]
				newPos := pos.Add(offset)
				slider.Inc()
				if u.world.pathgrid.CellIsFree(u.world.pathgrid.PosToCoord(newPos)) {
					pos = newPos
				}
			}
		}
		if u.pos.DistanceSquaredTo(pos) < 2 {
			return
		}
		p := u.world.BuildPath(u.pos, pos)
		if p.Partial {
			u.partialPathSteps = u.world.Rand().IntRange(3, 6)
		} else {
			u.partialPathSteps = -1
		}
		if p.Steps.Len() == 0 {
			return
		}
		u.path = p.Steps
		alignedPos := u.world.pathgrid.AlignPos(u.pos)
		if alignedPos.DistanceSquaredTo(u.pos) < 1 {
			u.waypoint = posMove(u.pos, u.path.Next())
		} else {
			u.waypoint = alignedPos
		}
		u.finalWaypoint = pos
		u.setDstRotation(u.pos.AngleToPoint(u.waypoint).Normalized())
	}
}

func (u *unit) setDstRotation(v gmath.Rad) {
	u.needRotate = true
	u.dstRotation = v
}

func (u *unit) sendConstructorTo(pos gmath.Vec) {
	u.waypoint = u.world.pathgrid.AlignPos(pos)
}

func (u *unit) sendCommanderTo(pos gmath.Vec) {
	u.waypoint = u.world.pathgrid.AlignPos(pos)

	alignedPos := u.world.pathgrid.AlignPos(pos)
	var slider gmath.Slider
	slider.SetBounds(0, len(groupOffsets)-1)
	slider.TrySetValue(u.world.Scene().Rand().IntRange(0, len(groupOffsets)-1))
	for _, gu := range u.group {
		offset := groupOffsets[slider.Value()]
		slider.Inc()
		gu.SendTo(alignedPos.Add(offset))
	}
}

func (u *unit) calcSpeed() float64 {
	return u.stats.Body.Speed
}

func (u *unit) setGroundUnitWaypoint(pos gmath.Vec) {
	u.waypoint = pos
	dstRotation := u.pos.AngleToPoint(u.waypoint).Normalized()
	u.setDstRotation(dstRotation)
	if u.turret != nil {
		if u.world.Rand().Chance(0.7) {
			if u.world.Rand().Chance(0.2) {
				dstRotation = (dstRotation + gmath.Rad(u.world.Rand().FloatRange(-0.3, 0.3))).Normalized()
			}
			u.turret.AlignRequest(dstRotation)
		}
	}
}

func (u *unit) moveGroundUnitToWaypoint(delta float64) {
	if u.needRotate {
		u.setRotation(u.rotation.RotatedTowards(u.dstRotation, u.stats.Body.RotationSpeed*gmath.Rad(delta)))
		if u.rotation == u.dstRotation {
			u.needRotate = false
		}
	}
	if u.needRotate {
		return
	}

	travelled := u.calcSpeed() * delta
	u.pos = u.pos.MoveTowards(u.waypoint, travelled)
	if u.pos != u.waypoint {
		return
	}

	const maxFinalWaypointDistSqr = (pathing.CellSize * pathing.CellSize * 4)
	if u.path.HasNext() {
		if u.partialPathSteps == 0 {
			if u.finalWaypoint.DistanceSquaredTo(u.pos) < maxFinalWaypointDistSqr {
				u.groundUnitStop()
				return
			}
			if u.world.Rand().Chance(0.4) {
				probe := u.pos.Add(gmath.RandElem(u.world.Rand(), groupOffsets))
				if u.world.pathgrid.CellIsFree(u.world.pathgrid.PosToCoord(probe)) {
					u.setGroundUnitWaypoint(probe)
					return
				}
			}
			u.SendTo(u.finalWaypoint.Add(gmath.RandElem(u.world.Rand(), groupOffsets)))
			return
		}
		if u.partialPathSteps > 0 {
			u.partialPathSteps--
		}
		d := u.path.Next()
		aligned := u.world.pathgrid.AlignPos(u.pos)
		u.setGroundUnitWaypoint(posMove(aligned, d))
	} else {
		if u.finalWaypoint.DistanceSquaredTo(u.pos) > maxFinalWaypointDistSqr {
			u.SendTo(u.finalWaypoint)
			return
		}
		u.groundUnitStop()
	}
}

func (u *unit) groundUnitStop() {
	u.waypoint = gmath.Vec{}
	u.needRotate = false
}

func (u *unit) moveToWaypoint(delta float64) {
	switch u.stats.Movement {
	case gamedata.UnitMovementHover:
		if entryTarget, ok := u.extra.(*constructorEntryTarget); ok {
			if entryTarget.site.IsDisposed() {
				u.extra = nil
				u.SendTo(u.world.pathgrid.AlignPos(u.pos))
				return
			}
		}
		travelled := u.calcSpeed() * delta
		u.pos = u.pos.MoveTowards(u.waypoint, travelled)
		if u.pos == u.waypoint {
			u.waypoint = gmath.Vec{}
			switch extra := u.extra.(type) {
			case *constructorEntryTarget:
				if !extra.site.AddConstructorToSite(u) {
					u.extra = nil
					u.SendTo(u.pos.Add(gmath.RandElem(u.world.Rand(), groupOffsets)))
				}
			case *constructionOrder:
				site := u.world.NewUnit(unitConfig{
					Stats: extra.siteStats,
					Pos:   u.pos,
				})
				site.hp = (u.hp / u.maxHP) * site.maxHP
				site.extra = extra.siteExtra
				site.AddConstructorToSite(u)
				u.world.runner.AddObject(site)
			}
		}

	case gamedata.UnitMovementGround:
		u.moveGroundUnitToWaypoint(delta)
	}
}
