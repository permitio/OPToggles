package httpserver

import (
	"log"
	"net/http"
)

func ServeHttp(listenAddr string) {
	go func() {
		log.Fatal(http.ListenAndServe(listenAddr, nil))
	}()
}
