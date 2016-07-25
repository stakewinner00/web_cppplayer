package main

// #include <cppplayer/commands.h>
import "C"

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"html/template"
  "io"
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

type fifo struct {
	DaemonPipe string
	ClientPipe string
}

type Options struct {
	MusicFolder string
	AutoStart   bool
	PidFile     string
	DbFile      string
	fifo
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
  log.Print("Reading line")

	f, err := os.OpenFile(opt.fifo.ClientPipe, os.O_RDONLY, os.ModeNamedPipe)
	check(err)

	reader := bufio.NewReader(f)

	line, err := reader.ReadString('\n')
	for err == io.EOF {
		line, err = reader.ReadString('\n')
	}
	check(err)

	f.Close()
  log.Print("Line readed")

	return line
}

func pipeWrite(s []byte) {
  log.Print("Writing line")
	err := ioutil.WriteFile(opt.fifo.DaemonPipe, s, 0666)
	check(err)
  log.Print("Line writed")
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

	err = cfg.Section("fifo").MapTo(&opt.fifo)
	check(err)

	log.Print(opt.DaemonPipe)

	opt.fifo.DaemonPipe = Expand(opt.fifo.DaemonPipe)
	opt.fifo.ClientPipe = Expand(opt.fifo.ClientPipe)
	opt.MusicFolder = Expand(opt.MusicFolder)
	opt.PidFile = Expand(opt.PidFile)
	opt.DbFile = Expand(opt.DbFile)

	log.Print(opt.fifo.DaemonPipe)
	log.Print(opt.fifo.ClientPipe)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index", "index", "")
}

func NextHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.NEXT})

	RenderTemplate(w, "index", "index", "")
}

func PrevHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.BACK})

	RenderTemplate(w, "index", "index", "")
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.PAUSE})

	RenderTemplate(w, "index", "index", "")
}

func GetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.VOLUME_GET})
	fmt.Fprint(w, pipeReadLine())
}

func GetTitleHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.GET_TITLE})
	fmt.Fprint(w, pipeReadLine())
}

func GetArtistHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.GET_ARTIST})
	fmt.Fprint(w, pipeReadLine())
}

func SetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	pipeWrite([]byte(string(C.VOLUME_SET) + urlPart[2]))
}

func GetRemainingHandler(w http.ResponseWriter, r *http.Request) {
	pipeWrite([]byte{C.TIME_GET_REMAINING})
	fmt.Fprint(w, pipeReadLine())
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
