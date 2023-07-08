package battle

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/pathing"
)

var groupOffsets = []gmath.Vec{
	{X: -gamedata.CellSize, Y: -gamedata.CellSize},
	{Y: -gamedata.CellSize},
	{X: +gamedata.CellSize, Y: -gamedata.CellSize},
	{X: +gamedata.CellSize},
	{X: +gamedata.CellSize, Y: +gamedata.CellSize},
	{Y: +gamedata.CellSize},
	{X: -gamedata.CellSize, Y: +gamedata.CellSize},
	{X: -gamedata.CellSize},
}

func posMove(pos gmath.Vec, d pathing.Direction) gmath.Vec {
	switch d {
	case pathing.DirRight:
		return pos.Add(gmath.Vec{X: pathing.CellSize})
	case pathing.DirDown:
		return pos.Add(gmath.Vec{Y: pathing.CellSize})
	case pathing.DirLeft:
		return pos.Add(gmath.Vec{X: -pathing.CellSize})
	case pathing.DirUp:
		return pos.Add(gmath.Vec{Y: -pathing.CellSize})
	default:
		return pos
	}
}
