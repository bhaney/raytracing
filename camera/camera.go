package camera

import (
	"math"

	"github.com/bhaney/raytracing/ray"
	"github.com/bhaney/raytracing/utils"
	"github.com/golang/geo/r3"
)

// Camera

type Camera struct {
	lensRadius      float64
	origin          r3.Vector
	u, v, w         r3.Vector
	horizontal      r3.Vector
	vertical        r3.Vector
	lowerLeftCorner r3.Vector
}

func NewCamera(aspectRatio, vfov, aperture, focusDistance float64, lookfrom, lookat, vup r3.Vector) *Camera {
	theta := utils.DegreesToRadians(vfov)
	h := math.Tan(theta / 2.)
	viewportHeight := 2.0 * h
	viewportWidth := aspectRatio * viewportHeight

	w := lookfrom.Sub(lookat).Normalize()
	u := vup.Cross(w).Normalize()
	v := w.Cross(u)

	origin := lookfrom
	horizontal := u.Mul(viewportWidth * focusDistance)
	vertical := v.Mul(viewportHeight * focusDistance)
	lowerLeftCorner := origin.Sub(horizontal.Mul(0.5)).Sub(vertical.Mul(0.5)).Sub(w.Mul(focusDistance))
	lensRadius := aperture / 2.
	return &Camera{
		lensRadius:      lensRadius,
		origin:          origin,
		u:               u,
		v:               v,
		w:               w,
		horizontal:      horizontal,
		vertical:        vertical,
		lowerLeftCorner: lowerLeftCorner,
	}
}

func (c *Camera) GetRay(s, t float64) *ray.Ray {
	rd := utils.RandomInUnitDisk().Mul(c.lensRadius)
	offset := c.u.Mul(rd.X).Add(c.v.Mul(rd.Y))

	return &ray.Ray{c.origin.Add(offset), c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.origin).Sub(offset)}
}
