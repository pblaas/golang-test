package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type Shoutbox struct {
	Name  string `form:"Username"`
	Shout string `form:"Shout"`
}

func Submitform(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Form["Username"] != nil {

		session, err := mgo.Dial("localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()

		//Optional. Switch the session to monotonic behavior
		session.SetMode(mgo.Monotonic, true)

		c := session.DB("mgodb").C("shouts")
		entry := &Shoutbox{
			Name:  r.FormValue("Username"),
			Shout: r.FormValue("Shout"),
		}
		err = c.Insert(entry)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	http.Handle("/input", http.FileServer(http.Dir("public/input")))
	http.HandleFunc("/submit", Submitform)
	http.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Println("Server started!")
	http.ListenAndServe(":3000", nil)
}
