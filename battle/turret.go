package battle

import (
	"math"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/gamedata"
)

type turret struct {
	world *worldState

	target *unit
	owner  *unit

	sprite     *ge.Sprite
	frameWidth float64

	rotation    gmath.Rad
	dstRotation gmath.Rad

	dangerCheckDelay  float64
	dangerModeTime    float64
	seekTargetDelay   float64
	targetChangeDelay float64
	fireAttemptDelay  float64
	reload            float64

	stats *gamedata.TurretStats

	needRotate bool

	config turretConfig
}

type turretConfig struct {
	Image resource.Image

	Pos *gmath.Vec

	InitialRotation gmath.Rad

	Owner *unit
}

func newTurret(world *worldState, config turretConfig) *turret {
	return &turret{
		world:      world,
		owner:      config.Owner,
		config:     config,
		rotation:   config.InitialRotation,
		stats:      config.Owner.stats.Turret,
		frameWidth: config.Image.DefaultFrameWidth,
	}
}

func (t *turret) Init(scene *ge.Scene) {
	t.sprite = ge.NewSprite(scene.Context())
	t.sprite.SetImage(t.config.Image)
	t.sprite.Pos.Base = t.config.Pos
	t.world.Stage().AddSprite(t.sprite)

	if t.owner.stats.Creep {
		t.sprite.SetHue(gmath.DegToRad(80))
	}

	t.setRotation(t.rotation)
}

func (t *turret) IsDisposed() bool {
	return t.sprite.IsDisposed()
}

func (t *turret) Dispose() {
	t.sprite.Dispose()
}

func (t *turret) Update(delta float64) {
	if t.needRotate {
		t.setRotation(t.rotation.RotatedTowards(t.dstRotation, t.stats.RotationSpeed*gmath.Rad(delta)))
		if t.rotation == t.dstRotation {
			t.needRotate = false
		}
	}

	t.reload = gmath.ClampMin(t.reload-delta, 0)

	if t.dangerModeTime > 0 {
		t.updateInDangerMode(delta)
	} else {
		t.updateInNormalMode(delta)
	}
}

func (t *turret) updateInNormalMode(delta float64) {
	t.dangerCheckDelay = gmath.ClampMin(t.dangerCheckDelay-delta, 0)

	if t.dangerCheckDelay == 0 {
		t.dangerCheckDelay = t.world.Rand().FloatRange(3.5, 8.5)
		t.maybeEnterDangerMode()
	}
}

func (t *turret) maybeChangeTarget() {
	// Is target still in lock-on range?
	if t.target.pos.DistanceSquaredTo(*t.config.Pos) > t.stats.TargetLockSqr {
		t.target = nil
		return
	}

	if t.world.Rand().Chance(0.1) && t.reload == 0 {
		t.target = nil
		return
	}
}

func (t *turret) updateInDangerMode(delta float64) {
	if t.target != nil && t.target.IsDisposed() {
		t.target = nil
	}

	t.fireAttemptDelay = gmath.ClampMin(t.fireAttemptDelay-delta, 0)

	t.targetChangeDelay = gmath.ClampMin(t.targetChangeDelay-delta, 0)
	if t.targetChangeDelay == 0 {
		t.targetChangeDelay = t.world.Rand().FloatRange(3.0, 6.5)
		if t.target != nil {
			t.maybeChangeTarget()
		}
	}

	t.seekTargetDelay = gmath.ClampMin(t.seekTargetDelay-delta, 0)
	if t.target == nil && t.seekTargetDelay == 0 {
		t.seekTargetDelay = t.world.Rand().FloatRange(1.0, 3.5)
		t.target = t.seekTarget()
	}

	if t.dangerModeTime > 0 && t.target == nil {
		t.dangerModeTime -= delta
		if t.dangerModeTime <= 0 {
			t.dangerModeTime = 0
			t.dangerCheckDelay = t.world.Rand().FloatRange(1.5, 4.5)
			return
		}
	}

	if t.target == nil || t.reload > 0 || t.fireAttemptDelay > 0 {
		return
	}
	if t.target.pos.DistanceSquaredTo(*t.config.Pos) > t.stats.RangeSqr {
		t.fireAttemptDelay = 0.5
		return
	}
	t.fireAttemptDelay = 0.15
	angleToTarget := t.owner.pos.AngleToPoint(t.target.pos).Normalized()
	canFire := t.rotation.AngleDelta2(angleToTarget) <= t.stats.MaxAngleDelta
	if !canFire {
		t.setDstRotation(angleToTarget)
		return
	}

	t.reload = t.stats.Reload * t.world.Rand().FloatRange(0.9, 1.1)
	if !t.stats.ProjectilePlaysSound {
		playSound(t.world, t.stats.AttackSound, t.owner.pos)
	}
	for i := 0; i < t.stats.BurstSize; i++ {
		fireDelay := float64(i) * t.stats.BurstDelay
		t.world.WalkNearbyTargets(t.owner.stats, t.target, func(target *unit) {
			p := t.world.NewProjectile(projectileConfig{
				Weapon:    t.stats,
				World:     t.world,
				Attacker:  t.owner,
				ToPos:     target.pos,
				Target:    target,
				FireDelay: fireDelay,
			})
			t.world.runner.AddProjectile(p)
		})

	}
}

func (t *turret) seekTarget() *unit {
	return t.world.FindTarget(*t.config.Pos, t.owner.stats)
}

func (t *turret) maybeEnterDangerMode() {
	if t.world.HasEnemyNearby(*t.config.Pos, t.owner.stats) {
		t.dangerModeTime = t.world.Rand().FloatRange(7, 15)
		t.seekTargetDelay = 0
	}
}

func (t *turret) setRotation(v gmath.Rad) {
	t.rotation = v
	spriteAngle := t.rotation.Normalized() - (gamedata.TankFrameAngleStep / 2)
	t.sprite.FrameOffset.X = t.frameWidth * math.Trunc(float64(spriteAngle/gamedata.TankFrameAngleStep))
}

func (t *turret) setDstRotation(v gmath.Rad) {
	t.needRotate = true
	t.dstRotation = v
}

func (t *turret) OnAttacked(attacker *unit) {
	if t.target == nil {
		t.target = attacker
		t.dangerModeTime = 8
		return
	}
}

func (t *turret) AlignRequest(rotation gmath.Rad) {
	if rotation == t.rotation {
		return
	}
	if t.dangerModeTime > 0 {
		return
	}

	t.setDstRotation(rotation)
}
