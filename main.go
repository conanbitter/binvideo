package main

func main() {
	image := ImageLoad("data/test1.png", -1)
	image.Save("data/test1gray.png")
}
