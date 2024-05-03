package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

func write(file io.Writer, data interface{}) {
	binary.Write(file, binary.LittleEndian, data)
}

var magic = [3]byte{'B', 'V', 1}

func EncodeVideo(files []string, outfile string, options *EncodingOptions, fps float32) {
	noise := ImageLoadRaw("data/noise.raw")
	temp := ImageLoad(files[0], -1)
	bw, bh := GetSizeInBlocks(temp)
	width := temp.Width
	height := temp.Height
	curve := GetHilbertCurve(bw, bh)

	bar := progressbar.NewOptions(len(files),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionUseANSICodes(true))

	var err error
	file, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Magic
	file.Write(magic[:])

	// Header
	write(file, uint16(width))
	write(file, uint16(height))
	write(file, uint16(len(files)))
	write(file, fps)

	var lastMono []ImageBlock = nil
	var lastGray []ImageBlock = nil

	for i, filename := range files {
		bar.Set(i)

		grayImage := ImageLoad(filename, 2.4)
		if grayImage.Width != width || grayImage.Height != height {
			panic(fmt.Errorf("wrong frame size: %dx%d instead of %dx%d", grayImage.Width, grayImage.Height, width, height))
		}

		grayBlocks := ImageToBlocks(grayImage)
		gray := ApplyCurve(grayBlocks, curve)

		monoImage := BinMask(grayImage, noise)
		monoBlocks := ImageToBlocks(monoImage)
		mono := ApplyCurve(monoBlocks, curve)

		enc := EncodeFrame(gray, mono, lastGray, options)
		
		lastMono = enc.Decode(lastMono)
		lastGray = enc.DecodeGray(gray, lastGray)

		frameData := enc.Pack()
		length := uint32(len(frameData))
		if enc.IsClean() {
			length |= 1 << 31
		}
		write(file, length)
		file.Write(frameData)
	}
	bar.Set(len(files))
	bar.Finish()
}
