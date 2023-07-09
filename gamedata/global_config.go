package gamedata

import (
	"math"

	"github.com/quasilyte/gmath"
)

const (
	NumTankSpriteFrames = 64
	TankFrameAngleStep  = gmath.Rad((2 * math.Pi) / float64(NumTankSpriteFrames))

	CellSize        = 40
	NumSegmentCells = 64

	MaxGroupSize = 8

	ConstructorEnergyCost = 60
	CommanderEnergyCost   = 80
)
