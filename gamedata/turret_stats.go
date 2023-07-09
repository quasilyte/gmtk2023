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

	ProductionTime float64

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
	MaxTargets int

	Reload float64

	HP float64

	RotationSpeed gmath.Rad

	Range         float64
	RangeSqr      float64
	TargetLockSqr float64

	Damage DamageValue

	FireOffset    gmath.Vec
	ArcPower      float64
	Accuracy      float64
	MaxAngleDelta gmath.Rad
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

var ScatterCannonStats = FinalizeTurretStats(&TurretStats{
	AttackSound:         assets.AudioShotLightCannon1,
	HP:                  25,
	RotationSpeed:       1.8,
	Range:               8 * CellSize,
	Accuracy:            0.75,
	MaxAngleDelta:       0.075,
	ImpactArea:          10,
	ProjectileImage:     assets.ImageProjectileLightCannon,
	ProjectileExplosion: ExplosionNormal,
	ProjectileSpeed:     250,
	Reload:              2.6,
	MaxTargets:          4,
	BurstSize:           1,
	Damage:              DamageValue{Health: 4},
})

var LightCannonStats = FinalizeTurretStats(&TurretStats{
	ProductionTime:      5,
	AttackSound:         assets.AudioShotLightCannon1,
	HP:                  10,
	RotationSpeed:       2.0,
	Range:               9 * CellSize,
	Accuracy:            0.8,
	MaxAngleDelta:       0.05,
	ImpactArea:          8,
	ProjectileImage:     assets.ImageProjectileLightCannon,
	ProjectileExplosion: ExplosionNormal,
	ProjectileSpeed:     250,
	Reload:              1.5,
	MaxTargets:          1,
	BurstSize:           1,
	Damage:              DamageValue{Health: 6},
})

var HurricaneStats = FinalizeTurretStats(&TurretStats{
	ProductionTime:      20,
	AttackSound:         assets.AudioShotHurricane1,
	HP:                  0,
	RotationSpeed:       2.0,
	Range:               15 * CellSize,
	Accuracy:            0.65,
	MaxAngleDelta:       0.15,
	ImpactArea:          14,
	ProjectileImage:     assets.ImageProjectileSmallMissile,
	ProjectileExplosion: ExplosionNormal,
	ProjectileSpeed:     240,
	Reload:              4.25,
	MaxTargets:          5,
	BurstSize:           1,
	ArcPower:            2.5,
	Damage:              DamageValue{Health: 6},
})

var GatlingStats = FinalizeTurretStats(&TurretStats{
	ProductionTime:  5,
	AttackSound:     assets.AudioShotGatling,
	HP:              10,
	RotationSpeed:   1.2,
	Range:           7 * CellSize,
	Accuracy:        0.5,
	MaxAngleDelta:   0.1,
	ImpactArea:      4,
	ProjectileImage: assets.ImageProjectileGatling,
	ProjectileSpeed: 400,
	Reload:          2.0,
	MaxTargets:      1,
	BurstSize:       3,
	BurstDelay:      0.1,
	Damage:          DamageValue{Health: 1},
})
