package common

import (
	"encoding/binary"
	"image"
	"image/png"
	"log"
	"math"
	"os"

	_ "image/jpeg"

	_ "golang.org/x/image/tiff"
)

type ImageData struct {
	Width  int
	Height int
	Data   []float64
}

func clipColor(data float64) uint8 {
	if data < 0.0 {
		data = 0.0
	}
	if data > 1.0 {
		data = 1.0
	}
	return uint8(data * 255)
}

func clipFloat(data float64) float64 {
	if data < 0.0 {
		return 0.0
	}
	if data > 1.0 {
		return 1.0
	}
	return data
}

func (img *ImageData) Save(filename string) {
	output := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{img.Width, img.Height}})
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			output.Pix[output.PixOffset(x, y)] = clipColor(img.Get(x, y))
		}
	}
	outf, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer outf.Close()
	if err = png.Encode(outf, output); err != nil {
		log.Printf("failed to encode: %v", err)
	}
}

func (img *ImageData) SaveRaw(filename string) {
	fo, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	binary.Write(fo, binary.LittleEndian, uint32(img.Width))
	binary.Write(fo, binary.LittleEndian, uint32(img.Height))
	binary.Write(fo, binary.LittleEndian, img.Data)
}

func (img *ImageData) GetBytes() [][]byte {
	cols := img.Width
	rows := img.Height / 8
	result := make([][]byte, rows)
	for y := 0; y < rows; y++ {
		result[y] = make([]byte, cols)
		for x := 0; x < cols; x++ {
			color := 0
			for s := 0; s < 8; s++ {
				if img.Get(x, y*8+s) > 0.5 {
					color = color | 1<<s
				}
			}
			result[y][x] = byte(color)
		}
	}
	return result
}

func (img *ImageData) Set(x int, y int, color float64) {
	img.Data[x+y*img.Width] = color
}

func (img *ImageData) Get(x int, y int) float64 {
	return img.Data[x+y*img.Width]
}

func (img *ImageData) GammaCorrection(gamma float64) {
	for i, c := range img.Data {
		if c <= 0.04045 {
			continue
		}
		img.Data[i] = math.Pow((c+0.055)/1.055, gamma)
	}
}

func NewImage(width int, height int) *ImageData {
	result := &ImageData{
		Width:  width,
		Height: height,
		Data:   make([]float64, width*height),
	}
	clear(result.Data)
	return result
}

func ImageFrom(other *ImageData) *ImageData {
	res := &ImageData{
		Width:  other.Width,
		Height: other.Height,
		Data:   make([]float64, other.Width*other.Height),
	}
	copy(res.Data, other.Data)
	return res
}

func ColorConvert(r uint32, g uint32, b uint32, gamma float64) float64 {
	rf := float64(r) / 65535.0
	gf := float64(g) / 65535.0
	bf := float64(b) / 65535.0

	c := 0.299*rf + 0.587*gf + 0.114*bf

	if gamma < 0 {
		return c
	}
	// Gamma correction
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, gamma)
}

func ColorConvertFloat(r float64, g float64, b float64) float64 {
	return 0.299*r + 0.587*g + 0.114*b
}

func ImageLoad(filename string, gamma float64) *ImageData {
	imgFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	result := NewImage((bounds.Max.X - bounds.Min.X), (bounds.Max.Y - bounds.Min.Y))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			result.Set(x-bounds.Min.X, y-bounds.Min.Y, ColorConvert(r, g, b, gamma))
		}
	}
	return result
}

func ImageLoadRaw(filename string) *ImageData {
	fo, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	result := &ImageData{}

	var width uint32
	var height uint32
	binary.Read(fo, binary.LittleEndian, &width)
	binary.Read(fo, binary.LittleEndian, &height)
	result.Width = int(width)
	result.Height = int(height)
	result.Data = make([]float64, result.Width*result.Height)
	binary.Read(fo, binary.LittleEndian, result.Data)

	return result
}

func ImageLoadDecompose(filename string) []*ImageData {
	imgFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	result := []*ImageData{NewImage(width, height), NewImage(width, height), NewImage(width, height)}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rf := float64(r) / 65535.0
			gf := float64(g) / 65535.0
			bf := float64(b) / 65535.0
			result[0].Set(x-bounds.Min.X, y-bounds.Min.Y, rf)
			result[1].Set(x-bounds.Min.X, y-bounds.Min.Y, gf)
			result[2].Set(x-bounds.Min.X, y-bounds.Min.Y, bf)
		}
	}
	return result
}

func ImageCompose(channels []*ImageData) *ImageData {
	result := NewImage(channels[0].Width, channels[0].Height)
	for i := range result.Data {
		result.Data[i] = ColorConvertFloat(
			channels[0].Data[i],
			channels[1].Data[i],
			channels[2].Data[i],
		)
	}
	return result
}

func ImageMerge(img1 *ImageData, img2 *ImageData, alpha float64) *ImageData {
	result := NewImage(img1.Width, img1.Height)
	for i := range img1.Data {
		result.Data[i] = clipFloat(img1.Data[i]*(1.0-alpha) + img2.Data[i]*alpha)
	}
	return result
}
