package utils

import (
	"math"
	"math/rand"

	"github.com/golang/geo/r3"
)

func DegreesToRadians(degrees float64) float64 {
	return (degrees * math.Pi) / 180.0
}

func RandomFromRange(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func RandomVector() r3.Vector {
	return r3.Vector{rand.Float64(), rand.Float64(), rand.Float64()}
}

func RandomVectorFromRange(min, max float64) r3.Vector {
	return r3.Vector{RandomFromRange(min, max), RandomFromRange(min, max), RandomFromRange(min, max)}
}

func RandomInUnitDisk() r3.Vector {
	for {
		p := r3.Vector{RandomFromRange(-1., 1.), RandomFromRange(-1., 1.), 0}
		if p.Norm2() >= 1 {
			continue
		}
		return p
	}
}

func RandomInUnitSphere() r3.Vector {
	for {
		p := RandomVectorFromRange(-1., 1.)
		if p.Norm2() >= 1. {
			continue
		}
		return p
	}
}

func RandomUnitVector() r3.Vector {
	return RandomInUnitSphere().Normalize()
}

func NearZero(v r3.Vector) bool {
	s := 1e-8
	return math.Abs(v.X) < s && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}
