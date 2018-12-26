package main

import (
	chars "./chars"
	"fmt"
	"github.com/nfnt/resize"
	jpeg "image/jpeg"
	"math/bits"
	"os"
	"os/exec"
)

// Map to memoize patterns which are already mapped to characters
var mem map[uint64]string

// Here <pattern> represents a 8x8 window of the image compressed into a 64 bit number
func printClosestChar(pattern uint64) {
	maxDistance := 100
	var bestLetter string

	if val, ok := mem[pattern]; ok {
		bestLetter = val
	} else {
		// Go through each character Mapping we have in chars.CharMap
		for k, v := range chars.CharMap {

			// Count the number of bits which are different in the pattern and the character
			// This count represents dissimilar these 2 8x8 images are.
			// Remember both are actually 8x8 images/patterns packed in 64 bit numbers.
			// Here we take the XOR of these two numbers, which gives us count of bits which are different.
			distance := bits.OnesCount64(v ^ pattern)

			// We need to store the character which is the most similar.
			// i.e. having the least number of different bits between it and the pattern.
			if distance < maxDistance {
				bestLetter = k
				maxDistance = distance
			}
		}

		// Memoize, so we can reuse the same letter (instead of iterating through chars.CharMap again)
		mem[pattern] = bestLetter
	}

	// Voila !
	fmt.Printf("%s ", bestLetter)
}

func printImage(path string, threshold uint32, width uint) {

	// Open the image present at <path>
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Some Error occured while opening %s: Erro: %v", path, err)
		return
	}

	// Try to read as JPEG
	img_big, err := jpeg.Decode(f)

	// Resize according to width
	// The scripts converts each 8 x 8 block of image to 1 character.
	// Therefore, in order to write X characters per line, the image should be resized to 8*X.
	// Which maybe bigger/smaller than the original image.
	// The aspect ratio of the image will be maintained because 0 is passed as the second argument.
	img := resize.Resize(width*8, 0, img_big, resize.Lanczos3)
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	// We need to scan (and draw) the image from left to right (and top to bottom)
	// Here winX, winY represents the top-left corner of the 8x8 window of the image, which will be
	// mapped to a single character.
	//
	// We move the window by 8 units, since we don't want to read the same pixels again.
	// Also, we move alone X axis first and then Y axis, because of the reasons stated earlier.
	for winY := 0; winY < h; winY += 8 {
		for winX := 0; winX < w; winX += 8 {

			// <pattern> will eventually be the packed form of the current 8 x 8 window
			// The packing will be similar to the one done in chars.CharMap
			// Refer the comment there for more details.
			var pattern uint64 = 0

			cnt := 63 // Start assigning values in <pattern> from the MSB.

			// Just go through all the pixels in the current window
			// While ensuring that we don't cross the image bounds.
			//
			// Again, the order of scanning is important as we we want to store
			// the top most line on the most significant 8 bits.
			for y := winY; y < winY+8 && y < h; y++ {
				for x := winX; x < winX+8 && x < w; x++ {

					// Read the R,G,B values of the image at pixel <x>,<y>
					r, g, b, _ := img.At(x, y).RGBA()

					// We need to somewhow represent this RGB values as 0/1
					// The threshold governs that.
					// For now, just adding R,G,B and comparing with the threshold.
					// The image can be inverted (in colors) by inverting this condition
					if r+g+b > threshold {
						pattern |= 1 << uint(cnt) // Set the <cnt>th bit in pattern as this pixel is above the threshold.
					}

					cnt-- // Move towards LSB
				}
			}

			// Figure out and print the character whose 8x8 representation is most similar to the current 8x8 window
			// (which is packed inside <pattern>)
			printClosestChar(pattern)
		}
		fmt.Println("")
	}
}

func main() {
	// Some random numbers for thresholding
	var threshold uint32 = 130000
	var step_size uint32 = 10000

	// Number of characters per line
	var width uint = 50

	// Initializing a map to memoize 8x8 window to character mappings
	mem = make(map[uint64]string)

	// forever
	for {

		// clear screen, so the user doesn't have to scroll
		exec.Command("clear").Run()

		// Display Image
		printImage(os.Args[1], threshold, width)

		// Some data
		fmt.Printf("Threshold: %d\tThreshold-StepSize: %d\tWidth: %d\n", threshold, step_size, width)
		fmt.Print("Increment/Decrement Threshold with j/k\n")
		fmt.Print("Increment/Decrement Threshold-StepSize with h/l\n")
		fmt.Print("Increment/Decrement Width with u/i\n")

		// Wait for user to input something.
		var input string
		fmt.Scanln(&input)

		// Just handle different input cases
		switch input {
		case "j":
			{
				threshold -= step_size
			}
		case "k":
			{
				threshold += step_size
			}
		case "h":
			{
				step_size -= 100
			}
		case "l":
			{
				step_size += 100
			}
		case "u":
			{
				width -= 2
			}
		case "i":
			{
				width += 2
			}
		case "q":
			{
				os.Exit(0)
			}
		}
	}
}
