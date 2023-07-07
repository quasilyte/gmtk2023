package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmtk2023/assets"
)

type BodyStats struct {
	Texture resource.Image
	Image   resource.ImageID

	TurretOffset float64

	Speed float64

	HP float64
}

var DestroyerBodyStats = &BodyStats{
	TurretOffset: -1,
	HP:           120,
	Speed:        50,
}

var BunkerBodyStats = &BodyStats{
	TurretOffset: -4,
	HP:           100,
	Image:        assets.ImageTowerBodyBunker,
}
