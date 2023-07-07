package controls

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
)

const (
	ActionUnknown input.Action = iota

	ActionPanRight
	ActionPanDown
	ActionPanLeft
	ActionPanUp

	ActionSelectUnit
	ActionSendUnit
)

func MakeHandler(ctx *ge.Context) *input.Handler {
	// gamepadKeymap := input.Keymap{
	// 	ActionPanRight:   {input.KeyGamepadLStickRight, input.KeyGamepadRight},
	// 	ActionPanDown:    {input.KeyGamepadLStickDown, input.KeyGamepadDown},
	// 	ActionPanLeft:    {input.KeyGamepadLStickLeft, input.KeyGamepadLeft},
	// 	ActionPanUp:      {input.KeyGamepadLStickUp, input.KeyGamepadUp},
	// 	ActionSelectUnit: {input.KeyGamepadA},
	// }

	keyboardKeymap := input.Keymap{
		ActionPanRight:   {input.KeyD, input.KeyRight},
		ActionPanDown:    {input.KeyS, input.KeyDown},
		ActionPanLeft:    {input.KeyA, input.KeyLeft},
		ActionPanUp:      {input.KeyW, input.KeyUp},
		ActionSelectUnit: {input.KeyMouseLeft},
		ActionSendUnit:   {input.KeyMouseRight},
	}

	combinedKeymap := input.Keymap{}

	// for a, keys := range gamepadKeymap {
	// 	combinedKeymap[a] = append(combinedKeymap[a], keys...)
	// }
	for a, keys := range keyboardKeymap {
		combinedKeymap[a] = append(combinedKeymap[a], keys...)
	}

	return ctx.Input.NewHandler(0, combinedKeymap)
}
