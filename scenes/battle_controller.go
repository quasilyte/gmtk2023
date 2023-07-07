package scenes

import (
	"fmt"

	"github.com/quasilyte/ge"
)

type BattleController struct {
	scene *ge.Scene
}

func NewBattleController() *BattleController {
	return &BattleController{}
}

func (c *BattleController) Init(scene *ge.Scene) {
	c.scene = scene

	fmt.Println("OK")
}

func (c *BattleController) Update(delta float64) {}
