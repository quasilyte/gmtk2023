package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmtk2023/assets"
)

type UnitStats struct {
	Movement UnitMovementKind

	Selectable bool
	Creep      bool

	Speed float64

	Image resource.ImageID
}

type UnitMovementKind int

const (
	UnitMovementGround UnitMovementKind = iota
	UnitMovementHover
)

func initUnitStats(stats *UnitStats) *UnitStats {
	return stats
}

var CommanderUnitStats = initUnitStats(&UnitStats{
	Movement:   UnitMovementHover,
	Image:      assets.ImageDroneCommander,
	Selectable: true,
	Speed:      40,
})
