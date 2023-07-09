package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/controls"
	"github.com/quasilyte/gmtk2023/eui"
	"github.com/quasilyte/gmtk2023/scenes"
	"github.com/quasilyte/gmtk2023/session"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "assemblox"
	ctx.WindowTitle = "Assemblox"
	ctx.WindowWidth = 1920 / 2
	ctx.WindowHeight = 1080 / 2
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	playerInput := controls.MakeHandler(ctx)

	state := session.NewState()
	state.UIResources = eui.PrepareResources(ctx.Loader)
	state.Input = playerInput

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}
