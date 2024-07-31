package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

var img_width int = 200
var img_height int = 200

var colour_black color.RGBA = color.RGBA{0, 0, 0, 0xff}
var colour_white color.RGBA = color.RGBA{0xff, 0xff, 0xff, 0xff}

func handle_error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func binary_string(data []byte) string {
	result := ""

	for i := 0; i < len(data); i++ {
		str := strconv.FormatUint(uint64(data[i]), 2)

		if len(str) < 8 {
			fill_count := 8 - len(str)
			result += strings.Repeat("0", fill_count)
		}

		result += str
	}

	return result
}

func decode_binary(binary string) []byte {
	data := make([]byte, len(binary)/8)

	for i := 0; i < len(binary)-6; i += 8 {
		num, _ := strconv.ParseInt(binary[i:i+8], 2, 8)
		data[i/8] = byte(num)
	}

	return data
}

func encode_img(data []byte, filename string) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img_width, img_height}})

	// first pixel -> r, g, b - e.g. 0xaabbcc
	// r = 0xaa0000 >> 16
	// b = 0x00bb00 >> 8
	// g = 0xcc >> 0

	binary := binary_string(data)
	file_size := len(binary)

	b1 := (file_size & 0xff0000) >> 16
	b2 := (file_size & 0xff00) >> 8
	b3 := (file_size & 0xff) >> 0
	size_colour := color.RGBA{uint8(b1), uint8(b2), uint8(b3), 0xff}

	img.Set(0, 0, size_colour)

	x, y := 1, 0
	for i := 0; i < len(binary); i++ {
		if string(binary[i]) == "1" {
			img.Set(x, y, colour_white)
		} else if string(binary[i]) == "0" {
			img.Set(x, y, colour_black)
		}

		x++

		if x >= img_width {
			x = 0
			y++
		}
	}

	return img
}

func decode_img(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	image, _ := png.Decode(file)

	r, g, b, _ := image.At(0, 0).RGBA()

	file_size := int((uint8(r) << 16) + (uint8(g) << 8) + uint8(b))

	data := ""

	x, y := 1, 0
	for i := 0; i < file_size; i++ {
		pixel_colour, _, _, _ := image.At(x, y).RGBA()

		if pixel_colour>>8 == 0xff {
			data += "1"
		} else if pixel_colour>>8 == 0 {
			data += "0"
		}

		x++

		if x >= img_width {
			x = 0
			y++
		}
	}

	bytes := decode_binary(data)

	return string(bytes)
}

func main() {
	arguments := os.Args
	path := arguments[2]

	data, err := os.ReadFile(path)
	handle_error(err)

	if arguments[1] == "encode" {
		info, err := os.Stat(path)
		handle_error(err)

		filename := info.Name()

		frame := encode_img(data, filename)

		f, _ := os.Create("outfile.png")
		png.Encode(f, frame)
	} else if arguments[1] == "decode" {
		output := decode_img(path)
		fmt.Println(output)
	}
}