package handlers

import (
	"github.com/labstack/echo"
	"image"
	"image/color"
	"strconv"
)

func Encrypt(c echo.Context) error {

	ori, hid, _ := fileData(c)

	bounds := ori.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	original := copyToNewImage(ori, width, height, bounds)
	hidden := copyToNewImage(hid, width, height, bounds)

	newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			oR, oG, oB := binary(original.RGBAAt(x, y))
			hR, hG, hB := binary(hidden.RGBAAt(x, y))

			R, _ := strconv.ParseUint(oR[0:4] + hR[0:4], 2, 32)
			G, _ := strconv.ParseUint(oG[0:4] + hG[0:4], 2, 32)
			B, _ := strconv.ParseUint(oB[0:4] + hB[0:4], 2, 32)

			colorData := color.Color(color.RGBA{uint8(R), uint8(G), uint8(B), 255})
			newImage.Set(x, y, colorData)
		}
	}

	_ = saveImage("./images/new.jpg", newImage)

	return c.File("images/new.jpg")
}
