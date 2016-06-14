package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

//Puerto donde escuchar
const PORT = "8998"

//TODO: Por el momento suponemos que los pipes estan en /tmp
//Config por defecto, estaría bien usar swig para reusar la clase del daemon
//const CONFIG_FILE = "~/.config/player++/daemon.conf"

//Page contiene datos de la página
type Page struct {
	Title   string
	Request string
}

const templatesPath = "./templates/"

var templates = template.Must(template.ParseGlob(templatesPath + "[a-z]*"))

//RenderTemplate renderiza la template con el texto adecuado
func RenderTemplate(w http.ResponseWriter, tmpl string, title string, r interface{}) {
	var err error
	switch request := r.(type) {
	case string:
		err = templates.ExecuteTemplate(w, tmpl, &Page{title, request})
	case error:
		//TODO: Si es un error se imprima en otro formato
		err = templates.ExecuteTemplate(w, tmpl, &Page{title, request.Error()})
	default:
		log.Print("Strange error ocurred")
	}

	if err != nil {
		log.Print("ExecuteTemplate: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RootHandler muestra el index
func RootHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index", "index", "")
}

func NextHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile("/tmp/dplayer++", []byte{2}, 0666)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, "index", "index", "")
}

func PrevHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile("/tmp/dplayer++", []byte{3}, 0666)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, "index", "index", "")
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile("/tmp/dplayer++", []byte{4}, 0666)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, "index", "index", "")
}

func main() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/next/", NextHandler)
	http.HandleFunc("/prev/", PrevHandler)
	http.HandleFunc("/pause/", PauseHandler)

	log.Print("Escuchando en el puerto " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
