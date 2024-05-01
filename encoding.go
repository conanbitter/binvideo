package main

import "math"

const (
	ENC_RAW   byte = 0
	ENC_SKIP  byte = 1
	ENC_BLACK byte = 2
	ENC_WHITE byte = 3
)

type ImageBlock [16]float64

type EncodedBlock struct {
	BlockType byte
	Count     int
	Data      []ImageBlock
}

const MaxLength = 0b111111

var FullWhite = ImageBlock{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
var FullBlack = ImageBlock{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

type EncodingChain []*EncodedBlock

func ImageToBlocks(image *ImageData) []ImageBlock {
	bw, bh := GetSizeInBlocks(image)
	result := make([]ImageBlock, bw*bh)
	for y := 0; y < bh; y++ {
		for x := 0; x < bw; x++ {
			for by := 0; by < 4; by++ {
				for bx := 0; bx < 4; bx++ {
					bi := bx + by*4
					bl := x + y*bw
					imx := x*4 + bx
					if imx >= image.Width {
						imx = image.Width - 1
					}
					imy := y*4 + by
					if imy >= image.Height {
						imy = image.Height - 1
					}
					imi := imx + imy*image.Width
					result[bl][bi] = image.Data[imi]
				}
			}
		}
	}
	return result
}

func BlocksToImage(blocks []ImageBlock, blockWidth int, blockHeight int) *ImageData {
	result := &ImageData{
		Width:  blockWidth * 4,
		Height: blockHeight * 4,
		Data:   make([]float64, blockWidth*blockHeight*16),
	}
	for y := 0; y < blockHeight; y++ {
		for x := 0; x < blockWidth; x++ {
			for by := 0; by < 4; by++ {
				for bx := 0; bx < 4; bx++ {
					bi := bx + by*4
					bl := x + y*blockWidth
					imx := x*4 + bx
					imy := y*4 + by
					imi := imx + imy*result.Width
					result.Data[imi] = blocks[bl][bi]
				}
			}
		}
	}
	return result
}

func CompareBlocks(a *ImageBlock, b *ImageBlock) float64 {
	var acc float64 = 0.0
	for i := range a {
		dist := a[i] - b[i]
		acc += math.Sqrt(dist * dist)
	}
	return acc
}

func EncodeFrame(gray []ImageBlock, mono []ImageBlock, prevGray []ImageBlock, treshold float64) EncodingChain {
	var last *EncodedBlock = nil
	result := make([]*EncodedBlock, 0)

	for i := range gray {
		if last != nil && last.Count <= MaxLength {
			switch last.BlockType {
			case ENC_WHITE:
				if (CompareBlocks(&gray[i], &FullWhite) < treshold) || (CompareBlocks(&mono[i], &FullWhite) < treshold) {
					last.Count++
					continue
				}
			case ENC_BLACK:
				if (CompareBlocks(&gray[i], &FullBlack) < treshold) || (CompareBlocks(&mono[i], &FullBlack) < treshold) {
					last.Count++
					continue
				}
			case ENC_SKIP:
				if CompareBlocks(&gray[i], &prevGray[i]) < treshold {
					last.Count++
					continue
				}
			}
		}

		if prevGray != nil && CompareBlocks(&gray[i], &prevGray[i]) < treshold {
			last = &EncodedBlock{
				BlockType: ENC_SKIP,
				Count:     1,
				Data:      nil,
			}
			result = append(result, last)
			continue
		}
		if (CompareBlocks(&gray[i], &FullWhite) < treshold) || (CompareBlocks(&mono[i], &FullWhite) < treshold) {
			last = &EncodedBlock{
				BlockType: ENC_WHITE,
				Count:     1,
				Data:      nil,
			}
			result = append(result, last)
			continue
		}
		if (CompareBlocks(&gray[i], &FullBlack) < treshold) || (CompareBlocks(&mono[i], &FullBlack) < treshold) {
			last = &EncodedBlock{
				BlockType: ENC_BLACK,
				Count:     1,
				Data:      nil,
			}
			result = append(result, last)
			continue
		}

		if last != nil && last.BlockType == ENC_RAW {
			last.Count++
			last.Data = append(last.Data, mono[i])
		} else {
			last = &EncodedBlock{
				BlockType: ENC_RAW,
				Count:     1,
				Data:      []ImageBlock{mono[i]},
			}
			result = append(result, last)
		}
	}

	return result
}

func (chain EncodingChain) GetSize() int {
	size := 0
	for _, block := range chain {
		if block.BlockType == ENC_RAW {
			size += 1 + 2*block.Count
		} else {
			size++
		}
	}
	return size
}

func (chain EncodingChain) Decode(prevMono []ImageBlock) []ImageBlock {
	result := make([]ImageBlock, 0, len(prevMono))
	ind := 0
	for _, block := range chain {
		switch block.BlockType {
		case ENC_WHITE:
			for i := 0; i < block.Count; i++ {
				result = append(result, FullWhite)
			}
			ind += block.Count
		case ENC_BLACK:
			for i := 0; i < block.Count; i++ {
				result = append(result, FullBlack)
			}
			ind += block.Count
		case ENC_SKIP:
			for i := 0; i < block.Count; i++ {
				result = append(result, prevMono[ind])
				ind++
			}
		case ENC_RAW:
			result = append(result, block.Data...)
			ind += block.Count
		}
	}
	return result
}

func GetSizeInBlocks(image *ImageData) (int, int) {
	return int(math.Ceil(float64(image.Width) / 4)), int(math.Ceil(float64(image.Height) / 4))
}
