package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Note struct {
	Email string `json:"email"`
	Title string `json:"title"`
	Time string  `json:"time"`
}

var noteStore = make(map[string]Note)
var id int

//HTTP POST api/notes
func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note
	fmt.Println(note)
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
	for _, v := range noteStore {
		notes = append(notes, v)
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
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr: ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}
