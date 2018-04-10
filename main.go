package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	rootPath = "/root/art_root/"
	port     = "3030"

	cYellow = color.New(color.FgYellow).SprintFunc()
	cRed    = color.New(color.FgRed).SprintFunc()
	cBold   = color.New(color.Bold).SprintFunc()
	cGreen  = color.New(color.FgGreen).SprintFunc()
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
	e.POST("/squad/update", setSquads)
	e.POST("/squad/update/:name", setSquad)

	// e.GET("/bin", upFile)
	e.Static("/bin", rootPath)

	e.Logger.Fatal(e.Start(":" + port))
}

func gitTrigger() {
	for _, a := range *artef {
		if a.Bin != "butler" {
			go func(a Artef) {
				log.Printf("[%s]:\t triggering git repo", cBold(a.Name))
				payload := url.Values{}
				payload.Add("token", a.Tok)
				payload.Add("ref", "master")
				_, err := http.PostForm(a.GitPath, payload)
				if err != nil {
					log.Println(err)
				}
			}(a)
		}
	}
}

func broadcastArtef(bin string) {
	t, _ := json.Marshal(artef)
	for _, s := range *squad {
		if s.Ips.V4 != "192.168.0.101" {
			needUpdate := false
			for _, ms := range s.Ms {
				if ms.Bin == bin {
					needUpdate = true
				}
			}
			if needUpdate {
				res, err := http.Post("http://"+s.Ips.V4+":9000/update/artef", "application/json", bytes.NewBuffer(t))
				if err != nil || res.StatusCode != 200 {
					log.Printf("[%s]: Failed broadcast to %s.\n", cRed(bin), cBold(s.Name))
				} else {
					log.Printf("[%s]: Sent artef to %s.\n", cGreen(bin), cBold(s.Name))
				}
			}
		}
	}
}

func broadcastSquad() {
	t, _ := json.Marshal(squad)
	for _, s := range *squad {
		for _, m := range s.Ms {
			if m.Bin == "gateway" {
				res, err := http.Post("http://"+s.Ips.V4+":"+m.Port+"/update", "application/json", bytes.NewBuffer(t))
				if res.StatusCode != 200 {
					log.Println("Failed update gateway on", s.Name, err)
					log.Printf("[%s]: Failed update gateway on.\n %s", cRed(s.Name), err)
				}
			}
		}
	}
}
