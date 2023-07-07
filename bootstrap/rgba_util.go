package bootstrap

import (
	"image/color"
)

func rgbaSub(clr color.RGBA, v uint8) color.RGBA {
	if clr.R <= v {
		clr.R = 0
	} else {
		clr.R -= v
	}
	if clr.G <= v {
		clr.G = 0
	} else {
		clr.G -= v
	}
	if clr.B <= v {
		clr.B = 0
	} else {
		clr.B -= v
	}
	return clr
}

func rgbaAdd(clr color.RGBA, v uint8) color.RGBA {
	if int(clr.R)+int(v) >= 0xff {
		clr.R = 255
	} else {
		clr.R += v
	}
	if int(clr.G)+int(v) >= 0xff {
		clr.G = 255
	} else {
		clr.G += v
	}
	if int(clr.B)+int(v) >= 0xff {
		clr.B = 255
	} else {
		clr.B += v
	}
	return clr
}
