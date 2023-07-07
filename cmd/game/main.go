package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/controls"
	"github.com/quasilyte/gmtk2023/gamedata"
	"github.com/quasilyte/gmtk2023/scenes"
	"github.com/quasilyte/gmtk2023/session"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "quasilte_gmtk2023"
	ctx.WindowTitle = "GMTK2023"
	ctx.WindowWidth = 1920 / 2
	ctx.WindowHeight = 1080 / 2
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	state := session.NewState()

	playerInput := controls.MakeHandler(ctx)

	config := &gamedata.BattleConfig{
		PlayerInput: playerInput,
		GameSpeed:   1,
	}

	if err := ge.RunGame(ctx, scenes.NewBattleController(state, config)); err != nil {
		panic(err)
	}
}
