package ray

import (
	"fmt"
	"io"
	"math"

	"github.com/bhaney/raytracing/utils"
	"github.com/golang/geo/r3"
)

// Rays
type Ray struct {
	Origin    r3.Vector
	Direction r3.Vector
}

func (r *Ray) At(t float64) r3.Vector {
	return r.Origin.Add(r.Direction.Mul(t))
}

// Color
type Color r3.Vector

func (c Color) Add(d Color) Color {
	return Color(r3.Vector(c).Add(r3.Vector(d)))
}

func (c Color) Mul(m float64) Color {
	return Color(r3.Vector(c).Mul(m))
}

func (c *Color) Write(w io.Writer, samplesPerPixel int) {
	r, g, b := c.X, c.Y, c.Z
	scale := 1.0 / float64(samplesPerPixel)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)
	fmt.Fprintf(w, "%d %d %d\n", int(256.*utils.Clamp(r, 0., .999)), int(256.*utils.Clamp(g, 0., .999)), int(256.*utils.Clamp(b, 0., .999)))
}
