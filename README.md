# ascii_gen
Generate basic ASCII art from Images.

# Working
The code maps each 8 x 8 window of pixels of the image to the character which is most similar to
that window.

The code has enough comments to get you going !
I recommend starting with chars/chars.go

# Demo
[![asciicast](https://asciinema.org/a/wM66IFO9ZaIVapOX3H5sDJkbt.svg)](https://asciinema.org/a/wM66IFO9ZaIVapOX3H5sDJkbt)

# Running
* clone this repo
* `go get github.com/nfnt/resize`
* `go run main.go /path/to/image.jpeg`

