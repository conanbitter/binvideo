package main

import (
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"

	"github.com/dustin/go-humanize"

	com "cybercon/common"
)

func LoadFrame(filename string, noise *com.ImageData, curve []int) (gray []ImageBlock, mono []ImageBlock) {
	grayImage := com.ImageLoad(filename, -1)
	//accent := GetAccentMask(grayImage, 2, 0.2, 0.3)
	grayImage.GammaCorrection(2.4)
	//ApplyAccent(grayImage, accent)
	grayBlocks := ImageToBlocks(grayImage)
	gray = ApplyCurve(grayBlocks, curve)
	monoImage := com.BinMask(grayImage, noise)
	monoBlocks := ImageToBlocks(monoImage)
	mono = ApplyCurve(monoBlocks, curve)

	//monoImage.Save("data/result_1raw.png")
	return
}

func main() {
	options := NewEncodingOptions(1.0)
	options.SetBlack(0.5)

	filelist := make([]string, 0)
	for i := 0; i <= 3486; i++ {
		filelist = append(filelist, fmt.Sprintf("data/video/%04d.tif", i))
	}
	EncodeVideo(filelist, "data/test.bvf", options, 16.0)
	os.Exit(0)

	const count = 3486
	noise := com.ImageLoadRaw("data/noise.raw")
	temp := com.ImageLoad("data/video/0100.tif", -1)
	bw, bh := GetSizeInBlocks(temp)
	width := temp.Width
	height := temp.Height
	curve := GetHilbertCurve(bw, bh)

	bar := progressbar.NewOptions(count,
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionUseANSICodes(true))

	var lastMono []ImageBlock = nil
	var lastGray []ImageBlock = nil
	var avgComp float64 = 0
	var totalSize uint64 = 0
	for i := 0; i <= count; i++ {
		bar.Set(i)

		gray, mono := LoadFrame(fmt.Sprintf("data/video/%04d.tif", i), noise, curve)
		res := EncodeFrame(gray, mono, lastGray, options)
		avgComp += float64(res.GetSize()) / (float64(width*height) / 8.0) * 100
		totalSize += uint64(res.GetSize())

		resBlocks := res.Decode(lastMono)
		resBlocks = RevertCurve(resBlocks, curve)
		resImage := BlocksToImage(resBlocks, bw, bh)
		resImage.Save(fmt.Sprintf("data/vidcompress3/%04d.png", i))
		if i%80 != 0 {
			lastMono = res.Decode(lastMono)
			lastGray = res.DecodeGray(gray, lastGray)
		} else {
			lastGray = nil
			lastMono = nil
		}

		//fmt.Println(i)
		//image := ImageLoad(fmt.Sprintf("data/video/%04d.tif", i), -1) //2.4
		//accent := GetAccentMask(image, 2, 0.2, 0.5)
		//image.GammaCorrection(2.4)
		//ApplyAccent(image, accent)
		//BinMaskAccent(image, noise, accent)
		//BinMask(image, noise)
		//image.Save(fmt.Sprintf("data/videobin/%04d.png", i))
		//imageg := ImageLoad(fmt.Sprintf("data/test%d.png", i), 2.4)
		//BinAdaptiveFull(image)
		//image.Save(fmt.Sprintf("data/test%dadfull.png", i))
		//BinAdaptiveLocal(image)
		//BinMask4(imageg, noise, image)

		//image.Save(fmt.Sprintf("data/test%da.png", i))
		//imageg.Save(fmt.Sprintf("data/test%db.png", i))
	}
	bar.Finish()

	avgComp /= float64(count)
	fmt.Printf("\n\nAverage compression: %2.f %%\n", avgComp)
	fmt.Printf("Total size: %s\n", humanize.Bytes(totalSize))
	//image := ImageLoad("data/video/1014.tif", -1)
	//accent := GetAccentMask(image, 2, 0.2, 0.5)
	//image.GammaCorrection(2.4)
	//ApplyAccent(image, accent)
	//BinMask(image, noise)
	//BinMaskAccent(image, noise, accent)
	//image.Save("data/videotest0.png")
	//noise := GenerateBlueNoise(1024, 768)
	//noise.SaveRaw("data/noise.raw")
	//noise.Save("data/noise.png")

	/*grayA, monoA := LoadFrame("data/video/0100.tif", noise, curve)
	grayB, monoB := LoadFrame("data/video/0101.tif", noise, curve)

	res := EncodeFrame(grayB, monoB, grayA, 0.5)
	fmt.Println(len(res), len(grayB))
	fmt.Printf("Compression: %2.f %% (%d/%d)\n", float64(res.GetSize())/float64(width*height)*100, res.GetSize(), width*height)

	resBlocks := res.Decode(monoA)
	resBlocks = RevertCurve(resBlocks, curve)
	resImage := BlocksToImage(resBlocks, bw, bh)
	resImage.Save("data/result_2compress.png")*/
}
