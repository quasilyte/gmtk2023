package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/pathing"
	"github.com/quasilyte/gmtk2023/viewport"
)

type worldState struct {
	Camera *viewport.Camera

	PlayerInput *input.Handler

	runner *Runner

	playerUnits playerUnits

	enemyUnits []*unit

	gridCounters map[int]uint8
	pathgrid     *pathing.Grid
	bfs          *pathing.GreedyBFS
}

type playerUnits struct {
	selectable    []*unit
	nonSelectable []*unit
	towers        []*unit
}

func newWorldState() *worldState {
	return &worldState{
		playerUnits: playerUnits{
			selectable:    make([]*unit, 0, 16),
			nonSelectable: make([]*unit, 0, 96),
			towers:        make([]*unit, 0, 10),
		},
		enemyUnits:   make([]*unit, 0, 64),
		gridCounters: make(map[int]uint8, 64),
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
	if u.IsTower() {
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
