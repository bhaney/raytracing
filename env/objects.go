package env

import (
	"math"

	"github.com/bhaney/raytracing/ray"
	"github.com/golang/geo/r3"
)

// Hittable
type Hittable interface {
	Hit(r *ray.Ray, tMin, tMax float64, rec *HitRecord) bool
}

type HitRecord struct {
	Point     r3.Vector
	Normal    r3.Vector
	MatPtr    Material
	T         float64
	FrontFace bool
}

func (h *HitRecord) Update(u *HitRecord) {
	h.Point = u.Point
	h.Normal = u.Normal
	h.MatPtr = u.MatPtr
	h.T = u.T
	h.FrontFace = u.FrontFace
}

func (h *HitRecord) SetFaceNormal(r *ray.Ray, outwardNormal r3.Vector) {
	h.FrontFace = r.Direction.Dot(outwardNormal) < 0
	if h.FrontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.Mul(-1.)
	}
}

// Hittable List

type HittableList []Hittable

func (hl HittableList) Hit(r *ray.Ray, tMin, tMax float64, rec *HitRecord) bool {
	tempRec := new(HitRecord)
	hitAnything := false
	closestSoFar := tMax
	for _, object := range hl {
		if object.Hit(r, tMin, closestSoFar, tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = *tempRec
		}
	}
	return hitAnything
}

// Spheres

type Sphere struct {
	Center r3.Vector
	Radius float64
	MatPtr Material
}

func (s *Sphere) Hit(r *ray.Ray, tMin, tMax float64, rec *HitRecord) bool {
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.Norm2()
	halfB := oc.Dot(r.Direction)
	c := oc.Norm2() - s.Radius*s.Radius
	discrim := halfB*halfB - a*c
	if discrim < 0 {
		return false
	}
	sqrtD := math.Sqrt(discrim)
	//find nearest root
	root := (-halfB - sqrtD) / a
	if root < tMin || tMax < root {
		root = (-halfB + sqrtD) / a
		if root < tMin || tMax < root {
			return false
		}
	}
	rec.T = root
	rec.Point = r.At(rec.T)
	outwardNormal := rec.Point.Sub(s.Center).Mul(1. / s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.MatPtr = s.MatPtr
	return true

}
