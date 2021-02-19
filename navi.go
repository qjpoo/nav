package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}
type home struct {
	Title string
}

const (
	Js_Dir  = "./js/"
	Css_Dir = "./css/"
	Img_Dir = "./img/"
)

func main() {
	//filesvr服务
	server := http.Server{
		Addr:        ":8002",
		Handler:     &Myhandler{},
		ReadTimeout: 100 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index
	mux["/js"] = jsFile
	mux["/css"] = cssFile
	mux["/img"] = imgFile
	fmt.Println("Hello, this is fancygo navi!")
	go server.ListenAndServe()

	//设置sigint信号
	close := make(chan os.Signal, 1)
	signal.Notify(close, os.Interrupt, os.Kill)
	<-close

	fmt.Println("Bye, fancygo webgame close")
}

func (*Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	if ok, _ := regexp.MatchString("/css/", r.URL.String()); ok {
		http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))).ServeHTTP(w, r)
	} else if ok, _ := regexp.MatchString("/js/", r.URL.String()); ok {
		http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))).ServeHTTP(w, r)
	} else if ok, _ := regexp.MatchString("/img/", r.URL.String()); ok {
		http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))).ServeHTTP(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	title := home{Title: "fancygo navi"}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, title)
}

func jsFile(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))).ServeHTTP(w, r)
}

func cssFile(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))).ServeHTTP(w, r)
}

func imgFile(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))).ServeHTTP(w, r)
}
