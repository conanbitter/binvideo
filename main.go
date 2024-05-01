package main

import "fmt"

func main() {
	const count = 3486
	noise := ImageLoadRaw("data/noise.raw")
	/*
		bar := progressbar.NewOptions(count,
			progressbar.OptionFullWidth(),
			progressbar.OptionShowCount(),
			progressbar.OptionUseANSICodes(true))

		for i := 0; i <= count; i++ {
			bar.Set(i)
			//fmt.Println(i)
			image := ImageLoad(fmt.Sprintf("data/video/%04d.tif", i), -1) //2.4
			//accent := GetAccentMask(image, 2, 0.2, 0.5)
			image.GammaCorrection(2.4)
			//ApplyAccent(image, accent)
			//BinMaskAccent(image, noise, accent)
			BinMask(image, noise)
			image.Save(fmt.Sprintf("data/videobin/%04d.png", i))
			//imageg := ImageLoad(fmt.Sprintf("data/test%d.png", i), 2.4)
			//BinAdaptiveFull(image)
			//image.Save(fmt.Sprintf("data/test%dadfull.png", i))
			//BinAdaptiveLocal(image)
			//BinMask4(imageg, noise, image)

			//image.Save(fmt.Sprintf("data/test%da.png", i))
			//imageg.Save(fmt.Sprintf("data/test%db.png", i))
		}
		bar.Finish()*/
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

	grayA := ImageLoad("data/video/0100.tif", -1)
	grayA.GammaCorrection(2.4)
	grayBlocksA, _, _ := ImageToBlocks(grayA)
	monoA := BinMask(grayA, noise)
	monoBlocksA, _, _ := ImageToBlocks(monoA)

	grayB := ImageLoad("data/video/0101.tif", -1)
	grayB.GammaCorrection(2.4)
	grayBlocksB, _, _ := ImageToBlocks(grayB)
	monoB := BinMask(grayB, noise)
	monoBlocksB, bw, bh := ImageToBlocks(monoB)
	res := EncodeFrame(grayBlocksB, monoBlocksB, grayBlocksA, 0.5)
	fmt.Println(len(res), len(grayBlocksB))
	fmt.Printf("Compression: %2.f %% (%d/%d)\n", float64(res.GetSize())/float64(grayB.Width*grayB.Height)*100, res.GetSize(), grayB.Width*grayB.Height)

	monoB.Save("data/result_1orig.png")

	resBlocks := res.Decode(monoBlocksA)
	resImage := BlocksToImage(resBlocks, bw, bh)
	resImage.Save("data/result_2compress.png")
}
