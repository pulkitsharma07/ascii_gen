# ascii_gen
A fast script to generate basic colorfull ASCII art from Images.

## Working
The code maps each 8x8 window of pixels of the image to the character which is most similar to
that window.

The code has enough comments to get you going, I recommend starting out with chars/chars.go !

## Demo
![Demo](https://github.com/pulkitsharma07/ascii_gen/raw/master/demo/render1545861845502.gif)


## Running/Installing
* clone this repo
* `go get github.com/nfnt/resize`  (Required fo resizing the image to suitable dimensions)
* `go run main.go /path/to/image.jpeg`
