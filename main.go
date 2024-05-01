package main

import "fmt"

func main() {
	noise := ImageLoadRaw("data/noise.raw")
	for i := 1; i <= 9; i++ {
		fmt.Println(i)
		image := ImageLoad(fmt.Sprintf("data/test%d.png", i), -1) //2.4
		imageg := ImageLoad(fmt.Sprintf("data/test%d.png", i), 2.4)
		//BinAdaptiveFull(image)
		//image.Save(fmt.Sprintf("data/test%dadfull.png", i))
		BinAdaptiveLocal2(image)
		BinMask4(imageg, noise, image)

		image.Save(fmt.Sprintf("data/test%da.png", i))
		imageg.Save(fmt.Sprintf("data/test%db.png", i))
	}
	//noise := GenerateBlueNoise(1024, 768)
	//noise.SaveRaw("data/noise.raw")
	//noise.Save("data/noise.png")
}
