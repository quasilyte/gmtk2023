package bootstrap

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/session"
)

func InitState(ctx *ge.Context, state *session.State) {
	type textureConfig struct {
		src        resource.ImageID
		depth      int
		colorLayer int
		dst        *resource.Image
	}

	tankTextureTasks := []textureConfig{
		{assets.ImageTankBodyMCV, 3, 1, &gamedata.MCVBodyStats.Texture},
		{assets.ImageTankBodyScout, 3, 1, &gamedata.ScoutBodyStats.Texture},
		{assets.ImageTankBodyFighter, 3, 1, &gamedata.FighterBodyStats.Texture},
		{assets.ImageTankBodyHunter, 3, 1, &gamedata.HunterBodyStats.Texture},
		{assets.ImageTankBodyDestroyer, 3, 1, &gamedata.DestroyerBodyStats.Texture},
		{assets.ImageTankBodyWheels, 3, 1, &gamedata.WheelsBodyStats.Texture},

		{assets.ImageTankTurretScatterCannon, 2, 2, &gamedata.ScatterCannonStats.Texture},
		{assets.ImageTankTurretLightCannon, 2, 2, &gamedata.LightCannonStats.Texture},
		{assets.ImageTankTurretGatling, 2, 2, &gamedata.GatlingStats.Texture},
	}

	s := ge.NewSprite(ctx)
	for _, task := range tankTextureTasks {
		s.SetImage(ctx.Loader.LoadImage(task.src))
		tex, frameHeight, frameWidth := createTexture(s, task.depth, task.colorLayer)
		*task.dst = resource.Image{
			Data:               tex,
			DefaultFrameWidth:  frameWidth,
			DefaultFrameHeight: frameHeight,
		}
	}
}

func createTexture(source *ge.Sprite, depth, colorLayer int) (*ebiten.Image, float64, float64) {
	source.Centered = true
	pos := gmath.Vec{
		X: source.ImageWidth() / 2,
		Y: source.ImageHeight() / 2,
	}
	source.Pos.Base = &pos
	angle := gmath.Rad(0)
	source.Rotation = &angle

	sides := gamedata.NumTankSpriteFrames
	width := int(source.ImageWidth()) + depth + 4
	height := int(source.ImageHeight()) + depth + 4

	result := ebiten.NewImage(width*sides, height+depth)
	tmpImage := ebiten.NewImage(width, height)

	offsetX := 0
	for i := 0; i < sides; i++ {
		tmpImage.Clear()
		source.Draw(tmpImage)
		addShading(tmpImage, depth, colorLayer)

		var options ebiten.DrawImageOptions
		options.GeoM.Translate(float64(offsetX)+1, float64(depth))
		result.DrawImage(tmpImage, &options)

		offsetX += width
		angle += gamedata.TankFrameAngleStep
	}
	return result, float64(width), float64(height)
}

func addShading(img *ebiten.Image, depth, colorLayer int) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	numPixels := 4 * width * height
	origPixels := make([]byte, numPixels)
	pixels := make([]byte, numPixels)
	img.ReadPixels(origPixels)
	copy(pixels, origPixels)

	getColor := func(x, y int) (color.RGBA, bool) {
		i := 4 * (y*width + x)
		if i > len(origPixels) || i < 0 {
			return color.RGBA{}, false
		}
		clr := color.RGBA{
			R: origPixels[i+0],
			G: origPixels[i+1],
			B: origPixels[i+2],
			A: origPixels[i+3],
		}
		return clr, true
	}
	setColor := func(x, y int, clr color.RGBA) {
		i := 4 * (y*width + x)
		pixels[i+0] = clr.R
		pixels[i+1] = clr.G
		pixels[i+2] = clr.B
		pixels[i+3] = clr.A
	}

	upperOutline := make([]uint8, width)
	for i := range upperOutline {
		upperOutline[i] = 0xff
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			clr, ok := getColor(x, y)
			if !ok || clr.A == 0 {
				continue
			}
			if y < int(upperOutline[x]) {
				upperOutline[x] = uint8(y)
			}
			for i := 1; i < depth+1; i++ {
				pix, ok := getColor(x, y+i)
				if ok && pix.A == 0 {
					shadingColor := clr
					shadingColor.A = 255
					switch colorLayer {
					case 0:
						shadingColor.R /= uint8(i) + 1
						shadingColor.G /= uint8(i) + 1
						shadingColor.B /= uint8(i) + 1
					default:
						// 1 -> 20
						// 2 -> 10
						decrease := uint8(i) * uint8(30-10*colorLayer)
						shadingColor = rgbaSub(shadingColor, decrease)
					}
					setColor(x, y+i, shadingColor)
				}
			}
		}
	}

	for x, y := range upperOutline {
		if y == 0xff {
			continue
		}
		clr, ok := getColor(x, int(y))
		if !ok {
			panic("unreachable")
		}
		setColor(x, int(y), rgbaSub(clr, 20))
		clr2, ok := getColor(x, int(y)+1)
		if !ok && clr.A == 0 {
			continue
		}
		setColor(x, int(y)+1, rgbaSub(clr2, 15))
	}

	img.WritePixels(pixels)
}
