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

var TowerConstruction = &UnitStats{
	Movement: UnitMovementNone,
	Body: &BodyStats{
		HP:    30,
		Image: assets.ImageTowerBodyBunker,
		Size:  unitSizeBig,
	},
	Selectable:         true,
	ConstructorsNeeded: 1,
}

var GeneratorUnitStats = &UnitStats{
	Movement: UnitMovementNone,
	Body: &BodyStats{
		HP:    70,
		Image: assets.ImageGenerator,
		Size:  unitSizeMedium,
	},
	Selectable:         true,
	ConstructorsNeeded: 1,
	ConstructionTime:   15,
}

var RepairDepotUnitStats = &UnitStats{
	Movement: UnitMovementNone,
	Body: &BodyStats{
		HP:    100,
		Image: assets.ImageRepairDepot,
		Size:  unitSizeMedium,
	},
	Selectable:         true,
	ConstructorsNeeded: 2,
	ConstructionTime:   25,
}

var CreepTankFactoryUnitStats = &UnitStats{
	Movement: UnitMovementNone,
	Body: &BodyStats{
		HP:    175,
		Image: assets.ImageCreepTankFactory,
		Size:  unitSizeVeryBig,
	},
	Creep: true,
}

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
	ConstructionTime:   30,
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
	ConstructionTime:   60,
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
