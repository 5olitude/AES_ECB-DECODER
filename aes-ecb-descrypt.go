package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
)

func main() {
	colors_needed := 254
	flip := true
	pix_width := 16

	image_location := "aes.bmp"
	block_size := 16

	bytes, err := os.ReadFile(image_location)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	bytes = bytes[:len(bytes)-(len(bytes)%block_size)]

	blocks := make([][]byte, 0)
	for i := 0; i < len(bytes); i += block_size {
		block := bytes[i : i+block_size]
		blocks = append(blocks, block)
	}

	block_counts := make(map[string]int)
	for _, block := range blocks {
		blockKey := string(block)
		block_counts[blockKey]++
	}

	type BlockFrequencifier struct {
		BlockKey string
		Count    int
	}

	block_frequencies := make([]BlockFrequencifier, 0, len(block_counts))
	for k, v := range block_counts {
		block_frequencies = append(block_frequencies, BlockFrequencifier{BlockKey: k, Count: v})
	}
	sort.Slice(block_frequencies, func(i, j int) bool {
		return block_frequencies[i].Count > block_frequencies[j].Count
	})

	colorMap := make(map[string]color.Color)
	palette := make([]color.Color, 0, colors_needed)
	palette = append(palette, color.White) // First color is white

	for i := 1; i < colors_needed-1; i++ {

		clr := color.RGBA{
			R: uint8((i * 50) % 255),
			G: uint8((i * 80) % 255),
			B: uint8((i * 110) % 255),
			A: 255,
		}
		palette = append(palette, clr)
	}
	palette = append(palette, color.Black) // Last color is black

	// Map blocks to colors
	for i, bf := range block_frequencies {
		var clr color.Color
		if i < colors_needed-1 {
			clr = palette[i]
		} else {
			clr = palette[len(palette)-1] // Assign black to less frequent blocks
		}
		colorMap[bf.BlockKey] = clr
	}

	// Map data blocks to color indices
	pixelData := make([]color.Color, 0, len(blocks)*pix_width)
	for _, block := range blocks {
		blockKey := string(block)
		clr, exists := colorMap[blockKey]
		if !exists {
			clr = color.Black
		}
		for i := 0; i < block_size/pix_width; i++ {
			pixelData = append(pixelData, clr)
		}
	}

	// Calculate image dimensions
	totalPixels := len(pixelData)
	width := int(math.Sqrt(float64(totalPixels)))
	height := totalPixels / width

	if totalPixels%width != 0 {
		height++
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for idx, clr := range pixelData {
		x := idx % width
		y := idx / width
		if flip {
			y = height - y - 1 // Correctly calculate the flipped y-coordinate
		}
		if y >= height {
			break
		}
		img.Set(x, y, clr)
	}

	dest_file := fmt.Sprintf("%s_aes.png", image_location)
	outFile, err := os.Create(dest_file)
	if err != nil {
		log.Fatalf("Failed to genearte image: %v", err)
	}
	defer outFile.Close()

	buf := bufio.NewWriter(outFile)
	if err := png.Encode(buf, img); err != nil {
		log.Fatalf("Encoder failed: %v", err)
	}
	if err := buf.Flush(); err != nil {
		log.Fatalf("File Witer Failed: %v", err)
	}

	fmt.Printf("Image saved to %s\n", dest_file)
}
