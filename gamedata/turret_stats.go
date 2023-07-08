package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
)

type ProjectileExplosionKind int

const (
	ExplosionNone ProjectileExplosionKind = iota
	ExplosionNormal
)

type TurretStats struct {
	Texture resource.Image

	ProjectileImage          resource.ImageID
	ProjectileSpeed          float64
	ProjectileRotateSpeed    float64
	ProjectilePlaysSound     bool
	ProjectileIsRound        bool
	ProjectileAlwaysExplodes bool
	ImpactArea               float64
	ImpactAreaSqr            float64
	ProjectileExplosion      ProjectileExplosionKind
	AttackSound              resource.AudioID

	BurstSize  int
	BurstDelay float64

	Reload float64

	HP float64

	RotationSpeed gmath.Rad

	Range         float64
	RangeSqr      float64
	TargetLockSqr float64

	Damage DamageValue

	FireOffset gmath.Vec
	ArcPower   float64
	Accuracy   float64
}

type DamageValue struct {
	Health float64
}

func FinalizeTurretStats(stats *TurretStats) *TurretStats {
	stats.ImpactAreaSqr = stats.ImpactArea * stats.ImpactArea
	stats.RangeSqr = stats.Range * stats.Range
	stats.TargetLockSqr = stats.RangeSqr + (CellSize * CellSize)
	return stats
}

var LightCannonStats = FinalizeTurretStats(&TurretStats{
	AttackSound:         assets.AudioShotLightCannon,
	HP:                  10,
	RotationSpeed:       2.0,
	Range:               9 * CellSize,
	Accuracy:            0.8,
	ImpactArea:          8,
	ProjectileImage:     assets.ImageProjectileGatling,
	ProjectileExplosion: ExplosionNormal,
	ProjectileSpeed:     250,
	Reload:              1.5,
	BurstSize:           1,
	Damage:              DamageValue{Health: 6},
})

var GatlingStats = FinalizeTurretStats(&TurretStats{
	AttackSound:     assets.AudioShotGatling,
	HP:              0,
	RotationSpeed:   1.2,
	Range:           7 * CellSize,
	Accuracy:        0.5,
	ImpactArea:      4,
	ProjectileImage: assets.ImageProjectileGatling,
	ProjectileSpeed: 320,
	Reload:          2.0,
	BurstSize:       3,
	BurstDelay:      0.1,
	Damage:          DamageValue{Health: 1},
})
