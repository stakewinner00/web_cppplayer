package main

// #include <cppplayer/commands.h>
import "C"

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

//TODO: Se debe parsear el enum de los comandos de tal forma que no sea necesario tener los valores harcodeados

//Puerto donde escuchar
const PORT = "8998"
const CONFIG_FILE = "~/.config/player++/daemon.conf"

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func pipeReadLine() string {
	f, err := os.Open(opt.ClientPipe)
	check(err)

	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	check(err)

	f.Close()

	return line
}

func Expand(path string) string {
	usr, err := user.Current()
	check(err)

	return strings.Replace(path, "~", usr.HomeDir, 1)
}

func LoadConfig() {
	cfg, err := ini.Load(Expand(CONFIG_FILE))
	check(err)

	cfg.NameMapper = ini.TitleUnderscore
	err = cfg.MapTo(&opt)
	check(err)

	log.Print(opt.DaemonPipe)

	opt.DaemonPipe = Expand(opt.DaemonPipe)
	opt.ClientPipe = Expand(opt.ClientPipe)
	opt.MusicFolder = Expand(opt.MusicFolder)
	opt.PidFile = Expand(opt.PidFile)
	opt.DbFile = Expand(opt.DbFile)

	log.Print(opt.DaemonPipe)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index", "index", "")
}

func NextHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.NEXT}, 0666)
	check(err)

	RenderTemplate(w, "index", "index", "")
}

func PrevHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.BACK}, 0666)
	check(err)

	RenderTemplate(w, "index", "index", "")
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.PAUSE}, 0666)
	check(err)

	RenderTemplate(w, "index", "index", "")
}

func GetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.VOLUME_GET}, 0666)
	check(err)

	line := pipeReadLine()

	fmt.Fprint(w, line)
}

func GetTitleHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.GET_TITLE}, 0666)
	check(err)

	line := pipeReadLine()

	fmt.Fprint(w, line)
}

func GetArtistHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.GET_ARTIST}, 0666)
	check(err)

	line := pipeReadLine()

	fmt.Fprint(w, line)
}

func SetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	err := ioutil.WriteFile(opt.DaemonPipe, []byte(string(C.VOLUME_SET)+urlPart[2]+"\n"), 0666)
	check(err)
}

func GetRemainingHandler(w http.ResponseWriter, r *http.Request) {
	err := ioutil.WriteFile(opt.DaemonPipe, []byte{C.TIME_GET_REMAINING}, 0666)
	check(err)

	line := pipeReadLine()

	fmt.Fprint(w, line)
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
	http.HandleFunc("/gettitle/", GetTitleHandler)
	http.HandleFunc("/getartist/", GetArtistHandler)
	http.HandleFunc("/getremaining/", GetRemainingHandler)

	log.Print("Escuchando en el puerto " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
