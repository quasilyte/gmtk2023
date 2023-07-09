package scenes

import (
	"os"
	"runtime"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmtk2023/assets"
	"github.com/quasilyte/gmtk2023/eui"
	"github.com/quasilyte/gmtk2023/session"
	"github.com/quasilyte/gmtk2023/styles"
)

type MainMenuController struct {
	state *session.State
}

func NewMainMenuController(state *session.State) *MainMenuController {
	return &MainMenuController{
		state: state,
	}
}

func (c *MainMenuController) Init(scene *ge.Scene) {
	bigFont := assets.BitmapFont3
	smallFont := assets.BitmapFont1

	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(320, 10, nil)
	root.AddChild(rowContainer)

	rowContainer.AddChild(eui.NewCenteredLabel("Assemblox", bigFont))
	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.SeparatorColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "PLAY", func() {
		// scene.Context().ChangeScene(NewPlayController(c.state))
	}))

	credits := eui.NewButton(c.state.UIResources, "CREDITS", func() {
		// scene.Context().ChangeScene(NewPlayController(c.state))
	})
	rowContainer.AddChild(credits)
	credits.GetWidget().Disabled = true

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))

	if runtime.GOARCH != "wasm" {
		rowContainer.AddChild(eui.NewButton(c.state.UIResources, "EXIT", func() {
			os.Exit(0)
		}))
	}

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}, styles.TransparentColor))
	rowContainer.AddChild(eui.NewCenteredLabel("GMTK 2023 edition", smallFont))

	initUI(scene, root)
}

func (c *MainMenuController) Update(delta float64) {}
