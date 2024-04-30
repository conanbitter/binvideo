package main

func main() {
	image := ImageLoad("data/test2.png", -1)
	BinAdaptiveFull(image)
	image.Save("data/test2adfull.png")
}
