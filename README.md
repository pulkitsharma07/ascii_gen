# ascii_gen
Generate basic ASCII art from Images.

## Working
The code maps each 8x8 window of pixels of the image to the character which is most similar to
that window.

The code has enough comments to get you going, I recommend starting out with chars/chars.go !

## Demo
![Demo](https://github.com/pulkitsharma07/ascii_gen/raw/master/demo/render1545861845502.gif)
* You try the following things to yield better results:
    * Change number of characters per line.
    * Change the threshold for binarization.
    * Try inverting the image, which sometimes leads to better results (as seen in the demo).

## Running/Installing
* clone this repo
* `go get github.com/nfnt/resize`
* `go run main.go /path/to/image.jpeg`
