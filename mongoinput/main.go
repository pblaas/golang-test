package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

type Shoutbox struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `form:"Username"`
	Shout     string        `form:"Shout"`
	Timestamp time.Time
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

	if r.Form["Username"] != nil {

		entry := &Shoutbox{
			Name:      r.FormValue("Username"),
			Shout:     r.FormValue("Shout"),
			Timestamp: time.Now(),
		}
		err = c.Insert(entry)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.Redirect(w, r, "/shouts", http.StatusTemporaryRedirect)
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

	query := c.Find(nil)
	//sh := Shoutbox{"Jeremy Saenz", "shout shout shout"}
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

func main() {
	//http.Handle("/input", http.FileServer(http.Dir("public/input")))
	http.HandleFunc("/submit", Submitform)
	http.HandleFunc("/shouts", Mainpage)
	http.Handle("/", http.FileServer(http.Dir("public/")))
	fmt.Println("Server started!")
	http.ListenAndServe(":3000", nil)
}
