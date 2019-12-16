package handlers

import (
	"github.com/labstack/echo"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strconv"
)

func copyToNewImage(original image.Image, width int, height int, bounds image.Rectangle) *image.RGBA {
	var newImage *image.RGBA
	newImage = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newImage, newImage.Bounds(), original, bounds.Min, draw.Src)

	return newImage
}

func saveImage(path string, image *image.RGBA) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	var opt jpeg.Options
	opt.Quality = 100

	_ = jpeg.Encode(dst, image, &opt)

	return nil
}

func fileData(c echo.Context) (image.Image, image.Image, error) {
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	var img1, img2 image.Image

	if file1, _ := c.FormFile("image1"); file1 != nil{
		src1, err := file1.Open()
		if err != nil {
			return nil, nil, err
		}
		defer src1.Close()

		img1, _, err = image.Decode(src1)
		if err != nil {
			return nil, nil, err
		}
	}


	if file2, _ := c.FormFile("image2"); file2 != nil {
		src2, err := file2.Open()
		if err != nil {
			return nil, nil, err
		}
		defer src2.Close()

		img2, _, err = image.Decode(src2)
		if err != nil {
			return nil, nil, err
		}
	}



	return img1, img2, nil
}

func decodeFileData(c echo.Context) (image.Image, error) {
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	var img image.Image

	if file, _ := c.FormFile("decodeImage"); file != nil {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		img, _, err = image.Decode(src)
		if err != nil {
			return nil, err
		}
	}
	return img, nil
}

func eightDigits(s string) string {
	for len(s) != 8 {
		s = "0" + s
	}
	return s
}

func binary(c color.RGBA) (string, string, string) {
	r := eightDigits(strconv.FormatUint(uint64(c.R), 2))
	g := eightDigits(strconv.FormatUint(uint64(c.G), 2))
	b := eightDigits(strconv.FormatUint(uint64(c.B), 2))

	return r, g, b
}
