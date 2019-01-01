# ascii_gen
A script to generate basic colorful ASCII art from Images.
The code has comments documenting it and containing explanations of basic stuff. I am developing this as an exercise to learn Go.

## Demo
![Demo](https://github.com/pulkitsharma07/pulkitsharma07.github.io/raw/development/assets/render1546308037772.gif)

## Working
The code maps each 8x8 window of pixels of the image to the character which is most similar to
that window.

The code has enough comments to get you going, I recommend starting out with `chars/chars.go` and then
move to the `main()` function in `main.go`.

## Running/Installing
* clone this repo
* `go get github.com/nfnt/resize`  (Required fo resizing the image to suitable dimensions)
* `go run main.go /path/to/image.jpeg`
