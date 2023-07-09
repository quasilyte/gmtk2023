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
		ImageUITowerSelector: {Path: "image/battleui/tower_selector.png"},
		ImageUILargeSelector: {Path: "image/battleui/large_selector.png"},
		ImageUITankSelector:  {Path: "image/battleui/tank_selector.png"},

		ImageUIDeconstuctIcon:      {Path: "image/battleui/deconstruct.png"},
		ImageUIConstructors1outof1: {Path: "image/battleui/constructors_1_outof_1.png"},
		ImageUIConstructors1outof2: {Path: "image/battleui/constructors_1_outof_2.png"},
		ImageUIConstructors2outof2: {Path: "image/battleui/constructors_2_outof_2.png"},
		ImageUIConstructors1outof3: {Path: "image/battleui/constructors_1_outof_3.png"},
		ImageUIConstructors2outof3: {Path: "image/battleui/constructors_2_outof_3.png"},
		ImageUIConstructors3outof3: {Path: "image/battleui/constructors_3_outof_3.png"},

		ImageTowerBodyBunker:      {Path: "image/building/tower_body_bunker.png"},
		ImageTowerBodyCreepBunker: {Path: "image/building/tower_body_creep_bunker.png"},

		ImageTankBodyMCV:       {Path: "image/tank/body_mcv.png"},
		ImageTankBodyScout:     {Path: "image/tank/body_scout.png"},
		ImageTankBodyFighter:   {Path: "image/tank/body_fighter.png"},
		ImageTankBodyHunter:    {Path: "image/tank/body_hunter.png"},
		ImageTankBodyDestroyer: {Path: "image/tank/body_destroyer.png"},
		ImageTankBodyWheels:    {Path: "image/tank/body_wheels.png"},

		ImageTankTurretScatterCannon: {Path: "image/tank/turret_scatter_cannon.png"},
		ImageTankTurretLightCannon:   {Path: "image/tank/turret_light_cannon.png"},
		ImageTankTurretHurricane:     {Path: "image/tank/turret_hurricane.png"},
		ImageTankTurretGatling:       {Path: "image/tank/turret_gatling.png"},

		ImageDroneConstructor: {Path: "image/drone/constructor.png"},
		ImageDroneCommander:   {Path: "image/drone/commander.png"},

		ImageGenerator:        {Path: "image/building/generator.png"},
		ImageRepairDepot:      {Path: "image/building/repair_depot.png"},
		ImageTankFactory:      {Path: "image/building/tank_factory.png", FrameWidth: 46},
		ImageHeavyTankFactory: {Path: "image/building/heavy_tank_factory.png", FrameWidth: 52},
		ImageCreepTankFactory: {Path: "image/building/creep_tank_factory.png", FrameWidth: 52},

		ImageProjectileSmallMissile: {Path: "image/projectile/projectile_small_missile.png"},
		ImageProjectileLightCannon:  {Path: "image/projectile/projectile_light_cannon.png"},
		ImageProjectileGatling:      {Path: "image/projectile/projectile_gatling.png"},

		ImageConstructorMerge:     {Path: "image/effect/constructor_merge.png", FrameWidth: 50},
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
	ImageUITowerSelector
	ImageUILargeSelector
	ImageUITankSelector

	ImageUIDeconstuctIcon
	ImageUIConstructors1outof1
	ImageUIConstructors1outof2
	ImageUIConstructors2outof2
	ImageUIConstructors1outof3
	ImageUIConstructors2outof3
	ImageUIConstructors3outof3

	ImageTowerBodyBunker
	ImageTowerBodyCreepBunker

	ImageTankBodyMCV
	ImageTankBodyScout
	ImageTankBodyFighter
	ImageTankBodyHunter
	ImageTankBodyDestroyer
	ImageTankBodyWheels

	ImageTankTurretScatterCannon
	ImageTankTurretLightCannon
	ImageTankTurretHurricane
	ImageTankTurretGatling

	ImageDroneConstructor
	ImageDroneCommander

	ImageGenerator
	ImageRepairDepot
	ImageTankFactory
	ImageHeavyTankFactory
	ImageCreepTankFactory

	ImageProjectileSmallMissile
	ImageProjectileLightCannon
	ImageProjectileGatling

	ImageConstructorMerge
	ImageSmallExplosion
	ImageBigExplosion
	ImageVerticalExplosion
	ImageBigVerticalExplosion
)
