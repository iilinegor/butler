package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	dbHost   = "localhost"
	dbPass   = ""
	colArtef *mgo.Collection
	colSquad *mgo.Collection
	session  *mgo.Session
	artef    = new([]Artef)
	squad    = new([]Squad)
)

func initDB() {
	log.Println("db " + dbHost + " connection...")
	session, err := mgo.Dial("mongodb://" + dbHost + ":27017/test")
	if err != nil {
		log.Panic("Can't connect to DB: ", err)
	}

	log.Println("db initiation...")

	colArtef = session.DB("test").C("artef")
	colSquad = session.DB("test").C("squad")
	colArtef.Find(bson.M{}).All(artef)
	colSquad.Find(bson.M{}).All(squad)

	t, _ := json.Marshal(artef)
	log.Println(string(t))
	gitTrigger()
}

func setVer(bin string) {
	colArtef.Update(bson.M{"bin": bin}, bson.M{"$inc": bson.M{"ver": 1}})
	colArtef.Find(bson.M{}).All(artef)
}

func getConfig(target string) (int, string) {
	switch target {
	case "artef":
		t, _ := json.Marshal(artef)
		return http.StatusOK, string(t)

	case "squad":
		t, _ := json.Marshal(squad)
		return http.StatusOK, string(t)

	default:
		for _, s := range *squad {
			if s.Name == target {
				t, _ := json.Marshal(s)
				return http.StatusOK, string(t)
			}
		}
	}
	return http.StatusOK, " "
}

func setConfig(c echo.Context, target string) (int, string) {
	switch target {
	case "artef":
		tmp := new(Artef)
		if err := c.Bind(tmp); err != nil {
			log.Println(err)
		}

		for _, a := range *artef {
			if tmp.Name == a.Name {
				colArtef.Update(bson.M{"name": tmp.Name}, bson.M{"$set": tmp})
				t, _ := json.Marshal(tmp)
				return http.StatusOK, string(t)
			}
		}
		colArtef.Insert(tmp)
		t, _ := json.Marshal(tmp)
		return http.StatusOK, string(t)

	case "squad":
		tmp := new(Squad)
		if err := c.Bind(tmp); err != nil {
			log.Println(err)
		}

		for _, a := range *squad {
			if tmp.Name == a.Name {
				colSquad.Update(bson.M{"name": tmp.Name}, bson.M{"$set": tmp})
				t, _ := json.Marshal(tmp)
				return http.StatusOK, string(t)
			}
		}
		colSquad.Insert(tmp)
		t, _ := json.Marshal(tmp)
		return http.StatusOK, string(t)

	default:
		tmp := Squad{}
		if err := c.Bind(&tmp); err != nil {
			log.Println(err)
		}

		for _, s := range *squad {
			if s.Name == target {
				t, _ := json.Marshal(tmp)
				res, err := http.Post("http://"+s.Ips.V4+":9000/update/squad", "application/json", bytes.NewBuffer(t))
				if res.StatusCode != 200 {
					log.Printf("[%s]: Failed update squad on. \n %s", cRed(s.Name), err)
				} else {
					body, _ := ioutil.ReadAll(res.Body)
					json.Unmarshal(body, &tmp)

					colSquad.Update(bson.M{"name": s.Name}, bson.M{"$set": tmp})
					colSquad.Find(bson.M{}).All(squad)

					broadcastSquad()
					return http.StatusOK, string(body)
				}
			}
		}
	}

	return http.StatusBadRequest, "Something wrong..."
}

func uniqName() string {
	got := false

	for _, name := range Names {
		got = false
		for _, s := range *squad {
			if s.Name == name {
				got = true
			}
		}
		if !got {
			return name
		}
	}
	return "no free names"
}
