package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
)

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
		num, _ := strconv.ParseUint(binary[i:i+8], 2, 8)
		data[i/8] = byte(num)
	}

	return data
}

func encode_img(data []byte, filename string) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img_width, img_height}})

	binary := binary_string(data)

	x, y := 0, 0
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

func decode_img(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	image, _ := png.Decode(file)

	data := ""

	x, y := 0, 0

	decoding := true
	for decoding {
		pixel_colour, _, _, alpha := image.At(x, y).RGBA()

		if alpha <= 0 {
			fmt.Println(x, y)
			break
		}

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

	return bytes
}
