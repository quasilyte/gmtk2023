package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
)

type BodyStats struct {
	Texture resource.Image
	Image   resource.ImageID

	ProductionTime float64

	TurretOffset float64

	Speed         float64
	RotationSpeed gmath.Rad

	Size gmath.Vec

	Heavy bool

	HP float64
}

var (
	unitSizeSmall   = gmath.Vec{X: 26, Y: 26}
	unitSizeMedium  = gmath.Vec{X: 32, Y: 32}
	unitSizeBig     = gmath.Vec{X: 48, Y: 48}
	unitSizeVeryBig = gmath.Vec{X: 60, Y: 60}
)

var WheelsBodyStats = &BodyStats{
	ProductionTime: 3,
	TurretOffset:   -1,
	Size:           unitSizeSmall,
	HP:             5,
	RotationSpeed:  1.0,
	Speed:          140,
}

var MCVBodyStats = &BodyStats{
	Size:          unitSizeBig,
	HP:            250,
	RotationSpeed: 0.8,
	Speed:         30,
	Heavy:         true,
}

var ScoutBodyStats = &BodyStats{
	ProductionTime: 5,
	TurretOffset:   -1,
	Size:           unitSizeSmall,
	HP:             20,
	RotationSpeed:  1.8,
	Speed:          100,
}

var FighterBodyStats = &BodyStats{
	ProductionTime: 10,
	TurretOffset:   -1,
	Size:           unitSizeMedium,
	HP:             35,
	RotationSpeed:  2.0,
	Speed:          80,
}

var HunterBodyStats = &BodyStats{
	ProductionTime: 14,
	TurretOffset:   -1,
	Size:           unitSizeBig,
	HP:             40,
	RotationSpeed:  1.6,
	Speed:          110,
	Heavy:          true,
}

var DestroyerBodyStats = &BodyStats{
	ProductionTime: 30,
	TurretOffset:   -1,
	Size:           unitSizeBig,
	HP:             120,
	RotationSpeed:  1,
	Speed:          50,
	Heavy:          true,
}

var BunkerBodyStats = &BodyStats{
	TurretOffset: -5,
	Size:         unitSizeBig,
	HP:           100,
	Image:        assets.ImageTowerBodyBunker,
}

var CreepBunkerBodyStats = &BodyStats{
	TurretOffset: -5,
	Size:         unitSizeBig,
	HP:           100,
	Image:        assets.ImageTowerBodyCreepBunker,
}
