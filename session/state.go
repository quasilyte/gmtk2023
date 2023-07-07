package session

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

type State struct {
	Resources *Resources
}

func NewState() *State {
	return &State{
		Resources: &Resources{},
	}
}

type Resources struct {
	TankBodyDestroyer resource.Image

	TankTurretLightCannon resource.Image
}
