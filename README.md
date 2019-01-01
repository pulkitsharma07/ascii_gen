# ascii_gen
A script to generate basic colorfull ASCII art from Images.

The code has ample comments documenting it.

## Demo
![Demo](https://github.com/pulkitsharma07/pulkitsharma07.github.io/raw/development/assets/render1546308037772.gif)

## Working
The code maps each 8x8 window of pixels of the image to the character which is most similar to
that window.

The code has enough comments to get you going, I recommend starting out with `chars/chars.go` and then
move to the `main()` function in `main.go`.

## Running/Installing
* clone this repo
* `go get github.com/nfnt/resize`  (Required for resizing the image to suitable dimensions)
* `go run main.go /path/to/image.jpeg`

*PS: I am developing this mainly to learn Go, I do not have any considerable experience in writing Go (or ASCII generators for that matter). Please feel free to point out issues in the code, style, or anything in general !*
