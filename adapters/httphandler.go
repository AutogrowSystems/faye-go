package adapters

import (
	// "code.google.com/p/go.net/websocket"
	"encoding/json"
	"github.com/autogrowsystems/faye-go"
	"github.com/autogrowsystems/faye-go/transport"
	"github.com/gorilla/websocket"
	"net/http"
)

/* HTTP handler that can be dropped into the standard http handlers list */
func FayeHandler(server faye.Server) http.Handler {
	// websocketHandler := websocket.Handler(transport.WebsocketServer(server))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {

			ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
			if _, ok := err.(websocket.HandshakeError); ok {
				http.Error(w, "Not a websocket handshake", 400)
				return
			} else if err != nil {
				server.Logger().Errorf("ERROR: %s", err)
				return
			}

			// Start the websocket server
			transport.WebsocketServer(server)(ws)
			server.Logger().Warnf("Websocket server stopped")
		} else {
			if r.Method == "POST" {
				var v interface{}
				dec := json.NewDecoder(r.Body)
				if err := dec.Decode(&v); err == nil {

					// start the long poll server
					transport.MakeLongPoll(v, server, w)
					server.Logger().Warnf("Long poll server stopped")
				} else {
					server.Logger().Fatalf("ERROR: %s", err)
				}
			}
		}
	})
}

// func handler(w http.ResponseWriter, r *http.Request) {
//     ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
//     if _, ok := err.(websocket.HandshakeError); ok {
//         http.Error(w, "Not a websocket handshake", 400)
//         return
//     } else if err != nil {
//         log.Println(err)
//         return
//     }
//     ... Use conn to send and receive messages.
// }
