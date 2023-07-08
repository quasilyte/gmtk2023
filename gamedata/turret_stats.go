package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
)

type TurretStats struct {
	Texture resource.Image

	HP float64

	RotationSpeed gmath.Rad
}

var LightCannonStats = &TurretStats{
	HP:            10,
	RotationSpeed: 2.0,
}

var GatlingStats = &TurretStats{
	HP:            0,
	RotationSpeed: 1.2,
}
