package common

import (
	"math"
	"math/rand"
)

func SprayTest(radius int, samples int) ImageData {
	size := radius*2 + 1
	result := NewImage(size, size)
	for i := 0; i < samples; i++ {
		dist := rand.Float64() * float64(radius)
		angle := rand.Float64() * math.Pi * 2
		x := int(math.Sin(angle) * dist)
		y := int(math.Cos(angle) * dist)
		result.Set(x+radius, y+radius, 1.0)
	}
	return *result
}

func spraySampleSimple(radius float64) (int, int) {
	dist := rand.Float64() * float64(radius)
	angle := rand.Float64() * math.Pi * 2
	x := int(math.Sin(angle) * dist)
	y := int(math.Cos(angle) * dist)
	return x, y
}

func spraySample(img *ImageData, x int, y int, radius float64) float64 {
	for {
		xs, ys := spraySampleSimple(radius)
		if xs == x || ys == y || x+xs < 0 || x+xs >= img.Width || y+ys < 0 || y+ys >= img.Height {
			continue
		}
		return img.Get(x+xs, y+ys)
	}
}

func EnhanceBW(img *ImageData, radius int, samples int, iterations int) *ImageData {
	r := float64(radius)
	result := NewImage(img.Width, img.Height)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			um := float64(0)
			its := float64(iterations)
			for it := 0; it < iterations; it++ {
				co := img.Get(x, y)
				cmin := co
				cmax := co
				for s := 0; s < samples; s++ {
					c := spraySample(img, x, y, r)
					if c < cmin {
						cmin = c
					}
					if c > cmax {
						cmax = c
					}
				}
				ri := cmax - cmin
				ui := 0.5
				if ri > 0 {
					ui = float64(co-cmin) / float64(ri)
				}
				um += ui / its
			}
			result.Set(x, y, um)
		}
	}
	return result
}
