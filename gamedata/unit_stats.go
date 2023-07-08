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

	ConstructorsNeeded int
	ConstructionTime   float64
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
	Selectable:         true,
	Large:              true,
	ConstructorsNeeded: 2,
	ConstructionTime:   20,
}

var HeavyTankFactoryUnitStats = &UnitStats{
	Movement: UnitMovementNone,
	Body: &BodyStats{
		HP:    300,
		Image: assets.ImageHeavyTankFactory,
		Size:  unitSizeVeryBig,
	},
	Selectable:         true,
	Large:              true,
	ConstructorsNeeded: 3,
	ConstructionTime:   40,
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
		Speed: 40,
		Image: assets.ImageDroneCommander,
		Size:  unitSizeSmall,
	},
	Selectable: true,
}
