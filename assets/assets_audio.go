package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerAudioResources(ctx *ge.Context) {
	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioExplosion1: {Path: "audio/explosion1.wav", Volume: -0.5},
		AudioExplosion2: {Path: "audio/explosion2.wav", Volume: -0.5},
		AudioExplosion3: {Path: "audio/explosion3.wav", Volume: -0.5},
		AudioExplosion4: {Path: "audio/explosion4.wav", Volume: -0.5},
		AudioExplosion5: {Path: "audio/explosion5.wav", Volume: -0.5},

		AudioShotLightCannon: {Path: "audio/shot_light_cannon.wav", Volume: -0.35},

		AudioShotGatling: {Path: "audio/shot_gatling.wav", Volume: -0.3},
	}

	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.LoadAudio(id)
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioExplosion1
	AudioExplosion2
	AudioExplosion3
	AudioExplosion4
	AudioExplosion5

	AudioShotLightCannon

	AudioShotGatling
)
