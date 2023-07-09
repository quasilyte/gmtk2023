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
	ctx.GameName = "assemblox"
	ctx.WindowTitle = "Assemblox"
	ctx.WindowWidth = 1920 / 2
	ctx.WindowHeight = 1080 / 2
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	state := session.NewState()

	playerInput := controls.MakeHandler(ctx)

	playerDesigns := gamedata.NewPlayerDesigns()
	playerDesigns.Tanks[0] = &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.ScoutBodyStats,
		Turret:   gamedata.HurricaneStats,
	}
	playerDesigns.Tanks[1] = &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.FighterBodyStats,
		Turret:   gamedata.LightCannonStats,
	}
	playerDesigns.Tanks[2] = &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.HunterBodyStats,
		Turret:   gamedata.ScatterCannonStats,
	}
	playerDesigns.Tanks[3] = &gamedata.UnitStats{
		Movement: gamedata.UnitMovementGround,
		Body:     gamedata.DestroyerBodyStats,
		Turret:   gamedata.ScatterCannonStats,
	}

	playerDesigns.Towers[0] = &gamedata.UnitStats{
		Movement:   gamedata.UnitMovementNone,
		Body:       gamedata.BunkerBodyStats,
		Turret:     gamedata.GatlingStats,
		Selectable: true,
	}
	playerDesigns.Towers[1] = &gamedata.UnitStats{
		Movement:   gamedata.UnitMovementNone,
		Body:       gamedata.BunkerBodyStats,
		Turret:     gamedata.LightCannonStats,
		Selectable: true,
	}

	config := &gamedata.BattleConfig{
		PlayerInput:   playerInput,
		GameSpeed:     1,
		PlayerDesigns: playerDesigns,
	}

	if err := ge.RunGame(ctx, scenes.NewBattleController(state, config)); err != nil {
		panic(err)
	}
}
