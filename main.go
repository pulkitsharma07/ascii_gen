package main

import (
	chars "./chars"
	"fmt"
	jpeg "image/jpeg"
	"math/bits"
	"os"
  "os/exec"
  "strconv"
)

var mem map[uint64]string

func getClosestChar(pattern uint64) {
	maxSimilarity := 0
	var bestLetter string

  if val, ok := mem[pattern]; ok {
    bestLetter = val
  } else {
    for k, v := range chars.CharMap {
      distance := bits.OnesCount64(^(v ^ pattern))

      if distance > maxSimilarity {
        bestLetter = k
        maxSimilarity = distance
      }
    }

    mem[pattern] = bestLetter
  }

	fmt.Printf("%s ", bestLetter)
}

func printImage(path string, threshold uint32) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Some Error occured while opening %s: Erro: %v", path, err)
		return
	}

	img, err := jpeg.Decode(f)

	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y


	for winY := 0; winY < h; winY += 8 {
		for winX := 0; winX < w; winX += 8 {
			var number uint64 = 0
			cnt := 63

			for y := winY; y < winY+8 && y < h; y++ {
				for x := winX; x < winX+8 && x < w; x++ {
					r, g, b, _ := img.At(x, y).RGBA()

					if r+g+b > threshold {
						number |= 1 << uint(cnt)
					}
					cnt--
				}
			}

			getClosestChar(number)
		}
		fmt.Println("")
	}
}

func main() {
  var threshold uint32 = 130000
  mem = make(map[uint64]string)

  for ;; {
    exec.Command("clear").Run()

    printImage(os.Args[1], threshold)

    fmt.Printf("Current Threshold: %d\n", threshold)
    fmt.Print("Change Threshold:")

    var input string
    fmt.Scanln(&input)

    if input == "q" {
      break
    }

    parsed, err := strconv.ParseInt(input, 10, 32)

    if err != nil {
      fmt.Printf("Invalid Input")
    } else {
      threshold = uint32(parsed)
    }
  }
}
