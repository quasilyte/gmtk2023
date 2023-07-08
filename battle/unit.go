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

	pos       gmath.Vec
	spritePos gmath.Vec
	turretPos gmath.Vec

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

	hp    float64
	maxHP float64

	needRotate bool
	disposed   bool

	EventDestroyed gsignal.Event[*unit]
}

type unitConfig struct {
	Stats *gamedata.UnitStats
	Pos   gmath.Vec
}

func newUnit(world *worldState, config unitConfig) *unit {
	u := &unit{
		stats: config.Stats,
		world: world,
		pos:   config.Pos,
	}
	u.maxHP = config.Stats.Body.HP
	if u.stats.Turret != nil {
		u.maxHP += config.Stats.Turret.HP
	}
	u.hp = u.maxHP
	u.hp *= world.Scene().Rand().FloatRange(0.2, 0.9)
	return u
}

func (u *unit) IsDisposed() bool {
	return u.disposed
}

func (u *unit) IsCommander() bool { return u.stats == gamedata.CommanderUnitStats }

func (u *unit) IsTower() bool { return u.stats.Movement == gamedata.UnitMovementNone }

func (u *unit) Dispose() {
	u.disposed = true
	u.sprite.Dispose()

	if u.IsTower() {
		u.world.UnmarkPos(u.pos)
	}
}

func (u *unit) updatePos() {
	u.spritePos.X = math.Round(u.pos.X)
	u.spritePos.Y = math.Round(u.pos.Y)
	u.turretPos = u.spritePos.Add(gmath.Vec{Y: u.stats.Body.TurretOffset})
}

func (u *unit) Init(scene *ge.Scene) {
	u.updatePos()

	if u.IsTower() {
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

	if u.stats.Turret != nil {
		u.turret = newTurret(u.world, turretConfig{
			Image: u.stats.Turret.Texture,
			Pos:   &u.turretPos,
		})
		u.world.runner.AddObject(u.turret)
	}
}

func (u *unit) Update(delta float64) {
	if u.anim != nil {
		u.anim.Tick(delta)
	}

	u.updatePos()

	if !u.waypoint.IsZero() {
		u.moveToWaypoint(delta)
	}
}

func (u *unit) setRotation(v gmath.Rad) {
	u.rotation = v
	spriteAngle := u.rotation.Normalized() - (gamedata.TankFrameAngleStep / 2)
	u.sprite.FrameOffset.X = 48 * math.Trunc(float64(spriteAngle/gamedata.TankFrameAngleStep))
}

func (u *unit) SendTo(pos gmath.Vec) {
	if u.IsCommander() {
		u.sendCommanderTo(pos)
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
					u.waypoint = probe
					u.setDstRotation(u.pos.AngleToPoint(u.waypoint).Normalized())
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
		u.waypoint = posMove(aligned, d)
		u.setDstRotation(u.pos.AngleToPoint(u.waypoint).Normalized())
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
		travelled := u.calcSpeed() * delta
		u.pos = u.pos.MoveTowards(u.waypoint, travelled)
		if u.pos == u.waypoint {
			u.waypoint = gmath.Vec{}
		}

	case gamedata.UnitMovementGround:
		u.moveGroundUnitToWaypoint(delta)
	}
}
