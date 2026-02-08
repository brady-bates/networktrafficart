package particle

import "image/color"

type Particle struct {
	X, Y   float32
	YDelta float32
	XSkew  float32
	Color  color.RGBA
	Size   float32
}

func NewParticle(x, y, yDelta float32, xskew float32, color color.RGBA, size float32) *Particle {
	return &Particle{
		X:      x,
		Y:      y,
		YDelta: yDelta,
		XSkew:  xskew,
		Color:  color,
		Size:   size,
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
