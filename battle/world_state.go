package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/pathing"
	"github.com/quasilyte/gmtk2023/viewport"
	"github.com/quasilyte/gsignal"
)

type worldState struct {
	Camera *viewport.Camera

	PlayerInput *input.Handler

	scene *ge.Scene

	runner *Runner

	playerUnits playerUnits

	enemyUnits []*unit

	projectilePool []*projectile

	gridCounters map[int]uint8
	pathgrid     *pathing.Grid
	bfs          *pathing.GreedyBFS

	EventUnitCreated gsignal.Event[*unit]
}

type playerUnits struct {
	selectable    []*unit
	nonSelectable []*unit
	towers        []*unit
}

func newWorldState(scene *ge.Scene) *worldState {
	return &worldState{
		scene: scene,
		playerUnits: playerUnits{
			selectable:    make([]*unit, 0, 16),
			nonSelectable: make([]*unit, 0, 96),
			towers:        make([]*unit, 0, 10),
		},
		enemyUnits:     make([]*unit, 0, 64),
		gridCounters:   make(map[int]uint8, 64),
		projectilePool: make([]*projectile, 0, 128),
	}
}

func (w *worldState) Rand() *gmath.Rand {
	return w.Scene().Rand()
}

func (w *worldState) Scene() *ge.Scene {
	return w.runner.scene
}

func (w *worldState) Stage() *viewport.Stage {
	return w.Camera.Stage
}

func (w *worldState) WalkNearbyTargets(stats *gamedata.UnitStats, target *unit, f func(*unit)) {
	num := stats.Turret.MaxTargets
	f(target)
	if num == 1 {
		return
	}
	num--
	if !stats.Creep {
		w.walkNearbyTargetsInSlice(target, num, w.enemyUnits, f)
		return
	}
	num = w.walkNearbyTargetsInSlice(target, num, w.playerUnits.towers, f)
	num = w.walkNearbyTargetsInSlice(target, num, w.playerUnits.nonSelectable, f)
	w.walkNearbyTargetsInSlice(target, num, w.playerUnits.selectable, f)
}

func (w *worldState) walkNearbyTargetsInSlice(target *unit, num int, slice []*unit, f func(*unit)) int {
	if num == 0 {
		return 0
	}
	const nearbyDistSqr = 1.25 * (gamedata.CellSize * gamedata.CellSize)
	randIterate(w.Rand(), slice, func(u *unit) bool {
		if u == target {
			return false
		}
		if u.pos.DistanceSquaredTo(target.pos) > nearbyDistSqr {
			return false
		}
		f(u)
		num--
		return num <= 0
	})
	return num
}

func (w *worldState) FindTarget(pos gmath.Vec, stats *gamedata.UnitStats) *unit {
	if stats.Creep {
		if target := w.findTargetInSlice(pos, stats, w.playerUnits.towers); target != nil {
			return target
		}
		if target := w.findTargetInSlice(pos, stats, w.playerUnits.nonSelectable); target != nil {
			return target
		}
		if target := w.findTargetInSlice(pos, stats, w.playerUnits.selectable); target != nil {
			return target
		}
		return nil
	}
	return w.findTargetInSlice(pos, stats, w.enemyUnits)
}

func (w *worldState) HasEnemyNearby(pos gmath.Vec, stats *gamedata.UnitStats) bool {
	if stats.Creep {
		return w.hasEnemyInSlice(pos, stats, w.playerUnits.towers) ||
			w.hasEnemyInSlice(pos, stats, w.playerUnits.nonSelectable) ||
			w.hasEnemyInSlice(pos, stats, w.playerUnits.selectable)
	}
	return w.hasEnemyInSlice(pos, stats, w.enemyUnits)
}

func (w *worldState) findTargetInSlice(pos gmath.Vec, stats *gamedata.UnitStats, slice []*unit) *unit {
	attackDistSqr := stats.Turret.TargetLockSqr
	return randIterate(w.Rand(), slice, func(u2 *unit) bool {
		return u2.pos.DistanceSquaredTo(pos) < attackDistSqr
	})
}

func (w *worldState) hasEnemyInSlice(pos gmath.Vec, stats *gamedata.UnitStats, slice []*unit) bool {
	dangerDistSqr := (stats.Turret.RangeSqr * 1.25) + (gamedata.CellSize * gamedata.CellSize)
	for _, u2 := range slice {
		if u2.pos.DistanceSquaredTo(pos) <= dangerDistSqr {
			return true
		}
	}
	return false
}

func (w *worldState) FindSelectable(pos gmath.Vec) *unit {
	if len(w.playerUnits.selectable) == 0 {
		return nil
	}
	minDistSqr := math.MaxFloat64
	var closestUnit *unit
	for _, u := range w.playerUnits.selectable {
		distSqr := u.pos.DistanceSquaredTo(pos)
		if distSqr > (24 * 24) {
			continue
		}
		if distSqr < minDistSqr {
			minDistSqr = distSqr
			closestUnit = u
		}
	}
	return closestUnit
}

func (w *worldState) MayBlockFactory(pos gmath.Vec) bool {
	posAbove := w.pathgrid.AlignPos(pos).Sub(gmath.Vec{Y: pathing.CellSize - 1})
	const checkDist = pathing.CellSize + 1
	for _, u := range w.playerUnits.selectable {
		check := false
		switch extra := u.extra.(type) {
		case *tankFactoryExtra:
			check = true
		case *constructionSiteExtra:
			if _, ok := extra.newUnitExtra.(*tankFactoryExtra); ok {
				check = true
			}
		}
		if check {
			if u.pos.DistanceSquaredTo(posAbove) < (checkDist * checkDist) {
				return true
			}
		}
	}
	return false
}

func (w *worldState) IsInnerPos(pos gmath.Vec) bool {
	aligned := w.pathgrid.AlignPos(pos)
	if aligned.X < pathing.CellSize || aligned.Y < pathing.CellSize {
		return false
	}
	rect := w.Camera.WorldRect
	if aligned.X > rect.Width()-pathing.CellSize || aligned.Y > rect.Height()-(2*pathing.CellSize) {
		return false
	}
	return true
}

func (w *worldState) FindConstructionSitePos(pos gmath.Vec) gmath.Vec {
	alignedPos := w.pathgrid.AlignPos(pos)
	if w.pathgrid.CellIsFree(w.pathgrid.PosToCoord(alignedPos)) {
		return alignedPos
	}
	offset := randIterate(w.Rand(), groupOffsets, func(offset gmath.Vec) bool {
		probe := alignedPos.Add(offset)
		if w.pathgrid.CellIsFree(w.pathgrid.PosToCoord(probe)) {
			return true
		}
		return false
	})
	if !offset.IsZero() {
		return alignedPos.Add(offset)
	}
	return gmath.Vec{}
}

func (w *worldState) FindConstructor(pos gmath.Vec) *unit {
	if len(w.playerUnits.selectable) == 0 {
		return nil
	}
	minDistSqr := math.MaxFloat64
	var closestUnit *unit
	for _, u := range w.playerUnits.selectable {
		distSqr := u.pos.DistanceSquaredTo(pos)
		if distSqr > (24 * 24) {
			continue
		}
		if distSqr < minDistSqr {
			minDistSqr = distSqr
			closestUnit = u
		}
	}
	return closestUnit
}

func (w *worldState) FindAssignable(pos gmath.Vec) *unit {
	if len(w.playerUnits.nonSelectable) == 0 {
		return nil
	}
	minDistSqr := math.MaxFloat64
	var closestUnit *unit
	for _, u := range w.playerUnits.nonSelectable {
		distSqr := u.pos.DistanceSquaredTo(pos)
		if distSqr > (24 * 24) {
			continue
		}
		if distSqr < minDistSqr {
			minDistSqr = distSqr
			closestUnit = u
		}
	}
	return closestUnit
}

func (w *worldState) findUnitSlice(u *unit) *[]*unit {
	if u.stats.Creep {
		return &w.enemyUnits
	}
	if u.stats.Selectable {
		return &w.playerUnits.selectable
	}
	if u.IsBuilding() {
		return &w.playerUnits.towers
	}
	return &w.playerUnits.nonSelectable
}

func (w *worldState) NewUnit(config unitConfig) *unit {
	u := newUnit(w, config)
	slice := w.findUnitSlice(u)
	*slice = append(*slice, u)
	u.EventDestroyed.Connect(nil, func(u *unit) {
		*slice = xslices.Remove(*slice, u)
	})
	w.EventUnitCreated.Emit(u)
	return u
}

func (w *worldState) BuildPath(from, to gmath.Vec) pathing.BuildPathResult {
	return w.bfs.BuildPath(w.pathgrid, w.pathgrid.PosToCoord(from), w.pathgrid.PosToCoord(to))
}

func (w *worldState) MarkPos(pos gmath.Vec) {
	w.MarkCell(w.pathgrid.PosToCoord(pos))
}

func (w *worldState) UnmarkPos(pos gmath.Vec) {
	w.UnmarkCell(w.pathgrid.PosToCoord(pos))
}

func (w *worldState) MarkCell(coord pathing.GridCoord) {
	key := w.pathgrid.CoordToIndex(coord)
	if v := w.gridCounters[key]; v == 0 {
		w.pathgrid.MarkCell(coord)
	}
	w.gridCounters[key]++
}

func (w *worldState) UnmarkCell(coord pathing.GridCoord) {
	key := w.pathgrid.CoordToIndex(coord)
	if v := w.gridCounters[key]; v == 1 {
		w.pathgrid.UnmarkCell(coord)
		delete(w.gridCounters, key)
	} else {
		w.gridCounters[key]--
	}
}

func (w *worldState) NewProjectile(config projectileConfig) *projectile {
	if len(w.projectilePool) != 0 {
		p := w.projectilePool[len(w.projectilePool)-1]
		initProjectile(p, config)
		w.projectilePool = w.projectilePool[:len(w.projectilePool)-1]
		return p
	}
	p := &projectile{}
	initProjectile(p, config)
	return p
}

func (w *worldState) FreeProjectileNode(p *projectile) {
	w.projectilePool = append(w.projectilePool, p)
}
