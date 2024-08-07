package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
)

var img_width int = 1920
var img_height int = 1080

var colour_black color.RGBA = color.RGBA{0, 0, 0, 0xff}
var colour_white color.RGBA = color.RGBA{0xff, 0xff, 0xff, 0xff}

func handle_error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	arguments := os.Args
	//arguments := []string{"", "decode", "outfile.png", "document.docx"}
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

		file, _ := os.Create(arguments[3])

		_, err := file.Write(output)
		if err != nil {
			fmt.Println(err)
		}

		file.Close()
	}
}
