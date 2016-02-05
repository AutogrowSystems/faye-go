package main

import (
	"github.com/AutogrowSystems/faye-go"
	"github.com/AutogrowSystems/faye-go/adapters"
	"log"
	"net/http"
)

func OurLoggingHandler(pattern string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %+v", pattern, *r.URL)
		h.ServeHTTP(w, r)
	})
}

var cfg = struct {
	Host   string
	Port   string
	Public string
}{"127.0.0.1", "8000", "src/github.com/AutogrowSystems/faye-go/public"}

func main() {
	// TODO: read config file

	fayeServer := faye.NewServer(faye.NewEngine())
	http.Handle("/faye", adapters.FayeHandler(fayeServer))

	// Also serve up some static files and show off
	// the wonderful go http handler chain
	http.Handle("/", OurLoggingHandler("/",
		http.FileServer(http.Dir(cfg.Public)),
	)) // TODO: put this in a config file

	err := http.ListenAndServe(cfg.Host+":"+cfg.Port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	fmt.Println("Faye server started on", cfg.Host+":"+cfg.Port)
}
