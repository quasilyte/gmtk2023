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

	HP float64
}

var ScoutBodyStats = &BodyStats{
	TurretOffset:  -1,
	HP:            20,
	RotationSpeed: 1.8,
	Speed:         100,
}

var DestroyerBodyStats = &BodyStats{
	TurretOffset:  -1,
	HP:            120,
	RotationSpeed: 1,
	Speed:         50,
}

var BunkerBodyStats = &BodyStats{
	TurretOffset: -4,
	HP:           100,
	Image:        assets.ImageTowerBodyBunker,
}
