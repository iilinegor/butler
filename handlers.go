package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func getFromRepo(c echo.Context) error {
	log.Println("Reciving started..")
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
	}

	log.Println(file.Filename + ": updating binary ")
	setVer(file.Filename)

	src, err := file.Open()
	if err != nil {
		log.Println(err)
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(rootPath + file.Filename)
	if err != nil {
		log.Println(err)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
	}

	go broadcastArtef()

	return c.String(http.StatusOK, file.Filename)
}

func regRunner(c echo.Context) error {
	tmpSquad := new(Squad)

	err := c.Bind(tmpSquad)
	if err != nil {
		return err
	}

	for _, s := range *squad {
		if s.Ips.V4 == tmpSquad.Ips.V4 {
			ts, _ := json.Marshal(s)
			return c.String(http.StatusBadRequest, string(ts))
		}
	}

	// Validating income config
	switch {
	case tmpSquad.Ips.V4 == "":
		return c.String(http.StatusBadRequest, tmpSquad.Ips.V4)

	case tmpSquad.Ips.V6 == "":
		return c.String(http.StatusBadRequest, tmpSquad.Ips.V4)
	}

	tmpSquad.Name = uniqName()

	colSquad.Insert(tmpSquad)
	*squad = append(*squad, *tmpSquad)
	colSquad.Find(bson.M{}).All(squad)

	t, err := json.Marshal(tmpSquad)

	return c.String(http.StatusOK, string(t))
}

func getArtef(c echo.Context) error {
	return c.String(getConfig("artef"))
}

func getSquads(c echo.Context) error {
	return c.String(getConfig("squad"))
}

func getSquad(c echo.Context) error {
	return c.String(getConfig(c.Param("name")))
}

func setArtef(c echo.Context) error {
	return c.String(setConfig(c, "artef"))
}

func setSquad(c echo.Context) error {
	return c.String(setConfig(c, "squad"))
}
