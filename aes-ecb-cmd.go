//go run main.go --image your_image.bmp --colors 256 --flip true --pix-width 16

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

	"github.com/spf13/cobra"
)

func main() {
	var colorsNeeded int
	var flip bool
	var pixWidth int
	var imageLocation string

	// Create a new root command
	var rootCmd = &cobra.Command{
		Use:   "aes-image-visualizer",
		Short: "Visualizes AES-encrypted images",
		Run: func(cmd *cobra.Command, args []string) {
			// Call the function to process the image
			processImage(imageLocation, colorsNeeded, flip, pixWidth)
		},
	}

	// Define flags for the command
	rootCmd.Flags().IntVarP(&colorsNeeded, "colors", "c", 254, "Number of colors needed in the palette")
	rootCmd.Flags().BoolVarP(&flip, "flip", "f", true, "Flip the image vertically")
	rootCmd.Flags().IntVarP(&pixWidth, "pix-width", "w", 16, "Width in bytes of the pixel data")
	rootCmd.Flags().StringVarP(&imageLocation, "image", "i", "aes.bmp", "Path to the input image")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func processImage(imageLocation string, colorsNeeded int, flip bool, pixWidth int) {
	blockSize := 16

	bytes, err := os.ReadFile(imageLocation)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	bytes = bytes[:len(bytes)-(len(bytes)%blockSize)]

	blocks := make([][]byte, 0)
	for i := 0; i < len(bytes); i += blockSize {
		block := bytes[i : i+blockSize]
		blocks = append(blocks, block)
	}

	blockCounts := make(map[string]int)
	for _, block := range blocks {
		blockKey := string(block)
		blockCounts[blockKey]++
	}

	type BlockFrequencifier struct {
		BlockKey string
		Count    int
	}

	blockFrequencies := make([]BlockFrequencifier, 0, len(blockCounts))
	for k, v := range blockCounts {
		blockFrequencies = append(blockFrequencies, BlockFrequencifier{BlockKey: k, Count: v})
	}
	sort.Slice(blockFrequencies, func(i, j int) bool {
		return blockFrequencies[i].Count > blockFrequencies[j].Count
	})

	colorMap := make(map[string]color.Color)
	palette := make([]color.Color, 0, colorsNeeded)
	palette = append(palette, color.White) // First color is white

	for i := 1; i < colorsNeeded-1; i++ {
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
	for i, bf := range blockFrequencies {
		var clr color.Color
		if i < colorsNeeded-1 {
			clr = palette[i]
		} else {
			clr = palette[len(palette)-1] // Assign black to less frequent blocks
		}
		colorMap[bf.BlockKey] = clr
	}

	// Map data blocks to color indices
	pixelData := make([]color.Color, 0, len(blocks)*pixWidth)
	for _, block := range blocks {
		blockKey := string(block)
		clr, exists := colorMap[blockKey]
		if !exists {
			clr = color.Black
		}
		for i := 0; i < blockSize/pixWidth; i++ {
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

	destFile := fmt.Sprintf("%s_aes.png", imageLocation)
	outFile, err := os.Create(destFile)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}
	defer outFile.Close()

	buf := bufio.NewWriter(outFile)
	if err := png.Encode(buf, img); err != nil {
		log.Fatalf("Encoder failed: %v", err)
	}
	if err := buf.Flush(); err != nil {
		log.Fatalf("File Writer Failed: %v", err)
	}

	fmt.Printf("Image saved to %s\n", destFile)
}

// example usage go run main.go --image your_image.bmp --colors 256 --flip true --pix-width 16
