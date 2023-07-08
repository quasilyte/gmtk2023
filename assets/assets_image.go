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
		ImageUITankSelector:  {Path: "image/battleui/tank_selector.png"},

		ImageTowerBodyBunker: {Path: "image/building/tower_body_bunker.png"},

		ImageTankBodyScout:     {Path: "image/tank/body_scout.png"},
		ImageTankBodyDestroyer: {Path: "image/tank/body_destroyer.png"},

		ImageTankTurretGatling:     {Path: "image/tank/turret_gatling.png"},
		ImageTankTurretLightCannon: {Path: "image/tank/turret_light_cannon.png"},

		ImageDroneCommander: {Path: "image/drone/commander.png"},
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
	ImageUITankSelector

	ImageTowerBodyBunker

	ImageTankBodyScout
	ImageTankBodyDestroyer

	ImageTankTurretGatling
	ImageTankTurretLightCannon

	ImageDroneCommander
)
