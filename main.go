package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	rootPath = "./art_root"
)

func init() {
	initDB()
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.POST("/repo", getFromRepo)

	e.POST("/reg", regRunner)

	e.GET("/artef", getArtef)
	e.GET("/squad", getSquads)
	e.GET("/squad/:name", getSquad)

	e.POST("/artef/update", setArtef)
	e.POST("/squad/update", setSquad)

	// e.GET("/bin", upFile)
	e.Static("/bin", rootPath)

	e.Logger.Fatal(e.Start(":3030"))

}

func gitTrigger() {
	for _, a := range *artef {
		log.Println(a.Name + ": triggering git repo")
		payload := url.Values{}
		payload.Add("token", a.Tok)
		payload.Add("ref", "master")
		_, err := http.PostForm(a.GitPath, payload)
		if err != nil {
			log.Println(err)
		}
	}
}
