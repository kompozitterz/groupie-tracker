package main

import (
	"fmt"
	groupie "groupie/utils"
	"net/http"
)

const PORT = ":8020"

func main() {
	fileServerCss := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServerCss))

	http.HandleFunc("/", groupie.MainPageHandler)
	http.HandleFunc("/band/", groupie.PageArtistHandler)
	//http.HandleFunc("/500", groupie.ErrorHandler)
	fmt.Println("access:http://localhost" + PORT)

	http.ListenAndServe(PORT, nil)
}
