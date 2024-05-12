package main

import (
	com "cybercon/common"
	"fmt"
)

func main() {
	noise := com.ImageLoadRaw("data/noise1024x768.raw")
	fmt.Println("Generating sampler")
	inimg := com.ImageLoad("data/video/0410.tif", 2.4)
	sampler := com.GenerateSampler(inimg.Width, inimg.Height, 1000, 10, 100)
	for i := 440; i <= 450; i++ {
		fmt.Println(i)
		inimg := com.ImageLoad(fmt.Sprintf("data/video/%04d.tif", i), 2.4)
		outimg := com.EnhanceBWSampled(inimg, sampler)
		outimg = com.BinMask(outimg, noise)
		//outimg := com.BinMask(inimg, noise)
		//outimg := com.EnhanceBW(inimg, 1000, 10, 10)
		outimg.Save(fmt.Sprintf("data/vidno/%04d.png", i))
	}
	//inimg.Save("data/0410_1.png")

	//fmt.Println(sampler[0])
	//chans := com.ImageLoadDecompose("data/test1.png")
	//chans[0].Save("data/test1_cr.png")
	//chans[1].Save("data/test1_cg.png")
	//chans[2].Save("data/test1_cb.png")
	/*chans2 := []*com.ImageData{
		com.EnhanceBW(chans[0], 300, 10, 10),
		com.EnhanceBW(chans[1], 300, 10, 10),
		com.EnhanceBW(chans[2], 300, 10, 10),
	}
	chans2[0].Save("data/test1_er.png")
	chans2[1].Save("data/test1_eg.png")
	chans2[2].Save("data/test1_eb.png")
	comp := com.ImageCompose(chans2)
	comp.Save("data/test1_bw2.png")*/
	//noise := com.GenerateBlueNoise(1024, 768)
	//noise.SaveRaw("data/noise1024x768.raw")
	/*noise := com.ImageLoadRaw("data/noise693x924.raw")
	inimg := com.ImageLoad("data/test1.png", 2.4)
	//outimg := com.BinMask(inimg, noise)
	//outimg.Save("data/test1_1.png")
	outimg := com.EnhanceBW(inimg, 1000, 10, 10)
	outimg = com.ImageMerge(inimg, outimg, 0.5)
	outimg = com.BinMask(outimg, noise)

	outimg.Save("data/test1_2.png")*/
}
