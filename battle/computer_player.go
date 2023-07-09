package battle

import (
	"math"

	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/pathing"
)

type computerPlayer struct {
	world *worldState

	factories []*unit

	factoryUnits map[*unit][]*unit

	reactiveDelay float64

	tmpUnitSlice []*unit

	attackDelay float64
}

func newComputerPlayer(world *worldState) *computerPlayer {
	return &computerPlayer{
		world:        world,
		tmpUnitSlice: make([]*unit, 0, 8),
		factoryUnits: make(map[*unit][]*unit, 8),
	}
}

func (p *computerPlayer) Init() {
	p.attackDelay = p.world.Rand().FloatRange(40, 60)

	p.world.EventUnitCreated.Connect(p, func(u *unit) {
		if !u.stats.Creep {
			return
		}

		if u.IsTankFactory() {
			p.factories = append(p.factories, u)
			u.EventDestroyed.Connect(p, func(u *unit) {
				p.factories = xslices.Remove(p.factories, u)
				delete(p.factoryUnits, u)
			})
		} else if u.factory != nil {
			p.factoryUnits[u.factory] = append(p.factoryUnits[u.factory], u)
			u.EventDestroyed.Connect(p, func(u *unit) {
				p.removeUnitFromFactory(u)
			})
		}

		isBuilding := u.IsBuilding()
		defendChance := 1.0
		if !isBuilding {
			defendChance = 0.7
		}
		if defendChance == 1.0 || p.world.Rand().Chance(defendChance) {
			u.EventAttacked.Connect(p, func(attacker *unit) {
				if !isBuilding && u.turret != nil {
					isAttacking := u.turret.reload != 0
					if !isAttacking && p.world.Rand().Chance(0.6) {
						u.SendTo(u.pos.MoveTowards(attacker.pos, pathing.CellSize*2))
					}
				}

				if p.reactiveDelay > 0 {
					p.reactiveDelay = gmath.ClampMin(p.reactiveDelay-0.2, 0)
					return
				}
				maxGroupSize := p.world.Rand().IntRange(4, 8)
				g := p.findGroupNear(u.pos, maxGroupSize)
				if len(g) > 0 {
					sendGroupTo(p.world, attacker.pos.Add(gmath.RandElem(p.world.Rand(), groupOffsets)), g)
					p.detachUnitsFromFactory(g)
					p.reactiveDelay = p.world.Rand().FloatRange(10, 20)
				} else {
					p.reactiveDelay = p.world.Rand().FloatRange(2, 4)
				}
			})
		}
	})
}

func (p *computerPlayer) removeUnitFromFactory(u *unit) {
	if u.factory == nil {
		return
	}
	units := p.factoryUnits[u.factory]
	if len(units) == 0 {
		return
	}
	units = xslices.Remove(units, u)
	p.factoryUnits[u.factory] = units
}

func (p *computerPlayer) detachUnitsFromFactory(g []*unit) {
	for _, gu := range g {
		p.removeUnitFromFactory(gu)
	}
}

func (p *computerPlayer) IsDisposed() bool { return false }

func (p *computerPlayer) Update(scaledDelta, _ float64) {
	p.reactiveDelay = gmath.ClampMin(p.reactiveDelay-scaledDelta, 0)

	p.checkProactive(scaledDelta)
}

func (p *computerPlayer) checkProactive(delta float64) {
	p.attackDelay = gmath.ClampMin(p.attackDelay-delta, 0)
	if p.attackDelay == 0 {
		p.attackDelay = p.maybeDoAttack()
	}
}

func (p *computerPlayer) findGroupNear(pos gmath.Vec, maxGroupSize int) []*unit {
	p.tmpUnitSlice = p.tmpUnitSlice[:0]

	factory := p.findNearestFactory(pos, 3)
	if factory != nil {
		factoryUnits := p.factoryUnits[factory]
		if len(factoryUnits) >= maxGroupSize {
			return factoryUnits[:maxGroupSize]
		} else {
			p.tmpUnitSlice = append(p.tmpUnitSlice, factoryUnits...)
		}
	}

	return p.tmpUnitSlice
}

func (p *computerPlayer) findNearestFactory(pos gmath.Vec, minUnits int) *unit {
	if len(p.factories) == 0 {
		return nil
	}

	var closestUnit *unit
	closestDistSqr := math.MaxFloat64
	randIterate(p.world.Rand(), p.factories, func(factory *unit) bool {
		if minUnits != 0 {
			units := p.factoryUnits[factory]
			if len(units) < minUnits {
				return false
			}
		}
		distSqr := factory.pos.DistanceSquaredTo(pos)
		if distSqr < closestDistSqr {
			closestDistSqr = distSqr
			closestUnit = factory
		}
		return false
	})

	return closestUnit
}

func (p *computerPlayer) maybeDoAttack() float64 {
	if len(p.factories) != 0 {
		groupSize := p.world.Rand().IntRange(3, 6)

		factory := gmath.RandElem(p.world.Rand(), p.factories)

		// Choose the attack target.
		bestTargetScore := 0.0
		var bestTarget *unit
		for _, u := range p.world.playerUnits.selectable {
			dist := u.pos.DistanceTo(factory.pos)
			distScoreMult := 1000.0 / dist
			kindScore := 5.0
			switch {
			case u.IsTankFactory():
				kindScore = 20.0
			case u.IsMCV():
				kindScore = 18.0
			case u.IsGenerator():
				kindScore = 16.0
			case u.IsRepairDepot():
				kindScore = 12.0
			case u.IsConstructor():
				kindScore = 11.0
			case u.IsCommander():
				kindScore = 10.0
			}
			score := (kindScore * distScoreMult) * u.world.Rand().FloatRange(0.9, 1.1)
			if score > bestTargetScore {
				bestTargetScore = score
				bestTarget = u
			}
		}
		if bestTarget == nil {
			return p.world.Rand().FloatRange(10, 15)
		}

		p.tmpUnitSlice = p.tmpUnitSlice[:0]
		randIterate(p.world.Rand(), p.factoryUnits[factory], func(u *unit) bool {
			p.tmpUnitSlice = append(p.tmpUnitSlice, u)
			return len(p.tmpUnitSlice) >= groupSize
		})
		if len(p.tmpUnitSlice) < groupSize {
			return p.world.Rand().FloatRange(3, 6)
		}
		sendGroupTo(p.world, bestTarget.pos, p.tmpUnitSlice)
		groupSizeExtraDelay := float64(groupSize) * 10
		return p.world.Rand().FloatRange(15, 30) + groupSizeExtraDelay
	}

	return p.world.Rand().FloatRange(5, 10)
}
