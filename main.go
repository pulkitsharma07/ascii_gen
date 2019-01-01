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

type windowProcessor struct {
	img                            image.Image
	winX, winY, w, h, buffI, buffJ int
}

type windowProcessorResult struct {
	c            string
	buffI, buffJ int
}

func (p windowProcessor) Run(inform chan windowProcessorResult) {
	//
	// First, we figure out the color is dominant in this 8x8 window.
	// Therefore, we can draw the character with that color in order to convey the color information.
	// These can be done in multiple ways: (mean/mode/median/maximum) value of R,G,Bs in the window
	//
	// Here, we are going to try out the mean color.
	avgColor := getMeanColorForWindow(p.img, p.winX, p.winY, p.w, p.h)

	// Not really 'intensity' in the proper sense. But some kind of value to indicate the "brightness"
	avgIntensity := uint32(avgColor.r + avgColor.g + avgColor.b/3)

	// Pack the current window into a 64 bit integer by performing binarization.
	// Details in the function definition.
	packedWindow := getPackedFormOfWindow(p.img, p.winX, p.winY, p.w, p.h, avgIntensity)

	// Figure out and print the character whose 8x8 representation is most similar to the current 8x8 window
	char := getClosestChar(packedWindow, avgColor)

	inform <- windowProcessorResult{char, p.buffI, p.buffJ}
}

// Converts R,G,B values back to 8 bits.
func (a *color) retrofy() *color {
	a.Div(0x101)
	return a
}

// Refer: https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
func getCharWithColor(bestChar string, c *color) string {
	c.retrofy()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.r, c.g, c.b, bestChar)
}

// Here <pattern> represents a 8x8 window of the image compressed into a 64 bit number
func getClosestChar(pattern uint64, c *color) string {
	maxDistance := 100
	var bestLetter string

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

	return getCharWithColor(bestLetter, c)
}

// TODO: Move this inside windowProcessor
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

// TODO: Move this inside windowProcessor
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

func displayBuffer(buffer [][]string) {
	for _, v := range buffer {
		for _, s := range v {
			fmt.Printf("%s", s)
		}
		fmt.Printf("\n")
	}
}

func printImage(path string, ascii_width uint) {
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
	width := ascii_width * 8

	// There interesting bit here is that, we are not preserving the aspect ratio of the image while
	// resizing. Specifically, we make the height half the value it is supposed to be wrt to the width.
	// This is done because, if we don't rescale the image, it will show up as squished in ASCII.
	//
	// Need to wrap my around as to why this value turned out to be 0.5
	// The generated ASCII image has somewhat similar aspect ratio (visually) to that of the source image.
	height := uint(float64(width) * 0.5 / aspect_ratio)

	img := resize.Resize(width, height, img_big, resize.Lanczos3)
	bounds = img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	// Create a 2D buffer of ASCII chars.
	// This is required the goroutines responsible for processing window can/will finish in a random
	// sequence.
	// Therefore, we can't just draw the chars on the screen at the end of each goroutine's execution.
	// We need a way to set characters in arbitrary location on the final image.
	// So, we use a buffer here.
	// Instead of drawing on the screen directly, the go routines will set the appropriate characters
	// in this buffer, and once  all go routines are done processing the image, we can finally draw the
	// buffer on the screen in a single go.
	buffer := make([][]string, h/8+1)
	for i := range buffer {
		buffer[i] = make([]string, w/8+1)
	}

	// Each go routine must know what coordinate in the buffer is it respondible for.
	// These variables are used to track that.
	buffI, buffJ := 0, 0

	// Common pattern to invoke multiple workers is to create 2 channels, one on which work is published,
	// and another one where the main go routine waits for all of the work which was generated to be completed.
	inform := make(chan windowProcessorResult)
	done := make(chan bool)
	numProcessors := 0

	// We need to scan (and draw) the image from left to right (and top to bottom)
	// Here winX, winY represents the top-left corner of the 8x8 window of the image, which will be
	// mapped to a single character.
	//
	// We move the window by 8 units, since we don't want to read the same pixels again.
	// Also, we move alone X axis first and then Y axis, because of the reasons stated earlier.
	for winY := 0; winY < h; winY += 8 {
		buffJ = 0
		for winX := 0; winX < w; winX += 8 {
			// Create a processor responsible for processing this 8x8 window
			processor := windowProcessor{img, winX, winY, w, h, buffI, buffJ}

			// Run the processor as a goroutine \m/
			//
			// <inform> is a channel where the processor will inform back with the appropriate ASCII
			// representation for this window.
			go processor.Run(inform)

			// Store count of processors for closing the channel
			numProcessors++
			buffJ++
		}
		buffI++
	}

	// Start a go routine which collects the results from all of the processor go routines
	// and writes to the buffer.
	// This ensures that there are no race-conditions while accessing the buffer.
	//
	// Refer:
	// https://gobyexample.com/closing-channels
	go func() {
		resultsReceived := 0
		for {
			result, more := <-inform
			if more {
				resultsReceived++
				buffer[result.buffI][result.buffJ] = result.c

				// Close channel once all information is received.
				// On the next call to `<- inform`, we will exit the loop, as <more> will not be true.
				if resultsReceived == numProcessors {
					close(inform)
				}
			} else {
				break
			}
		}

		// Display the buffer contents on the screen
		displayBuffer(buffer)

		// Push to done channel
		done <- true
	}()

	<-done
}

func main() {
	// Number of characters per line
	var width uint = 120

	imgs := os.Args[1:]

	for _, img_path := range imgs {
		// Display Image
		printImage(img_path, width)
	}
}
