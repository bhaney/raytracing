package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/bhaney/raytracing/camera"
	"github.com/bhaney/raytracing/env"
	"github.com/bhaney/raytracing/ray"
	"github.com/bhaney/raytracing/utils"
	"github.com/golang/geo/r3"
)

// Scenes

func ColorRay(r *ray.Ray, world env.Hittable, depth int) ray.Color {
	rec := new(env.HitRecord)
	// If we've exceed the ray bounce limit, no more light is gathered
	if depth <= 0 {
		return ray.Color{0, 0, 0}
	}
	if world.Hit(r, 0.001, math.Inf(1), rec) { // Remember, Go is pass-by-value!
		scattered := new(ray.Ray)
		attenuation := new(ray.Color)
		if rec.MatPtr.Scatter(r, scattered, rec, attenuation) {
			resultVec := ColorRay(scattered, world, depth-1)
			return ray.Color{attenuation.X * resultVec.X, attenuation.Y * resultVec.Y, attenuation.Z * resultVec.Z}
		}
		return ray.Color{0, 0, 0}
	}
	unit := r.Direction.Normalize()
	t := 0.5 * (unit.Y + 1.)
	c1, c2 := r3.Vector{1., 1., 1.}, r3.Vector{0.5, 0.7, 1.0}
	return ray.Color(c1.Mul(1. - t).Add(c2.Mul(t)))
}

func PracticeScene() env.HittableList {
	world := make(env.HittableList, 0)

	// material
	materialGround := env.Lambertian{ray.Color{0.8, 0.8, 0.0}}
	materialCenter := env.Lambertian{ray.Color{0.1, 0.2, 0.5}}
	materialLeft := env.Dielectric{1.5}
	materialRight := env.Metal{ray.Color{0.8, 0.6, 0.2}, 0.0}
	world = append(world, &env.Sphere{r3.Vector{0, -100.5, -1}, 100, &materialGround})
	world = append(world, &env.Sphere{r3.Vector{0, 0, -1}, 0.5, &materialCenter})
	world = append(world, &env.Sphere{r3.Vector{-1, 0, -1}, 0.5, &materialLeft})
	world = append(world, &env.Sphere{r3.Vector{-1, 0, -1}, -0.4, &materialLeft})
	world = append(world, &env.Sphere{r3.Vector{1, 0, -1}, 0.5, &materialRight})
	/*
		R := math.Cos(math.Pi / 4.)
		materialLeft := env.Lambertian{Color{0.0, 0.0, 1.0}}
		materialRight := env.Lambertian{Color{1.0, 0.0, 0.0}}
		world = append(world, &env.Sphere{r3.Vector{-R, 0, -1}, R, &materialLeft})
		world = append(world, &env.Sphere{r3.Vector{R, 0, -1}, R, &materialRight})
	*/
	return world

}

func RandomScene() env.HittableList {
	world := make(env.HittableList, 0)

	groundMaterial := env.Lambertian{ray.Color{0.5, 0.5, 0.5}}
	world = append(world, &env.Sphere{r3.Vector{0, -1000., 0}, 1000, &groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := r3.Vector{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}
			if center.Sub(r3.Vector{4, 0.2, 0}).Norm() > 0.9 {
				var sphereMaterial env.Material
				if chooseMat < 0.8 {
					// diffuse
					v1, v2 := utils.RandomVector(), utils.RandomVector()
					albedo := ray.Color{v1.X * v2.X, v1.Y * v2.Y, v1.Z * v2.Z}
					sphereMaterial = &env.Lambertian{albedo}
					world = append(world, &env.Sphere{center, 0.2, sphereMaterial})
				} else if chooseMat < 0.95 {
					// metal
					albedo := ray.Color(utils.RandomVectorFromRange(0.5, 1.))
					fuzz := utils.RandomFromRange(0., 0.5)
					sphereMaterial = &env.Metal{albedo, fuzz}
					world = append(world, &env.Sphere{center, 0.2, sphereMaterial})
				} else {
					sphereMaterial = &env.Dielectric{1.5}
					world = append(world, &env.Sphere{center, 0.2, sphereMaterial})
				}
			}
		}
	}

	material1 := env.Dielectric{1.5}
	world = append(world, &env.Sphere{r3.Vector{0, 1, 0}, 1.0, &material1})
	material2 := env.Lambertian{ray.Color{0.4, 0.2, 0.1}}
	world = append(world, &env.Sphere{r3.Vector{-4, 1, 0}, 1.0, &material2})
	material3 := env.Metal{ray.Color{0.7, 0.6, 0.5}, 0.0}
	world = append(world, &env.Sphere{r3.Vector{4, 1, 0}, 1.0, &material3})

	return world

}

// Main

func main() {
	// Image
	/*
		aspectRatio := 3.0 / 2.0
		imageWidth := 1200.
		imageHeight := imageWidth / aspectRatio
		samplesPerPixel := 500
		maxDepth := 50
	*/
	aspectRatio := 16.0 / 9.0
	imageWidth := 400.
	imageHeight := imageWidth / aspectRatio
	samplesPerPixel := 100
	maxDepth := 50

	// Camera
	/*
		vfov := 20.
		lookfrom := r3.Vector{13, 2, 3}
		lookat := r3.Vector{0, 0, 0}
		vup := r3.Vector{0, 1, 0}
		distToFocus := 10.0
		aperture := 0.1
	*/
	vfov := 20.
	lookfrom := r3.Vector{3, 3, 2}
	lookat := r3.Vector{0, 0, -1}
	vup := r3.Vector{0, 1, 0}
	distToFocus := lookfrom.Sub(lookat).Norm()
	aperture := 0.1

	cam := camera.NewCamera(aspectRatio, vfov, aperture, distToFocus, lookfrom, lookat, vup)

	// World
	//world := RandomScene()
	world := PracticeScene()
	// Render

	fmt.Printf("P3\n %d %d\n255\n", int(imageWidth), int(imageHeight))

	for j := int(imageHeight) - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d \n", j)
		for i := int(imageWidth) - 1; i >= 0; i-- {
			c := ray.Color{0, 0, 0}
			for s := 0; s < samplesPerPixel; s++ {
				u := (float64(i) + rand.Float64()) / float64(imageWidth-1)
				v := (float64(j) + rand.Float64()) / float64(imageHeight-1)
				r := cam.GetRay(u, v)
				c = c.Add(ColorRay(r, world, maxDepth))
			}
			c.Write(os.Stdout, samplesPerPixel)
		}
	}
}
