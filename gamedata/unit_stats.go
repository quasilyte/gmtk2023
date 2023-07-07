package gamedata

import (
	"github.com/quasilyte/gmtk2023/assets"
)

type UnitStats struct {
	Movement UnitMovementKind

	Body   *BodyStats
	Turret *TurretStats

	Selectable bool
	Creep      bool
}

type UnitMovementKind int

const (
	UnitMovementGround UnitMovementKind = iota
	UnitMovementHover
	UnitMovementNone
)

var CommanderUnitStats = &UnitStats{
	Movement: UnitMovementHover,
	Body: &BodyStats{
		HP:    30,
		Speed: 60,
		Image: assets.ImageDroneCommander,
	},
	Selectable: true,
}
