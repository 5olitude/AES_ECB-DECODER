# **AES ECB Decoder in Go**

An example of decoding an image encrypted using **AES** in **ECB** (Electronic Codebook) mode with **Go**. This demonstrates how AES-ECB works and explains why it's considered insecure for real-world usage.

## **What is AES with ECB Encryption?**

**AES** (Advanced Encryption Standard) is a symmetric block cipher that encrypts data in blocks.  
The AES block size is always **128 bits**, which is equal to **16 bytes**.

### **Process in Each Round:**

Each encryption round in AES involves the following steps:

1. **SubBytes**: A non-linear substitution step where each byte in the block is replaced using a substitution table (S-box).
2. **ShiftRows**: Rows of the state (the 16-byte block treated as a 4x4 matrix) are shifted.
3. **MixColumns**: A mixing operation combining the bytes in each column of the state (except in the final round).
4. **AddRoundKey**: The round key, derived from the encryption key, is XOR-ed with the block.

### **Ciphertext:**

After all the rounds are complete, a **16-byte (128-bit) block of ciphertext** is produced.  
If the plaintext exceeds 16 bytes, AES processes the data in **multiple 16-byte blocks**, generating 16 bytes of ciphertext for each block.

---

## **ECB (Electronic Codebook) Mode:**

ECB is one of the simplest block cipher modes of operation. In ECB mode:

- The plaintext is divided into **16-byte blocks**.
- Each block is **independently encrypted** using the same key.
- If the plaintext is not a multiple of 16 bytes, padding is applied to the final block.
- The resulting ciphertext blocks are concatenated to form the final encrypted message.

### **How AES-ECB Works:**

1. **Plaintext** is divided into 16-byte blocks.  
2. **Each block** is independently encrypted using AES and the same key.
3. The **resulting ciphertext** blocks are concatenated.

   ![Wikipedia ECB Penguin Image](https://github.com/5olitude/AES_ECB-DECODER/blob/7a2a45ad849d19343d4fe402bba6ac78275a88f3/Screens.png)
---

## **Real-World Example: ECB Penguin Attack**

The **ECB Penguin** is a famous example illustrating the weakness of AES-ECB mode. When an image (like the famous penguin image) is encrypted using ECB, **repeating patterns** in the plaintext are easily visible in the ciphertext. This happens because identical plaintext blocks are encrypted into identical ciphertext blocks, making the encrypted image resemble the original.  
Learn more about the **ECB Penguin** from this [GitHub link](https://github.com/robertdavidgraham/ecb-penguin).

You can also refer to the following Wikipedia illustration:



### **Other Real-World Attacks**:

One example is the **Adobe Password Database Leak**, where identical passwords produced identical ciphertexts. This allowed attackers to locate frequently reused passwords, leading to significant data breaches.

Here is a comic from **XKCD** that humorously illustrates this vulnerability:

![XKCD Encryptic Comic](https://imgs.xkcd.com/comics/encryptic_2x.png)

For a detailed explanation, refer to this [StackExchange discussion](https://crypto.stackexchange.com/questions/14487/can-someone-explain-the-ecb-penguin).


### CODE LOGIC IN GOLANG 

```go
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
        
	image_location := "aes.bmp"  //please specify the path of the image you wanna decrypt
	block_size := 16    //the block size in ecb mode as we learned 16 bytes


        //Redaing the image as bytes 
	bytes, err := os.ReadFile(image_location)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}


       //Adjusts the byte slice to remove any remaining bytes that do not form a complete block of size
       //block_size (i.e., 16 bytes). This ensures that only full blocks are processed.

	bytes = bytes[:len(bytes)-(len(bytes)%block_size)] 


      // Creates a slice of byte slices (blocks) to hold the individual 16-byte blocks. It loops through the bytes slice, slicing it into blocks and appending each          block to blocks.

	blocks := make([][]byte, 0)
	for i := 0; i < len(bytes); i += block_size {
		block := bytes[i : i+block_size]
		blocks = append(blocks, block)
	}

        // initializes a map block_counts to keep track of how many times each block appears. It iterates over the blocks, converting each block to a string key and incrementing its count in the map.

	block_counts := make(map[string]int)
	for _, block := range blocks {
		blockKey := string(block)
		block_counts[blockKey]++
	}

       // Defines a struct BlockFrequencifier to hold a block's key and its frequency count.
	type BlockFrequencifier struct {
		BlockKey string
		Count    int
	}

       // Initializes a slice block_frequencies to hold the frequency of each block. It iterates over block_counts, creating an instance of BlockFrequencifier for          // each block and appending it to block_frequencies

	block_frequencies := make([]BlockFrequencifier, 0, len(block_counts))
	for k, v := range block_counts {
		block_frequencies = append(block_frequencies, BlockFrequencifier{BlockKey: k, Count: v})
	}


        //Sorts the block_frequencies slice in descending order based on the Count field. This helps prioritize the most frequent blocks.
	sort.Slice(block_frequencies, func(i, j int) bool {
		return block_frequencies[i].Count > block_frequencies[j].Count
	})


       //Initializes a colorMap to associate block strings with colors. Also, it creates a palette slice to hold the colors, starting with white as the first color.
	colorMap := make(map[string]color.Color)
	palette := make([]color.Color, 0, colors_needed)
	palette = append(palette, color.White) // First color is white

       //This loop generates a color palette with the specified number of colors. Each color is generated with varying values for red, green, and blue components,         //ensuring diversity. Finally, black is added as the last color.
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







```
