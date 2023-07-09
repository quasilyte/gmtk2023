package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/eui"
)

func initUI(scene *ge.Scene, root *widget.Container) {
	uiObject := eui.NewSceneObject(root)
	scene.AddGraphics(uiObject)
	scene.AddObject(uiObject)
}
