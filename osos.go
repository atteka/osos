package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var session *mgo.Session
var noteStore = make(map[string]Note)

type (
	Note struct {
		Id    bson.ObjectId `bson:"_id,omitempty"`
		Email string        `json:"email"`
		Title string        `json:"title"`
		Time  string        `json:"time"`
	}

	DataStore struct {
		session *mgo.Session
	}
)

func (d *DataStore) Close() {
	d.session.Close()
}

func (d *DataStore) C(name string) *mgo.Collection {
	return d.session.DB("taskdb").C(name)
}

func NewDataStore() *DataStore {
	ds := &DataStore{
		session: session.Copy(),
	}
	return ds
}

//HTTP POST api/notes
func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	ds := NewDataStore()
	defer ds.Close()

	c := ds.C("notes")
	err = c.Insert(&note)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(note)

	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)

}

//HTTP GET api/notes
func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	ds := NewDataStore()
	defer ds.Close()
	c := ds.C("notes")

	iter := c.Find(nil).Iter()
	result := Note{}
	for iter.Next(&result) {
		notes = append(notes, result)
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//HTTP PUT api/notes/{id}
func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	k := vars["id"]
	var noteToUpdate Note
	err = json.NewDecoder(r.Body).Decode(&noteToUpdate)
	if err != nil {
		panic(err)
	}
	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
		noteStore[k] = noteToUpdate
	} else {
		log.Printf("could not find key")
	}
	w.WriteHeader(http.StatusNoContent)
}

//HTTP DELETE api/notes/{id}
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
	} else {
		log.Printf("Could not find key")
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}
