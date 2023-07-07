package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

type BodyStats struct {
	Texture resource.Image

	HP float64
}

var DestroyerBodyStats = &BodyStats{
	HP: 120,
}
