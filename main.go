package main

import (
	"github.com/Zett-8/steganography/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "assets")

	e.POST("/enc", handlers.Encrypt)

	e.POST("/dec", handlers.Decrypt)

	e.Logger.Fatal(e.Start(":8888"))
}

