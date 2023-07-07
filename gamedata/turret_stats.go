package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

type TurretStats struct {
	Texture resource.Image

	HP float64
}

var LightCannonStats = &TurretStats{
	HP: 10,
}
