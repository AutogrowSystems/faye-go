package main

import (
	"flag"
	"github.com/autogrowsystems/faye-go"
	"github.com/autogrowsystems/faye-go/adapters"
	l "github.com/cenkalti/log"
	"net/http"
)

func OurLoggingHandler(pattern string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.NewLogger("http").Infof("%v: %+v", pattern, *r.URL)
		h.ServeHTTP(w, r)
	})
}

type config struct {
	Host   string
	Port   string
	Public string
}

func main() {
	// TODO: read config file
	var cfg config

	flag.StringVar(&cfg.Port, "p", "8000", "Port number to serve on")
	flag.StringVar(&cfg.Host, "h", "127.0.0.1", "Port number to serve on")
	flag.StringVar(&cfg.Public, "public", "src/github.com/autogrowsystems/faye-go/public", "Port number to serve on")
	flag.Parse()

	engineLog := l.NewLogger("engine")
	serverLog := l.NewLogger("server")
	httpLog := l.NewLogger("http")

	engineLog.SetLevel(l.DEBUG)
	serverLog.SetLevel(l.DEBUG)
	httpLog.SetLevel(l.DEBUG)

	fayeServer := faye.NewServer(serverLog, faye.NewEngine(engineLog))
	http.Handle("/faye", adapters.FayeHandler(fayeServer))

	httpLog.Infoln("Mounted faye server")

	// Also serve up some static files and show off
	// the wonderful go http handler chain
	http.Handle("/", OurLoggingHandler("/",
		http.FileServer(http.Dir(cfg.Public)),
	)) // TODO: put this in a config file

	httpLog.Infoln("Mounted file server")

	httpLog.Infoln("Listening on", cfg.Host+":"+cfg.Port)
	err := http.ListenAndServe(cfg.Host+":"+cfg.Port, nil)
	if err != nil {
		httpLog.Fatalln("Failed to start the server: " + err.Error())
	}
}
