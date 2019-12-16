package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strconv"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "assets")

	e.POST("/st", func(c echo.Context) error {

		ori, hid, err := encryptFileData(c)

		bounds := ori.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y

		var original *image.RGBA
		original = image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(original, original.Bounds(), ori, bounds.Min, draw.Src)

		var hidden *image.RGBA
		hidden = image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(hidden, hidden.Bounds(), hid, bounds.Min, draw.Src)

		newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {

				oR, oG, oB := binary(original.RGBAAt(x, y))
				hR, hG, hB := binary(hidden.RGBAAt(x, y))

				R, _ := strconv.ParseUint(oR[0:4] + hR[0:4], 2, 64)
				G, _ := strconv.ParseUint(oG[0:4] + hG[0:4], 2, 64)
				B, _ := strconv.ParseUint(oB[0:4] + hB[0:4], 2, 64)
				//R, _ := strconv.ParseUint(hR[0:4] + "0000", 2, 64)
				//G, _ := strconv.ParseUint(hG[0:4] + "0000", 2, 64)
				//B, _ := strconv.ParseUint(hB[0:4] + "0000", 2, 64)

				colorData := color.Color(color.RGBA{uint8(R), uint8(G), uint8(B), 255})
				newImage.Set(x, y, colorData)

				if x==150 && y==150 {
					fmt.Println(oR)
					fmt.Println(oR[0:4] + hR[0:4])
					fmt.Println(R)
					fmt.Println(colorData)
				}
			}
		}

		dst, err := os.Create("./images/new1.jpg")
		if err != nil {
			return err
		}
		defer dst.Close()

		var opt jpeg.Options
		opt.Quality = 100

		_ = jpeg.Encode(dst, newImage, &opt)

		return c.File("images/new1.jpg")
	})

	e.POST("/de", func(c echo.Context) error {
		img, _ := decryptFileData(c)

		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y

		var original *image.RGBA
		original = image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(original, original.Bounds(), img, bounds.Min, draw.Src)

		newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				oR, oG, oB := binary(original.RGBAAt(x, y))

				R, _ := strconv.ParseUint(oR[4:] + "0000", 2, 64)
				G, _ := strconv.ParseUint(oG[4:] + "0000", 2, 64)
				B, _ := strconv.ParseUint(oB[4:] + "0000", 2, 64)

				colorData := color.Color(color.RGBA{uint8(R), uint8(G), uint8(B), 255})
				newImage.Set(x, y, colorData)

				if x==150 && y==150 {
					fmt.Println(oR)
					fmt.Println(R)
					fmt.Println(colorData)
				}
			}
		}

		dst, err := os.Create("./images/de1.jpg")
		if err != nil {
			return err
		}
		defer dst.Close()

		var opt jpeg.Options
		opt.Quality = 100

		_ = jpeg.Encode(dst, newImage, &opt)

		return c.File("images/de1.jpg")
	})

	e.Logger.Fatal(e.Start(":8888"))
}

func encryptFileData(c echo.Context) (image.Image, image.Image, error) {
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	file1, err := c.FormFile("image1")
	if err != nil {
		return nil, nil, err
	}

	file2, err := c.FormFile("image2")
	if err != nil {
		return nil, nil, err
	}

	src1, err := file1.Open()
	if err != nil {
		return nil, nil, err
	}
	defer src1.Close()

	src2, err := file2.Open()
	if err != nil {
		return nil, nil, err
	}
	defer src2.Close()

	img1, _, err := image.Decode(src1)
	if err != nil {
		return nil, nil, err
	}

	img2, _, err := image.Decode(src2)
	if err != nil {
		return nil, nil, err
	}

	return img1, img2, nil
}

func decryptFileData(c echo.Context) (image.Image, error) {
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := c.FormFile("image")
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
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