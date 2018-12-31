package main

import (
	chars "./chars"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	jpeg "image/jpeg"
	"math/bits"
	"os"
)

type color struct {
	r, g, b uint32
}

func (a *color) Add(v *color) *color {
	a.r += v.r
	a.g += v.g
	a.b += v.b

	return a
}

func (a *color) Div(number uint32) *color {
	a.r /= number
	a.g /= number
	a.b /= number

	return a
}

func (a *color) RGB() (uint32, uint32, uint32) {
	return a.r, a.g, a.b
}

// Map to memoize patterns which are already mapped to characters
var mem map[uint64]string

// <color>
// Color should be according to https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit
func printCharWithColor(bestChar string, c *color) {
	eight_bit_color := c.Div(0x101)
	fmt.Printf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", eight_bit_color.r, eight_bit_color.g, eight_bit_color.b, bestChar)
}

// Here <pattern> represents a 8x8 window of the image compressed into a 64 bit number
func printClosestChar(pattern uint64, c *color) {
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

	printCharWithColor(bestLetter, c)
}

func getPackedFormOfWindow(img image.Image, winX, winY, w, h int, threshold uint32) uint64 {
	// <pattern> will eventually be the packed form of the current 8 x 8 window
	// The packing will be similar to the one done in chars.CharMap
	// Refer the comment there for more details.
	var pattern uint64 = 0

	// Start assigning values in <pattern> from the MSB.
	// <cnt> indicates the bit to currently set/unset
	cnt := 63

	for y := winY; y < winY+8 && y < h; y++ {
		for x := winX; x < winX+8 && x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// We need to somewhow represent this RGB value as 0/1.
			// This is known as 'binarization', There can be multiple ways to do this.
			// Overall, it depends on some value (<threshold> here), which governs whether this pixel/bit
			// will be a 0 or a 1.
			if r+g+b >= threshold {
				pattern |= 1 << uint(cnt) // Set the <cnt>th bit in pattern as this pixel is above the threshold.
			}

			cnt-- // Move towards LSB
		}
	}

	return pattern
}

func getMeanColorForWindow(img image.Image, winX, winY, w, h int) *color {
	colorAccum := &color{0, 0, 0}

	// Just go through all the pixels in the current window
	// While ensuring that we don't cross the image bounds.
	//
	// Again, the order of scanning is important as we we want to store
	// the top most line on the most significant 8 bits.
	for y := winY; y < winY+8 && y < h; y++ {
		for x := winX; x < winX+8 && x < w; x++ {
			// Read the R,G,B values of the image at pixel <x>,<y>
			r, g, b, _ := img.At(x, y).RGBA()
			colorAccum.Add(&color{r, g, b})
		}
	}

	return colorAccum.Div(64)
}

func printImage(path string, width uint) {

	// Open the image present at <path>
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Some Error occured while opening %s: Erro: %v", path, err)
		return
	}

	// Try to read as JPEG
	img_big, err := jpeg.Decode(f)
	bounds := img_big.Bounds()
	aspect_ratio := float64(bounds.Max.X) / float64(bounds.Max.Y)

	// Resize according to width
	// The scripts converts each 8 x 8 block of image to 1 character.
	// Therefore, in order to write X characters per line, the image should be resized to 8*X.
	// Which maybe bigger/smaller than the original image.
	// The aspect ratio of the image will be maintained because 0 is passed as the second argument.
	img := resize.Resize(width*8, uint(float64(width)*3.7/aspect_ratio), img_big, resize.Lanczos3)
	bounds = img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	// We need to scan (and draw) the image from left to right (and top to bottom)
	// Here winX, winY represents the top-left corner of the 8x8 window of the image, which will be
	// mapped to a single character.
	//
	// We move the window by 8 units, since we don't want to read the same pixels again.
	// Also, we move alone X axis first and then Y axis, because of the reasons stated earlier.
	for winY := 0; winY < h; winY += 8 {
		for winX := 0; winX < w; winX += 8 {
			//
			// First, we figure out the color is dominant in this 8x8 window.
			// Therefore, we can draw the character with that color in order to convey the color information.
			// These can be done in multiple ways: (mean/mode/median/maximum) value of R,G,Bs in the window
			//
			// Here, we are going to try out the mean color.
			avgColor := getMeanColorForWindow(img, winX, winY, w, h)

			// Not really 'intensity' in the proper sense. But some kind of value to indicate the "brightness"
			avgIntensity := uint32(avgColor.r + avgColor.g + avgColor.b/3)

			// Pack the current window into a 64 bit integer by performing binarization.
			// Details in the function definition.
			packedWindow := getPackedFormOfWindow(img, winX, winY, w, h, avgIntensity)

			// Figure out and print the character whose 8x8 representation is most similar to the current 8x8 window
			printClosestChar(packedWindow, avgColor)
		}
		fmt.Println("")
	}
}

func main() {
	// Number of characters per line
	var width uint = 50

	// Initializing a map to memoize 8x8 window to character mappings
	mem = make(map[uint64]string)

	// forever
	for {

		// Clear Screen, https://stackoverflow.com/a/22892171
		print("\033[H\033[2J")

		imgs := os.Args[1:]

		for _, img_path := range imgs {
			// Display Image
			printImage(img_path, width)
		}

		// Some data
		fmt.Printf("Width: %d\n", width)
		fmt.Print("Increment/Decrement Width with u/i\n")
		fmt.Print("Enter q to exit.\n")

		// Wait for user to input something.
		var input string
		fmt.Scanln(&input)

		// Just handle different input cases
		switch input {
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
