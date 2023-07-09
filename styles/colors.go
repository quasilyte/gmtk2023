package styles

import (
	"image/color"

	"github.com/quasilyte/ge"
)

var (
	HealthBarColor   = ge.RGB(0x26cd61)
	ProgressBarColor = ge.RGB(0x8b95d6)

	UnitPanelBgColor        = withAlpha(InterfaceBgColor, 150)
	UnitPanelBgOutlineColor = withAlpha(HealthBarColor, 200)

	InterfaceBgColor = ge.RGB(0x1e3a27)
)

func withAlpha(clr color.RGBA, a uint8) color.RGBA {
	clr.A = a
	return clr
}
