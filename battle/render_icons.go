package battle

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/gamedata"
)

// TODO: these functions should go away from the battle package.

func renderTowerIcon(scene *ge.Scene, icon *ebiten.Image, design *gamedata.UnitStats) {
	icon.Clear()
	icon.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 255})

	bodyTexture := scene.LoadImage(design.Body.Image)

	iconWidth := icon.Bounds().Dx()
	iconHeight := icon.Bounds().Dy()

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(
		float64((iconWidth-bodyTexture.Data.Bounds().Dx())/2),
		float64((iconHeight-bodyTexture.Data.Bounds().Dy())/2),
	)
	icon.DrawImage(bodyTexture.Data, &drawOptions)

	towerTurretTexture := design.Turret.Texture
	drawOptions.GeoM.Translate(0, design.Body.TurretOffset)
	icon.DrawImage(towerTurretTexture.Data, &drawOptions)
}

func renderFactoryIcon(scene *ge.Scene, icon *ebiten.Image, design *gamedata.UnitStats) {
	icon.Clear()
	icon.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 255})

	iconWidth := icon.Bounds().Dx()
	iconHeight := icon.Bounds().Dy()

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(4, 4)
	icon.DrawImage(scene.LoadImage(assets.ImageTankFactory).Data, &drawOptions)

	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(iconWidth-48), float64(iconHeight-48))

	tankBodyTexture := design.Body.Texture
	tankBodyFrame := tankBodyTexture.Data.SubImage(image.Rectangle{
		Min: image.Point{X: 2 * int(tankBodyTexture.DefaultFrameWidth)},
		Max: image.Point{
			X: 2*int(tankBodyTexture.DefaultFrameWidth) + int(tankBodyTexture.DefaultFrameWidth),
			Y: int(tankBodyTexture.DefaultFrameHeight),
		},
	}).(*ebiten.Image)
	icon.DrawImage(tankBodyFrame, &drawOptions)

	tankTurretTexture := design.Turret.Texture
	tankTurretFrame := tankTurretTexture.Data.SubImage(image.Rectangle{
		Min: image.Point{X: 54 * int(tankTurretTexture.DefaultFrameWidth)},
		Max: image.Point{
			X: 54*int(tankTurretTexture.DefaultFrameWidth) + int(tankTurretTexture.DefaultFrameWidth),
			Y: int(tankTurretTexture.DefaultFrameHeight),
		},
	}).(*ebiten.Image)
	icon.DrawImage(tankTurretFrame, &drawOptions)
}