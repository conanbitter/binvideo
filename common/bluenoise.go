package common

import (
	"math"
	"math/rand"

	"github.com/schollz/progressbar/v3"
)

const sigma = 1.5
const divisor = sigma * sigma * 2

type point struct {
	X int
	Y int
}

func generateLut(width int, height int) *ImageData {
	res := NewImage(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var xd int
			var yd int
			if x < width/2 {
				xd = x
			} else {
				xd = width - x
			}
			if y < height/2 {
				yd = y
			} else {
				yd = height - y
			}
			distance := float64(xd*xd + yd*yd)
			res.Set(
				x,
				y,
				math.Exp(-distance/divisor),
			)
		}
	}
	//res.Set(0, 0, math.MaxFloat64)
	return res
}

func startFill(image *ImageData, mask *ImageData, lut *ImageData, quantity float32) int {
	// Jittered grid method
	cols := int(float32(image.Width) * quantity)
	cell_width := image.Width / cols
	rows := int(float32(image.Height) * quantity)
	cell_height := image.Height / rows
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			xr := x*cell_width + rand.Intn(cell_width)
			yr := y*cell_height + rand.Intn(cell_height)
			applyPoint(xr, yr, true, image, mask, lut)
		}
	}
	return cols * rows
}

func applyPoint(x int, y int, add bool, image *ImageData, mask *ImageData, lut *ImageData) {
	var sign float64
	if add {
		image.Set(x, y, 1.0)
		sign = 1.0
	} else {
		image.Set(x, y, 0.0)
		sign = -1.0
	}

	for yp := 0; yp < mask.Height; yp++ {
		ylut := yp - y
		if ylut < 0 {
			ylut += mask.Height
		}
		for xp := 0; xp < mask.Width; xp++ {
			xlut := xp - x
			if xlut < 0 {
				xlut += mask.Width
			}
			mask.Set(xp, yp, mask.Get(xp, yp)+lut.Get(xlut, ylut)*sign)
		}
	}
}

func findVoid(image *ImageData, mask *ImageData) (int, int) {
	var (
		voidX     int
		voidY     int
		minEnergy float64 = math.MaxFloat64
	)
	for y := 0; y < image.Height; y++ {
		for x := 0; x < image.Width; x++ {
			enrg := mask.Get(x, y)
			if image.Get(x, y) < 0.5 && enrg < minEnergy {
				minEnergy = enrg
				voidX = x
				voidY = y
			}
		}
	}
	return voidX, voidY
}

func findCluster(image *ImageData, mask *ImageData) (int, int) {
	var (
		clusterX  int
		clusterY  int
		maxEnergy float64 = math.SmallestNonzeroFloat64
	)
	for y := 0; y < image.Height; y++ {
		for x := 0; x < image.Width; x++ {
			enrg := mask.Get(x, y)
			if image.Get(x, y) > 0.5 && enrg > maxEnergy {
				maxEnergy = enrg
				clusterX = x
				clusterY = y
			}
		}
	}
	return clusterX, clusterY
}

func generateBlueNoisePoints(width int, height int) []point {
	bar := progressbar.NewOptions(width*height,
		progressbar.OptionFullWidth(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionUseANSICodes(true))

	res := NewImage(width, height)
	energyMask := NewImage(width, height)
	lut := generateLut(width, height)

	firstPointsCount := startFill(res, energyMask, lut, 0.1)
	points := make([]point, firstPointsCount)
	//points[1] = Point{X: 1, Y: 1}

	//Step 1

	bar.Describe("Step 1")

	ind := 0
	for {
		var (
			voidX    int
			voidY    int
			clusterX int
			clusterY int
		)

		clusterX, clusterY = findCluster(res, energyMask)
		applyPoint(clusterX, clusterY, false, res, energyMask, lut)

		voidX, voidY = findVoid(res, energyMask)
		applyPoint(voidX, voidY, true, res, energyMask, lut)
		if voidX == clusterX && voidY == clusterY {
			break
		}
		ind++
	}

	bar.Describe("Step 2")
	bar.Set(0)

	// Step 2

	step2temp := ImageFrom(res)
	step2mask := ImageFrom(energyMask)

	for i := firstPointsCount - 1; i >= 0; i-- {
		var (
			clusterX int
			clusterY int
		)
		bar.Set(firstPointsCount - i)
		clusterX, clusterY = findCluster(step2temp, step2mask)
		points[i] = point{X: clusterX, Y: clusterY}
		applyPoint(clusterX, clusterY, false, step2temp, step2mask, lut)
	}

	bar.Describe("Step 3")

	// Step 3
	for c := firstPointsCount; c < width*height/2; c++ {
		var (
			voidX int
			voidY int
		)
		bar.Set(c)
		voidX, voidY = findVoid(res, energyMask)
		points = append(points, point{X: voidX, Y: voidY})
		applyPoint(voidX, voidY, true, res, energyMask, lut)
	}

	bar.Describe("Step 4")

	// Step 4
	negative := NewImage(width, height)
	negEnergy := NewImage(width, height)
	for y := 0; y < negative.Height; y++ {
		for x := 0; x < negative.Width; x++ {
			if res.Get(x, y) < 0.5 {
				applyPoint(x, y, true, negative, negEnergy, lut)
			}
		}
	}
	for c := width * height / 2; c < width*height; c++ {
		var (
			clusterX int
			clusterY int
		)
		bar.Set(c)
		clusterX, clusterY = findCluster(negative, negEnergy)
		points = append(points, point{X: clusterX, Y: clusterY})
		applyPoint(clusterX, clusterY, false, negative, negEnergy, lut)
	}

	bar.Finish()
	return points
}

func GenerateBlueNoise(width int, height int) *ImageData {
	points := generateBlueNoisePoints(width, height)
	noise := NewImage(width, height)
	for i, point := range points {
		noise.Set(point.X, point.Y, float64(i)/float64(len(points)))
	}

	return noise
}
