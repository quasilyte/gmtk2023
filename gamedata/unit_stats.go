package gamedata

import (
	"github.com/quasilyte/gmtk2023/assets"
)

type UnitStats struct {
	Movement UnitMovementKind

	Body   *BodyStats
	Turret *TurretStats

	Selectable bool
	Large      bool
	Creep      bool
}

type UnitMovementKind int

const (
	UnitMovementGround UnitMovementKind = iota
	UnitMovementHover
	UnitMovementNone
)

var TankFactoryUnitStats = &UnitStats{
	Movement: UnitMovementNone,
	Body: &BodyStats{
		HP:    200,
		Image: assets.ImageTankFactory,
		Size:  unitSizeVeryBig,
	},
	Selectable: true,
	Large:      true,
}

var ConstructorUnitStats = &UnitStats{
	Movement: UnitMovementHover,
	Body: &BodyStats{
		HP:    20,
		Speed: 50,
		Image: assets.ImageDroneConstructor,
		Size:  unitSizeSmall,
	},
	Selectable: true,
}

var CommanderUnitStats = &UnitStats{
	Movement: UnitMovementHover,
	Body: &BodyStats{
		HP:    30,
		Speed: 60,
		Image: assets.ImageDroneCommander,
		Size:  unitSizeSmall,
	},
	Selectable: true,
}
