//go:build ignore
// +build ignore

package main

var Time float

func tex2pixCoord(texCoord vec2) vec2 {
	pixSize := imageSrcTextureSize()
	originTexCoord, _ := imageSrcRegionOnTexture()
	actualTexCoord := texCoord - originTexCoord
	actualPixCoord := actualTexCoord * pixSize
	return actualPixCoord
}

func pix2texCoord(actualPixCoord vec2) vec2 {
	pixSize := imageSrcTextureSize()
	actualTexCoord := actualPixCoord / pixSize
	originTexCoord, _ := imageSrcRegionOnTexture()
	texCoord := actualTexCoord + originTexCoord
	return texCoord
}

func shaderRand(pixCoord vec2) int {
	return int(pixCoord.x+pixCoord.y) * int(pixCoord.y*5)
}

func applyPixPick(pixCoord vec2, dist float, m, hash int) vec2 {
	dir := hash % m
	if dir == int(0) {
		pixCoord.x += dist
	} else if dir == int(1) {
		pixCoord.x -= dist
	} else if dir == int(2) {
		pixCoord.y += dist
	} else if dir == int(3) {
		pixCoord.y -= dist
	}
	return pixCoord
}

func gridPattern(v, colorMult vec4, hash int, p, pixSize, originTexPos vec2) vec4 {
	posHash := int(p.x+p.y) * int(p.y*5)
	state := posHash % hash
	if state == int(1) {
		p.x += 1.0
	} else if state == int(2) {
		p.x -= 1.0
	} else if state == int(3) {
		p.y += 1.0
	} else if state == int(4) {
		p.y -= 1.0
	} else {
		return v
	}
	return imageSrc0At(p/pixSize+originTexPos) * colorMult
}

func sourceSize() vec2 {
	pixSize := imageSrcTextureSize()
	_, srcRegion := imageSrcRegionOnTexture()
	return pixSize.x * srcRegion
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	c := imageSrc0UnsafeAt(texCoord)
	if c.a == 0 {
		return c
	}

	actualPixPos := tex2pixCoord(texCoord)

	sizes := sourceSize()
	width := sizes.x // Изображение квадратное, поэтому достаточно width

	initialY := -2.0
	offsetY := width * 0.15 * Time
	circleCenter := vec2(width*0.5, initialY-offsetY)
	dist := distance(actualPixPos, circleCenter)

	progress := 1.4 - Time
	if dist > ((width * 0.95) * progress) {
		return c
	}

	spread := 0
	colorMultiplier := vec4(0)

	if dist > ((width * 0.85) * progress) {
		spread = 15
		colorMultiplier = vec4(1.0, 1.1, 1.3, 1.0)
	} else if dist > ((width * 0.75) * progress) {
		spread = 11
		colorMultiplier = vec4(0.9, 1.2, 1.6, 1.0)
	} else if dist > ((width * 0.65) * progress) {
		spread = 7
		colorMultiplier = vec4(0.8, 1.4, 2.0, 1.0)
	} else if dist > ((width * 0.62) * progress) {
		spread = 6
		colorMultiplier = vec4(0.25, 0.25, 0.25, 1.0)
	} else {
		return vec4(0)
	}

	h := shaderRand(actualPixPos)
	p := applyPixPick(actualPixPos, 1, spread, h)
	if p == actualPixPos {
		return c
	}
	return imageSrc0At(pix2texCoord(p)) * colorMultiplier
}
