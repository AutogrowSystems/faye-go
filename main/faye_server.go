package main

import (
	"github.com/AutogrowSystems/faye-go"
	"github.com/AutogrowSystems/faye-go/adapters"
	l "github.com/cenkalti/log"
	"net/http"
)

func OurLoggingHandler(pattern string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.NewLogger("http").Infof("%v: %+v", pattern, *r.URL)
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

	engineLog := l.NewLogger("engine")
	serverLog := l.NewLogger("server")
	httpLog := l.NewLogger("http")

	fayeServer := faye.NewServer(serverLog, faye.NewEngine(engineLog))
	http.Handle("/faye", adapters.FayeHandler(fayeServer))

	// Also serve up some static files and show off
	// the wonderful go http handler chain
	http.Handle("/", OurLoggingHandler("/",
		http.FileServer(http.Dir(cfg.Public)),
	)) // TODO: put this in a config file

	err := http.ListenAndServe(cfg.Host+":"+cfg.Port, nil)
	if err != nil {
		httpLog.Fatalln("Failed to start the server: " + err.Error())
	}

	httpLog.Infoln("Faye server started on", cfg.Host+":"+cfg.Port)
}
