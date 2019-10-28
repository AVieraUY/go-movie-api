package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (this *Message) setStatus(data string) {
	this.Status = data
}

func (this *Message) setMessage(data string) {
	this.Message = data
}

var collection = getSession().DB("curso_go").C("movies")

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}

	return session
}

func responseMovie(w http.ResponseWriter, status int, result Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

func responseMovies(w http.ResponseWriter, status int, results []Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola desde mi servidor GO.")
}

func movieList(w http.ResponseWriter, r *http.Request) {
	var results []Movie
	err := collection.Find(nil).Sort("-_id").All(&results)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(results)
	}

	responseMovies(w, 200, results)
}

func movieShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idMovie := params["id"]

	if !bson.IsObjectIdHex(idMovie) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(idMovie)
	result := Movie{}
	err := collection.FindId(oid).One(&result)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	responseMovie(w, 200, result)
}

func movieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var movieData Movie
	err := decoder.Decode(&movieData)

	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}

	defer r.Body.Close()

	err = collection.Insert(movieData)

	if err != nil {
		panic(err)
	}

	responseMovie(w, 200, movieData)
}

func movieUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idMovie := params["id"]

	if !bson.IsObjectIdHex(idMovie) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(idMovie)

	decoder := json.NewDecoder(r.Body)

	var movieData Movie
	err := decoder.Decode(&movieData)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	defer r.Body.Close()

	document := bson.M{"_id": oid}
	change := bson.M{"$set": movieData}
	err = collection.Update(document, change)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if err != nil {
		w.WriteHeader(404)
		return
	}

	responseMovie(w, 200, movieData)
}

func movieDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idMovie := params["id"]

	if !bson.IsObjectIdHex(idMovie) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(idMovie)
	err := collection.RemoveId(oid)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	//result := Message{"success", "La pelicula con ID " + idMovie + " ha sido eliminada."}
	message := new(Message)
	message.setStatus("success")
	message.setMessage("La película ha sido eliminada correctamente.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(message)
}
