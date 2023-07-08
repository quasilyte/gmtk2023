package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
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

func randIterate[T any](rand *gmath.Rand, slice []T, f func(x T) bool) T {
	var result T
	if len(slice) == 0 {
		return result
	}
	if len(slice) == 1 {
		// Don't use rand() if there is only 1 element.
		x := slice[0]
		if f(x) {
			result = x
		}
		return result
	}

	var slider gmath.Slider
	slider.SetBounds(0, len(slice)-1)
	slider.TrySetValue(rand.IntRange(0, len(slice)-1))
	inc := rand.Bool()
	for i := 0; i < len(slice); i++ {
		x := slice[slider.Value()]
		if inc {
			slider.Inc()
		} else {
			slider.Dec()
		}
		if f(x) {
			result = x
			break
		}
	}
	return result
}

func playSound(world *worldState, id resource.AudioID, pos gmath.Vec) {
	if world.Camera.ContainsPos(pos) {
		numSamples := assets.NumSamples(id)
		if numSamples == 1 {
			world.scene.Audio().PlaySound(id)
		} else {
			soundIndex := world.Rand().IntRange(0, numSamples-1)
			sound := resource.AudioID(int(id) + soundIndex)
			world.scene.Audio().PlaySound(sound)
		}
	}
}

func playExplosionSound(world *worldState, pos gmath.Vec) {
	playSound(world, assets.AudioExplosion1, pos)
}

func spriteRect(pos gmath.Vec, width, height float64) gmath.Rect {
	offset := gmath.Vec{X: width * 0.5, Y: height * 0.5}
	return gmath.Rect{
		Min: pos.Sub(offset),
		Max: pos.Add(offset),
	}
}

func createAreaExplosion(world *worldState, rect gmath.Rect, allowVertical bool) {
	// FIXME: Rect.Center() does not work properly in gmath.
	center := gmath.Vec{
		X: rect.Max.X - rect.Width()*0.5,
		Y: rect.Max.Y - rect.Height()*0.5,
	}
	size := rect.Width() * rect.Height()
	minExplosions := gmath.ClampMin(size/120.0, 1)
	numExplosions := world.Rand().IntRange(int(minExplosions), int(minExplosions*1.3))
	above := !allowVertical
	for numExplosions > 0 {
		offset := gmath.Vec{
			X: world.Rand().FloatRange(-rect.Width()*0.4, rect.Width()*0.4),
			Y: world.Rand().FloatRange(-rect.Height()*0.4, rect.Height()*0.4),
		}
		if numExplosions >= 4 && world.Rand().Chance(0.4) {
			numExplosions -= 4
			world.runner.AddObject(newEffectNode(world, center.Add(offset), above, assets.ImageBigExplosion))
		} else {
			numExplosions--
			if allowVertical && world.Rand().Chance(0.4) {
				effect := newEffectNode(world, center.Add(offset), above, assets.ImageVerticalExplosion)
				world.runner.AddObject(effect)
				effect.anim.SetSecondsPerFrame(0.035)
			} else {
				effect := newEffectNode(world, center.Add(offset), above, assets.ImageSmallExplosion)
				world.runner.AddObject(effect)
			}
		}
	}
	playExplosionSound(world, center)
}
