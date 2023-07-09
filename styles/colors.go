package styles

import (
	"image/color"

	"github.com/quasilyte/ge"
)

var (
	HealthBarColor   = ge.RGB(0x26cd61)
	ProgressBarColor = ge.RGB(0x8b95d6)

	FontColor         = ge.RGB(0x8ae6a2)
	NormalTextColor   = FontColor
	DisabledTextColor = ge.RGB(0x5e7364)

	SeparatorColor = DisabledTextColor

	TransparentColor = color.RGBA{}

	UnitPanelBgColor        = withAlpha(ge.RGB(0x23252c), 255/2)
	UnitPanelBgOutlineColor = ge.RGB(0x1f9941)

	InterfaceBgColor = ge.RGB(0x1e3a27)
)

func withAlpha(clr color.RGBA, a uint8) color.RGBA {
	clr.A = a
	return clr
}
