package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		index,
	},
	Route{
		"MovieList",
		"GET",
		"/peliculas",
		movieList,
	},
	Route{
		"MovieShow",
		"GET",
		"/pelicula/{id}",
		movieShow,
	},
	Route{
		"MovieAdd",
		"POST",
		"/pelicula",
		movieAdd,
	},
	Route{
		"MovieUpdate",
		"PUT",
		"/pelicula/{id}",
		movieUpdate,
	},
	Route{
		"MovieDelete",
		"DELETE",
		"/pelicula/{id}",
		movieDelete,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Name(route.Name).
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandleFunc)
	}

	return router
}
