package main

func BinFixed(image *ImageData, treshold float64) {
	for i, point := range image.Data {
		if point > treshold {
			image.Data[i] = 1.0
		} else {
			image.Data[i] = 0.0
		}
	}
}

func meanTreshold(image *ImageData) float64 {
	valacc := float64(0)
	for _, v := range image.Data {
		valacc += v
	}
	return valacc / float64(image.Width*image.Height)
}

func BinAdaptiveFull(image *ImageData) {
	treshold := meanTreshold(image)
	BinFixed(image, treshold)
}
