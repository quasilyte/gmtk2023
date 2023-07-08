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

		ImageTankTurretLightCannon: {Path: "image/tank/turret_light_cannon.png"},
		ImageTankTurretGatling:     {Path: "image/tank/turret_gatling.png"},

		ImageDroneCommander: {Path: "image/drone/commander.png"},

		ImageProjectileLightCannon: {Path: "image/projectile/projectile_light_cannon.png"},
		ImageProjectileGatling:     {Path: "image/projectile/projectile_gatling.png"},

		ImageSmallExplosion:       {Path: "image/effect/small_explosion.png", FrameWidth: 32},
		ImageBigExplosion:         {Path: "image/effect/big_explosion.png", FrameWidth: 64},
		ImageVerticalExplosion:    {Path: "image/effect/vertical_explosion.png", FrameWidth: 50},
		ImageBigVerticalExplosion: {Path: "image/effect/big_vertical_explosion.png", FrameWidth: 38},
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

	ImageTankTurretLightCannon
	ImageTankTurretGatling

	ImageDroneCommander

	ImageProjectileLightCannon
	ImageProjectileGatling

	ImageSmallExplosion
	ImageBigExplosion
	ImageVerticalExplosion
	ImageBigVerticalExplosion
)
