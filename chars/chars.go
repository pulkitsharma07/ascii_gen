package chars

/*
  This Map is created from the C source code on:
  https://opengameart.org/content/8x8-ascii-bitmap-font-with-c-source
  .
  .
  and some awk magic....

  Details:
  Each character below is represented by a 8 x 8 matrix of bits.

  Think of this as an 8x8 black and white image. A pixel the bit present at that position
  is set 0, and it is white if the bit is set to 1.

  For example the character 'P' can be represented as:

  01111000
  01000100
  01111000
  01000000
  01000000    <- Notice how all the 1's form the shape of P
  01000000
  00000000
  00000000


  Now, since we have 8x8 pixels (bits) , we can "pack" these in a 64 bit number.
  These can be packed in multiple ways, like storing row1, then row2.. so on, or
  you can store column1, column2, etc..

  The C-code linked above has stored the bits by each row, starting from the top most row.

  After "packing", we'll get the following 64 bit number:
  0111100001000100011110000100000001000000010000000000000000000000

  Notice, as we move from MSB to LSB we are going through the image from left to right (and top to bottom).
  Converting this 64 bit number to decimal, we get: 8666183800318853120

  Again, the code mentioned above has chosen to represent this number in hexadecimal notation instead.
  Therefore, 'P' can be represented as 0x7844784040400000

  It is quite common to represent series of bits as hexadecimal numbers.
  One reason can be that they convey more information with less number of characters.
*/
var CharMap map[string]uint64 = map[string]uint64{
	" ":  0x0,
	"!":  0x808080800080000,
	"\"": 0x2828000000000000,
	"#":  0x287C287C280000,
	"$":  0x81E281C0A3C0800,
	"%":  0x6094681629060000,
	"&":  0x1C20201926190000,
	"'":  0x808000000000000,
	"(":  0x810202010080000,
	")":  0x1008040408100000,
	"*":  0x2A1C3E1C2A000000,
	"+":  0x8083E08080000,
	",":  0x81000,
	"-":  0x3C00000000,
	".":  0x80000,
	"/":  0x204081020400000,
	"0":  0x1824424224180000,
	"1":  0x8180808081C0000,
	"2":  0x3C420418207E0000,
	"3":  0x3C420418423C0000,
	"4":  0x81828487C080000,
	"5":  0x7E407C02423C0000,
	"6":  0x3C407C42423C0000,
	"7":  0x7E04081020400000,
	"8":  0x3C423C42423C0000,
	"9":  0x3C42423E023C0000,
	":":  0x80000080000,
	";":  0x80000081000,
	"<":  0x6186018060000,
	"=":  0x7E007E000000,
	">":  0x60180618600000,
	"?":  0x3844041800100000,
	"@":  0x3C449C945C201C,
	"A":  0x1818243C42420000,
	"B":  0x7844784444780000,
	"C":  0x3844808044380000,
	"D":  0x7844444444780000,
	"E":  0x7C407840407C0000,
	"F":  0x7C40784040400000,
	"G":  0x3844809C44380000,
	"H":  0x42427E4242420000,
	"I":  0x3E080808083E0000,
	"J":  0x1C04040444380000,
	"K":  0x4448507048440000,
	"L":  0x40404040407E0000,
	"M":  0x4163554941410000,
	"N":  0x4262524A46420000,
	"O":  0x1C222222221C0000,
	"P":  0x7844784040400000,
	"Q":  0x1C222222221C0200,
	"R":  0x7844785048440000,
	"S":  0x1C22100C221C0000,
	"T":  0x7F08080808080000,
	"U":  0x42424242423C0000,
	"V":  0x8142422424180000,
	"W":  0x4141495563410000,
	"X":  0x4224181824420000,
	"Y":  0x4122140808080000,
	"Z":  0x7E040810207E0000,
	"[":  0x3820202020380000,
	"\\": 0x4020100804020000,
	"]":  0x3808080808380000,
	"^":  0x1028000000000000,
	"_":  0x7E0000,
	"`":  0x1008000000000000,
	"a":  0x3C023E463A0000,
	"b":  0x40407C42625C0000,
	"c":  0x1C20201C0000,
	"d":  0x2023E42463A0000,
	"e":  0x3C427E403C0000,
	"f":  0x18103810100000,
	"g":  0x344C44340438,
	"h":  0x2020382424240000,
	"i":  0x800080808080000,
	"j":  0x800180808080870,
	"k":  0x20202428302C0000,
	"l":  0x1010101010180000,
	"m":  0x665A42420000,
	"n":  0x2E3222220000,
	"o":  0x3C42423C0000,
	"p":  0x5C62427C4040,
	"q":  0x3A46423E0202,
	"r":  0x2C3220200000,
	"s":  0x1C201804380000,
	"t":  0x103C1010180000,
	"u":  0x2222261A0000,
	"v":  0x424224180000,
	"w":  0x81815A660000,
	"x":  0x422418660000,
	"y":  0x422214081060,
	"z":  0x3C08103C0000,
	"{":  0x1C103030101C0000,
	"|":  0x808080808080800,
	"}":  0x38080C0C08380000,
	"~":  0x324C000000,
}
