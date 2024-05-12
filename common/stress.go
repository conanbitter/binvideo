package common

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/schollz/progressbar/v3"
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

func spraySampleIndex(width int, height int, x int, y int, radius float64) int {
	for {
		xs, ys := spraySampleSimple(radius)
		if xs == x || ys == y || (x+xs) < 0 || (x+xs) >= width || (y+ys) < 0 || (y+ys) >= height {
			continue
		}
		return x + xs + (y+ys)*width
	}
}

func GenerateSampler(width int, height int, radius int, samples int, iterations int) [][][]int {
	bar := progressbar.NewOptions(width*height,
		progressbar.OptionFullWidth(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionUseANSICodes(true))
	r := float64(radius)
	result := make([][][]int, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			xy := x + y*width
			bar.Set(xy)
			result[xy] = make([][]int, iterations)
			for it := 0; it < iterations; it++ {
				result[xy][it] = make([]int, samples)
				for s := 0; s < samples; s++ {
					result[xy][it][s] = spraySampleIndex(width, height, x, y, r)
				}
			}
		}
	}
	bar.Finish()
	fmt.Println()
	return result
}

func EnhanceBW(img *ImageData, radius int, samples int, iterations int) *ImageData {
	r := float64(radius)
	result := NewImage(img.Width, img.Height)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			um := float64(0)
			its := float64(iterations)
			co := img.Get(x, y)
			for it := 0; it < iterations; it++ {
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

func EnhanceBW2(img *ImageData, radius int, samples int, iterations int) *ImageData {
	r := float64(radius)
	result := NewImage(img.Width, img.Height)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			um := float64(0)
			rm := float64(0)
			co := img.Get(x, y)
			for it := 0; it < iterations; it++ {
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
				um += ui
				rm += float64(ri)
			}
			um /= float64(iterations)
			rm /= float64(iterations)
			emin := co - um*rm
			emax := emin + rm
			result.Set(x, y, (co-emin)/(emax-emin))
		}
	}
	return result
}

func EnhanceBWSampled(img *ImageData, sampler [][][]int) *ImageData {
	result := NewImage(img.Width, img.Height)
	for ind, its := range sampler {
		um := float64(0)
		co := img.Data[ind]
		for _, iter := range its {
			cmin := co
			cmax := co
			for _, sample := range iter {
				c := img.Data[sample]
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
			um += ui
		}
		um /= float64(len(its))
		result.Data[ind] = um
	}
	return result
}
