package main

import (
    "fmt"
    "github.com/nfnt/resize"
    "image"
    "image/jpeg"
    "math/bits"
    "os"

    "github.com/username/ascii_gen/chars" // Fixed the relative import path
)

// Defining a custom class for Color operations
type color struct {
    r, g, b uint32
}

func (a *color) Add(v *color) *color {
    a.r += v.r
    a.g += v.g
    a.b += v.b

    return a
}

func (a *color) retrofy() *color {
    a.Div(0x101)
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

func getCharWithColor(bestChar string, c *color) string {
    c.retrofy()
    return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.r, c.g, c.b, bestChar)
}

func getClosestChar(pattern uint64, c *color) string {
    maxDistance := 100
    var bestLetter string

    for k, v := range chars.CharMap {
        distance := bits.OnesCount64(v ^ pattern)
        if distance < maxDistance {
            bestLetter = k
            maxDistance = distance
        }
    }

    return getCharWithColor(bestLetter, c)
}

func getPackedFormOfWindow(img image.Image, winX, winY, w, h int, threshold uint32) uint64 {
    var pattern uint64 = 0
    cnt := 63

    for y := winY; y < winY+8 && y < h; y++ {
        for x := winX; x < winX+8 && x < w; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            if r+g+b >= threshold {
                pattern |= 1 << uint(cnt)
            }
            cnt--
        }
    }

    return pattern
}

func getMeanColorForWindow(img image.Image, winX, winY, w, h int) *color {
    colorAccum := &color{0, 0, 0}
    for y := winY; y < winY+8 && y < h; y++ {
        for x := winX; x < winX+8 && x < w; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            colorAccum.Add(&color{r, g, b})
        }
    }

    return colorAccum.Div(64)
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
    avgColor := getMeanColorForWindow(p.img, p.winX, p.winY, p.w, p.h)
    avgIntensity := uint32(avgColor.r + avgColor.g + avgColor.b/3)
    packedWindow := getPackedFormOfWindow(p.img, p.winX, p.winY, p.w, p.h, avgIntensity)
    char := getClosestChar(packedWindow, avgColor)
    inform <- windowProcessorResult{char, p.buffI, p.buffJ}
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
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("Error opening %s: %v", path, err)
        return
    }

    img_big, err := jpeg.Decode(f)
    bounds := img_big.Bounds()
    aspect_ratio := float64(bounds.Max.X) / float64(bounds.Max.Y)
    width := ascii_width * 8
    height := uint(float64(width) * 0.45 / aspect_ratio)
    img := resize.Resize(width, height, img_big, resize.Lanczos3)

    bounds = img.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y
    buffer := make([][]string, h/8+1)
    for i := range buffer {
        buffer[i] = make([]string, w/8+1)
    }

    buffI, buffJ := 0, 0
    inform := make(chan windowProcessorResult)
    done := make(chan bool)
    numProcessors := 0

    for winY := 0; winY < h; winY += 8 {
        buffJ = 0
        for winX := 0; winX < w; winX += 8 {
            processor := windowProcessor{img, winX, winY, w, h, buffI, buffJ}
            go processor.Run(inform)
            numProcessors++
            buffJ++
        }
        buffI++
    }

    go func() {
        resultsReceived := 0
        for {
            result, more := <-inform
            if more {
                resultsReceived++
                buffer[result.buffI][result.buffJ] = result.c
                if resultsReceived == numProcessors {
                    close(inform)
                }
            } else {
                break
            }
        }

        displayBuffer(buffer)
        done <- true
    }()

    <-done
}

func main() {
    var width uint = 150
    imgs := os.Args[1:]
    for _, img_path := range imgs {
        printImage(img_path, width)
    }
}
