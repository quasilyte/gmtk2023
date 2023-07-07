package battle

import "github.com/quasilyte/gmtk2023/gamedata"

type Runner struct {
	config *gamedata.BattleConfig
}

func NewRunner(config *gamedata.BattleConfig) *Runner {
	return &Runner{
		config: config,
	}
}
