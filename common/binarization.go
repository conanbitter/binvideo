package common

import "math"

func BinFixed(image *ImageData, treshold float64) {
	for i, point := range image.Data {
		if point > treshold {
			image.Data[i] = 1.0
		} else {
			image.Data[i] = 0.0
		}
	}
}

func meanValue(data []float64) float64 {
	valacc := float64(0)
	for _, v := range data {
		valacc += v
	}
	return valacc / float64(len(data))
}

func BinAdaptiveFull(image *ImageData) {
	treshold := meanValue(image.Data)
	BinFixed(image, treshold)
}

func getKernel(pixel int, image *ImageData, radius int) []float64 {
	x := pixel % image.Width
	y := pixel / image.Width

	x1 := x - radius
	if x1 < 0 {
		x1 = 0
	}
	x2 := x + radius
	if x2 > image.Width-1 {
		x2 = image.Width - 1
	}
	y1 := y - radius
	if y1 < 0 {
		y1 = 0
	}
	y2 := y + radius
	if y2 > image.Height-1 {
		y2 = image.Height - 1
	}

	result := make([]float64, 0)
	for j := y1; j <= y2; j++ {
		for i := x1; i <= x2; i++ {
			result = append(result, image.Data[i+j*image.Width])
		}
	}
	return result
}

func devValue(data []float64, mean float64) float64 {
	valacc := float64(0)
	for _, v := range data {
		diff := v - mean
		valacc += diff * diff
	}
	return math.Sqrt(valacc / float64(len(data)))
}

const (
	kernelRadius = 2            //5
	paramK       = float64(0.2) //0.2
	paramR       = float64(128) / 256
	lumMax       = float64(0.9)
	lumMin       = float64(0.1)
)

func BinAdaptiveLocal(image *ImageData) {
	newImage := make([]float64, len(image.Data))
	for i, point := range image.Data {
		if point < lumMin {
			newImage[i] = 0.0
			continue
		}
		if point > lumMax {
			newImage[i] = 1.0
			continue
		}
		kernel := getKernel(i, image, kernelRadius)
		mean := meanValue(kernel)
		dev := devValue(kernel, mean)
		treshold := mean * (1 + paramK*(dev/paramR-1))
		if point > treshold {
			newImage[i] = 1.0
		} else {
			newImage[i] = 0.0
		}
	}
	image.Data = newImage
}

func BinAdaptiveLocal2(image *ImageData) {
	newImage := make([]float64, len(image.Data))
	for i, point := range image.Data {
		kernel := getKernel(i, image, kernelRadius)
		mean := meanValue(kernel)
		dev := devValue(kernel, mean)
		treshold := mean * (1 + paramK*(dev/paramR-1))
		newImage[i] = (point - treshold) * 15
		/*if point > treshold+0.1 {
			newImage[i] = 1.0
		} else if point < treshold-0.1 {
			newImage[i] = 0.0
		} else {
			newImage[i] = point
		}*/

	}
	image.Data = newImage
}

func BinMask(image *ImageData, mask *ImageData) *ImageData {
	result := NewImage(image.Width, image.Height)
	for i, point := range image.Data {
		if point > mask.Data[i] {
			result.Data[i] = 1.0
		} else {
			result.Data[i] = 0.0
		}
	}
	return result
}

func BinMask2(image *ImageData, mask *ImageData, bin *ImageData) {
	for i, point := range image.Data {
		if point*bin.Data[i] > mask.Data[i] {
			image.Data[i] = 1.0
		} else {
			image.Data[i] = 0.0
		}
	}
}

func BinMask3(image *ImageData, mask *ImageData, bin *ImageData, k float64) {
	for i, point := range image.Data {
		s := math.Pow(bin.Data[i], k)
		if point*s > mask.Data[i] {
			image.Data[i] = 1.0
		} else {
			image.Data[i] = 0.0
		}
	}
}

func BinMask4(image *ImageData, mask *ImageData, bin *ImageData) {
	for i, point := range image.Data {
		s := math.Sqrt(bin.Data[i])
		if point*s > mask.Data[i] {
			image.Data[i] = 1.0
		} else {
			image.Data[i] = 0.0
		}
	}
}

func GetAccentMask(image *ImageData, kernelRadius int, k float64, amplitude float64) *ImageData {
	const R = float64(0.5)
	const level = float64(15)
	newImage := make([]float64, len(image.Data))
	for i, point := range image.Data {
		kernel := getKernel(i, image, kernelRadius)
		mean := meanValue(kernel)
		dev := devValue(kernel, mean)
		treshold := mean * (1 + k*(dev/R-1))
		newImage[i] = math.Pow((point-treshold)*level, amplitude)
	}
	return &ImageData{
		Width:  image.Width,
		Height: image.Height,
		Data:   newImage,
	}
}

func ApplyAccent(image *ImageData, accent *ImageData) {
	for i, point := range image.Data {
		image.Data[i] = point * accent.Data[i]
	}
}
