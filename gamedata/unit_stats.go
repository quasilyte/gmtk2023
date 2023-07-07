package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmtk2023/assets"
)

type UnitStats struct {
	Movement UnitMovementKind

	Body   *BodyStats
	Turret *TurretStats

	Selectable bool
	Creep      bool

	HP    float64
	Speed float64

	Image resource.ImageID
}

type UnitMovementKind int

const (
	UnitMovementGround UnitMovementKind = iota
	UnitMovementHover
)

var CommanderUnitStats = &UnitStats{
	Movement:   UnitMovementHover,
	Image:      assets.ImageDroneCommander,
	Selectable: true,
	Speed:      40,
	HP:         30,
}
