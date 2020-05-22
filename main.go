package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

func main()  {
	router := mux.NewRouter()
	router.HandleFunc("/", route)
	router.HandleFunc("/mode", mode)
	router.HandleFunc("/mono", Mono)
	router.HandleFunc("/mosaic", Mosaic)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":" + os.Getenv("PORT"), router)
	//http.ListenAndServe(":8080", router)
}

func mode(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error: ", err)
	}
	if r.Form["mode"][0] == "mosaic" {
		t, _ := template.ParseFiles("templates/studio1.html")
		t.Execute(w, nil)
	} else if r.Form["mode"][0] == "mono" {
		t, _ := template.ParseFiles("templates/studio2.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("templates/home.html")
		t.Execute(w, nil)
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/home.html")
	t.Execute(w, nil)
}