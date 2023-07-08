package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerAudioResources(ctx *ge.Context) {
	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioExplosion1: {Path: "audio/explosion1.wav", Volume: -0.65},
		AudioExplosion2: {Path: "audio/explosion2.wav", Volume: -0.65},
		AudioExplosion3: {Path: "audio/explosion3.wav", Volume: -0.65},
		AudioExplosion4: {Path: "audio/explosion4.wav", Volume: -0.65},
		AudioExplosion5: {Path: "audio/explosion5.wav", Volume: -0.65},

		AudioUnitAck1: {Path: "audio/unit_ack1.wav", Volume: -0.55},
		AudioUnitAck2: {Path: "audio/unit_ack2.wav", Volume: -0.55},
		AudioUnitAck3: {Path: "audio/unit_ack3.wav", Volume: -0.55},
		AudioUnitAck4: {Path: "audio/unit_ack4.wav", Volume: -0.55},
		AudioUnitAck5: {Path: "audio/unit_ack5.wav", Volume: -0.55},
		AudioUnitAck6: {Path: "audio/unit_ack6.wav", Volume: -0.55},

		AudioShotLightCannon1: {Path: "audio/shot_light_cannon1.wav", Volume: -0.3},
		AudioShotLightCannon2: {Path: "audio/shot_light_cannon2.wav", Volume: -0.3},
		AudioShotLightCannon3: {Path: "audio/shot_light_cannon3.wav", Volume: -0.3},

		AudioShotGatling: {Path: "audio/shot_gatling.wav", Volume: -0.8},
	}

	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioExplosion1:
		return 5
	case AudioShotLightCannon1:
		return 3
	case AudioUnitAck1:
		return 6
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioExplosion1
	AudioExplosion2
	AudioExplosion3
	AudioExplosion4
	AudioExplosion5

	AudioUnitAck1
	AudioUnitAck2
	AudioUnitAck3
	AudioUnitAck4
	AudioUnitAck5
	AudioUnitAck6

	AudioShotLightCannon1
	AudioShotLightCannon2
	AudioShotLightCannon3

	AudioShotGatling
)
