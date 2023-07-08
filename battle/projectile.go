package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/gamedata"
)

type projectile struct {
	attacker  *unit
	target    *unit
	pos       gmath.Vec
	toPos     gmath.Vec
	fireDelay float64
	weapon    *gamedata.TurretStats
	world     *worldState

	rotation gmath.Rad

	arcProgressionScaling float64
	arcProgression        float64
	arcStart              gmath.Vec
	arcFrom               gmath.Vec
	arcTo                 gmath.Vec

	disposed bool
	sprite   *ge.Sprite
}

type projectileConfig struct {
	Weapon     *gamedata.TurretStats
	World      *worldState
	Attacker   *unit
	ToPos      gmath.Vec
	Target     *unit
	FireDelay  float64
	FireOffset gmath.Vec
}

func initProjectile(p *projectile, config projectileConfig) {
	*p = projectile{
		weapon:    config.Weapon,
		attacker:  config.Attacker,
		toPos:     config.ToPos,
		target:    config.Target,
		fireDelay: config.FireDelay,
		world:     config.World,
	}
	p.initPos()
}

func (p *projectile) Init(scene *ge.Scene) {
	if p.weapon.ArcPower != 0 {
		arcPower := p.weapon.ArcPower

		speed := p.weapon.ProjectileSpeed
		p.rotation = -math.Pi / 2
		if p.toPos.Y >= p.pos.Y {
			arcPower *= 0.3
			speed *= 1.5
		}
		dist := p.pos.DistanceTo(p.toPos)
		t := dist / speed
		p.arcProgressionScaling = 1.0 / t
		power := gmath.Vec{Y: dist * arcPower}
		p.arcFrom = p.pos.Add(power)
		p.arcTo = p.toPos.Add(power)
		p.arcStart = p.pos
	} else if p.weapon.ProjectileRotateSpeed == 0 {
		p.rotation = p.pos.AngleToPoint(p.toPos)
	} else {
		p.rotation = scene.Rand().Rad()
	}

	p.sprite = scene.NewSprite(p.weapon.ProjectileImage)
	p.sprite.Pos.Base = &p.pos
	p.sprite.Rotation = &p.rotation
	p.world.Stage().AddSpriteAbove(p.sprite)
	p.sprite.Visible = false

	if p.weapon.Accuracy != 1.0 {
		missChance := 1.0 - p.weapon.Accuracy
		if missChance != 0 && scene.Rand().Chance(missChance) {
			dist := p.pos.DistanceTo(p.toPos)
			// 100 => 25
			// 200 => 50
			// 400 => 100
			offsetValue := gmath.Clamp(dist*0.25, 24, 140)
			p.toPos = p.toPos.Add(scene.Rand().Offset(-offsetValue, offsetValue))
		} else if p.arcProgressionScaling != 0 {
			p.toPos = p.toPos.Add(scene.Rand().Offset(-8, 8))
		}
	}

	if p.fireDelay == 0 && p.weapon.ProjectilePlaysSound {
		p.playFireSound()
	}
}

func (p *projectile) IsDisposed() bool { return p.disposed }

func (p *projectile) playFireSound() {
	playSound(p.world, p.weapon.AttackSound, p.pos)
}

func (p *projectile) initPos() {
	p.pos = p.attacker.pos.Add(p.weapon.FireOffset)
	p.pos = p.pos.MoveInDirection(20, p.attacker.turret.rotation)
}

func (p *projectile) Update(delta float64) {
	if p.fireDelay > 0 {
		if p.attacker.IsDisposed() || p.attacker.turret.target != p.target {
			p.Dispose()
			return
		}
		p.fireDelay -= delta
		if p.fireDelay <= 0 {
			p.sprite.Visible = true
			p.initPos()
			p.arcStart = p.pos
			if p.weapon.ProjectilePlaysSound {
				p.playFireSound()
			}
		} else {
			return
		}
	}

	travelled := p.weapon.ProjectileSpeed * delta

	if p.arcProgressionScaling == 0 {
		if p.pos.DistanceTo(p.toPos) <= travelled {
			p.detonate()
			return
		}
		p.pos = p.pos.MoveTowards(p.toPos, travelled)
		if p.weapon.ProjectileRotateSpeed != 0 {
			p.rotation += gmath.Rad(delta * p.weapon.ProjectileRotateSpeed)
		}
		p.sprite.Visible = true
		return
	}

	p.arcProgression += delta * p.arcProgressionScaling
	if p.arcProgression >= 1 {
		p.detonate()
		return
	}
	newPos := p.arcStart.CubicInterpolate(p.arcFrom, p.toPos, p.arcTo, p.arcProgression)
	if !p.weapon.ProjectileIsRound {
		p.rotation = p.pos.AngleToPoint(newPos)
	}
	p.pos = newPos
	p.sprite.Visible = true
}

func (p *projectile) Dispose() {
	p.sprite.Dispose()
	p.disposed = true
}

func (p *projectile) createExplosion() {
	explosionPos := p.pos.Add(p.world.Rand().Offset(-4, 4))
	above := false

	switch p.weapon.ProjectileExplosion {
	case gamedata.ExplosionNormal:
		explosion := newEffectNode(p.world, explosionPos, above, assets.ImageSmallExplosion)
		p.world.runner.AddObject(explosion)
		playExplosionSound(p.world, explosionPos)
	}
}

func (p *projectile) detonate() {
	p.Dispose()

	if p.target.IsDisposed() || p.toPos.DistanceSquaredTo(p.target.pos) > p.weapon.ImpactAreaSqr {
		if p.weapon.ProjectileAlwaysExplodes {
			p.createExplosion()
		}
		return
	}

	dmg := calcDamage(p.target, p.weapon.Damage)
	p.target.OnDamage(dmg, p.attacker)
	if p.weapon.ProjectileExplosion != gamedata.ExplosionNone {
		p.createExplosion()
	}
}
