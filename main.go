package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
)

//Puerto donde escuchar
const PORT = "8998"

//TODO: Se debe parsear el enum de los comandos de tal forma que no sea necesario tener los valores harcodeados

//TODO: Por el momento suponemos que los pipes estan en /tmp
//Config por defecto, estaría bien usar swig para reusar la clase del daemon
//const CONFIG_FILE = "~/.config/player++/daemon.conf"

//Page contiene datos de la página
type Page struct {
	Title   string
	Request string
}

type Options struct {
	DaemonPipe  string
	ClientPipe  string
	MusicFolder string
	AutoStart   bool
	PidFile     string
	DbFile      string
}

var opt Options

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
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{2}, 0666)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, "index", "index", "")
}

func PrevHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{3}, 0666)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, "index", "index", "")
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{4}, 0666)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, "index", "index", "")
}

func GetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{16}, 0666)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(opt.ClientPipe)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	f.Close()
	fmt.Fprint(w, line)
}

func SetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{15}, 0666)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(opt.DaemonPipe, []byte(urlPart[2]+"\n"), 0666)
	if err != nil {
		panic(err)
	}
}

func Expand(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return strings.Replace(path, "~", usr.HomeDir, 1)
}

func LoadConfig() {
	cfg, err := ini.Load(Expand("~/.config/player++/daemon.conf"))
	if err != nil {
		panic(err)
	}

	cfg.NameMapper = ini.TitleUnderscore
	err = cfg.MapTo(&opt)
	if err != nil {
		panic(err)
	}

	log.Print(opt.DaemonPipe)

	opt.DaemonPipe = Expand(opt.DaemonPipe)
	opt.ClientPipe = Expand(opt.ClientPipe)
	opt.MusicFolder = Expand(opt.MusicFolder)
	opt.PidFile = Expand(opt.PidFile)
	opt.DbFile = Expand(opt.DbFile)

	log.Print(opt.DaemonPipe)
}

func main() {
	LoadConfig()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/next/", NextHandler)
	http.HandleFunc("/prev/", PrevHandler)
	http.HandleFunc("/pause/", PauseHandler)
	http.HandleFunc("/getvolume/", GetVolumeHandler)
	http.HandleFunc("/setvolume/", SetVolumeHandler)

	log.Print("Escuchando en el puerto " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
