package handlers

import (
	"github.com/labstack/echo"
	"image"
	"image/color"
	"strconv"
)

func Decrypt(c echo.Context) error {
	img, _ := decodeFileData(c)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	original := copyToNewImage(img, width, height, bounds)

	newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			oR, oG, oB := binary(original.RGBAAt(x, y))

			R, _ := strconv.ParseUint(oR[4:] + "0000", 2, 32)
			G, _ := strconv.ParseUint(oG[4:] + "0000", 2, 32)
			B, _ := strconv.ParseUint(oB[4:] + "0000", 2, 32)

			colorData := color.Color(color.RGBA{uint8(R), uint8(G), uint8(B), 255})
			newImage.Set(x, y, colorData)
		}
	}

	_ = saveImage("./images/de1.jpg", newImage)

	return c.File("images/de1.jpg")
}
