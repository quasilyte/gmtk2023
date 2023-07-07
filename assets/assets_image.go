package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerImageResources(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageBackgroundTiles: {Path: "image/tiles.png"},

		ImageUIDroneSelector: {Path: "image/battleui/drone_selector.png"},

		ImageTankBodyDestroyer: {Path: "image/tank/body_destroyer.png"},

		ImageTankTurretLightCannon: {Path: "image/tank/turret_light_cannon.png"},

		ImageDroneCommander: {Path: "image/drone/commander.png", FrameWidth: 11},
	}

	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageBackgroundTiles

	ImageUIDroneSelector

	ImageTankBodyDestroyer

	ImageTankTurretLightCannon

	ImageDroneCommander
)
