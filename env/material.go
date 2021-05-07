package env

import (
	"math"
	"math/rand"

	"github.com/bhaney/raytracing/ray"
	"github.com/bhaney/raytracing/utils"
	"github.com/golang/geo/r3"
)

// Material
type Material interface {
	Scatter(rIn, scattered *ray.Ray, rec *HitRecord, attenuation *ray.Color) bool
}

type Lambertian struct {
	Albedo ray.Color
}

func (l *Lambertian) Scatter(rIn, scattered *ray.Ray, rec *HitRecord, attenuation *ray.Color) bool {
	scatterDirection := rec.Normal.Add(utils.RandomUnitVector())
	if utils.NearZero(scatterDirection) {
		scatterDirection = rec.Normal
	}
	*scattered = ray.Ray{rec.Point, scatterDirection}
	*attenuation = l.Albedo
	return true
}

type Metal struct {
	Albedo ray.Color
	Fuzz   float64
}

func (m *Metal) Scatter(rIn, scattered *ray.Ray, rec *HitRecord, attenuation *ray.Color) bool {
	reflected := Reflect(rIn.Direction.Normalize(), rec.Normal)
	*scattered = ray.Ray{rec.Point, reflected.Add(utils.RandomInUnitSphere().Mul(m.Fuzz))}
	*attenuation = m.Albedo
	return scattered.Direction.Dot(rec.Normal) > 0
}

type Dielectric struct {
	IR float64
}

func (d *Dielectric) Reflectance(cosine, refIdx float64) float64 {
	// Schlick's approximation
	r0 := (1. - refIdx) / (1. + refIdx)
	r0 = r0 * r0
	return r0 + (1.-r0)*math.Pow(1-cosine, 5.)
}

func (d *Dielectric) Scatter(rIn, scattered *ray.Ray, rec *HitRecord, attenuation *ray.Color) bool {
	*attenuation = ray.Color{1.0, 1.0, 1.0}
	var refractionRatio float64
	if rec.FrontFace {
		refractionRatio = 1. / d.IR
	} else {
		refractionRatio = d.IR
	}
	unitDirection := rIn.Direction.Normalize()
	cosTheta := math.Min(unitDirection.Mul(-1).Dot(rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	cannotRefract := refractionRatio*sinTheta > 1.0
	var direction r3.Vector
	if cannotRefract || d.Reflectance(cosTheta, refractionRatio) > rand.Float64() {
		direction = Reflect(unitDirection, rec.Normal)
	} else {
		direction = Refract(unitDirection, rec.Normal, refractionRatio)
	}
	*scattered = ray.Ray{rec.Point, direction}
	return true
}

func Reflect(v, n r3.Vector) r3.Vector {
	return v.Sub(n.Mul(2. * v.Dot(n)))
}

func Refract(uv, n r3.Vector, etaiOverEtat float64) r3.Vector {
	cosTheta := math.Min(uv.Mul(-1).Dot(n), 1.0)
	rOutPerp := uv.Add(n.Mul(cosTheta)).Mul(etaiOverEtat)
	rOutParallel := n.Mul(-math.Sqrt(math.Abs(1.0 - rOutPerp.Norm2())))
	return rOutPerp.Add(rOutParallel)
}
