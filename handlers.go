package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func getFromRepo(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	log.Println(file.Filename + ": updating binary ")
	setVer(file.Filename)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("art_root/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.String(http.StatusOK, file.Filename)
}

func getArtef(c echo.Context) error {
	return c.String(getConfig("artef"))
}

func getSquad(c echo.Context) error {
	return c.String(getConfig("squad"))
}

func setArtef(c echo.Context) error {
	return c.String(setConfig(c, "artef"))
}

func setSquad(c echo.Context) error {
	return c.String(setConfig(c, "squad"))
}
