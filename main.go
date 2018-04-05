package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	rootPath = "/root/art_root/"
	port     = "3030"
)

func init() {
	flag.StringVar(&port, "p", port, "port")
	flag.Parse()

	log.Println("Init butler...")

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

	e.Logger.Fatal(e.Start(":" + port))
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
