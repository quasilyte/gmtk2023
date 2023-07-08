package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
)

type BodyStats struct {
	Texture resource.Image
	Image   resource.ImageID

	TurretOffset float64

	Speed         float64
	RotationSpeed gmath.Rad

	Size gmath.Vec

	HP float64
}

var (
	unitSizeSmall  = gmath.Vec{X: 26, Y: 26}
	unitSizeMedium = gmath.Vec{X: 32, Y: 32}
	unitSizeBig    = gmath.Vec{X: 48, Y: 48}
)

var ScoutBodyStats = &BodyStats{
	TurretOffset:  -1,
	Size:          unitSizeSmall,
	HP:            20,
	RotationSpeed: 1.8,
	Speed:         100,
}

var DestroyerBodyStats = &BodyStats{
	TurretOffset:  -1,
	Size:          unitSizeBig,
	HP:            120,
	RotationSpeed: 1,
	Speed:         50,
}

var BunkerBodyStats = &BodyStats{
	TurretOffset: -4,
	Size:         unitSizeBig,
	HP:           100,
	Image:        assets.ImageTowerBodyBunker,
}
