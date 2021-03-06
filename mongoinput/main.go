package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type Shoutbox struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `form:"Username"`
	Shout      string        `form:"Shout"`
	Remoteaddr string
	Timestamp  time.Time
}

type Wordcount struct {
	Word  string `bson:"_id"`
	Value int    `bson:"value"`
}

func Submitform(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("mgodb").C("shouts")
	var shoutname string
	if r.FormValue("Shout") != "" {

		s := strings.Split(r.RemoteAddr, ":")
		ip := s[0]
		if r.FormValue("Username") == "" {
			shoutname = "Anonymous"
		} else {
			shoutname = r.FormValue("Username")
		}
		entry := &Shoutbox{
			Name:       shoutname,
			Shout:      r.FormValue("Shout"),
			Remoteaddr: ip,
			Timestamp:  time.Now(),
		}
		err = c.Insert(entry)
		if err != nil {
			log.Fatal(err)

		}
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func Mainpage(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("mgodb").C("shouts")

	query := c.Find(nil).Sort("-$natural")
	//query := c.Find(nil)
	var entries []Shoutbox
	if err := query.All(&entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func wcloudpage(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("mgodb").C("word_count")

	query := c.Find(nil)
	var words []Wordcount
	if err := query.All(&words); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "wcloud.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, words); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	//http.Handle("/input", http.FileServer(http.Dir("public/input")))
	http.HandleFunc("/submit", Submitform)
	//http.HandleFunc("/shouts", Mainpage)
	//http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, r.URL.Path[1:])
	//})
	http.HandleFunc("/wcloud", wcloudpage)
	http.HandleFunc("/", Mainpage)
	fmt.Println("Server started!")
	http.ListenAndServe(":3000", nil)
}
