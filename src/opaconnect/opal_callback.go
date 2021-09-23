package opaconnect

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		request = append(request, "\n")
		request = append(request, string(body))
	}

	// Return the request as a string
	return strings.Join(request, "\n")
}

type Trigger interface {
	waitForTrigger(ctx context.Context) error
}

type OpalCBServer struct {
	token string
	triggerChan chan bool
}

func (ocbs *OpalCBServer) waitForTrigger(ctx context.Context) error {
	select {
	case <-ocbs.triggerChan:
		// Empty the trigger channel, avoiding redundant updates
		for i := 0; i < len(ocbs.triggerChan); i++ {
			_ = <-ocbs.triggerChan
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (ocbs *OpalCBServer) registerCallback() error {
	return nil
}

func (ocbs *OpalCBServer) serveOPALCallback() error {
	go func() {
		http.HandleFunc("/opa_update", func(w http.ResponseWriter, r *http.Request) {
			//log.Printf(formatRequest(r))
			//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
			log.Printf("Got update callback from OPAL")
			ocbs.triggerChan <- true
		})

		// TODO: Port should be configurable
		// TODO: HTTPS
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	return nil
}

func NewOpalCBServer() (*OpalCBServer, error) {
	ocbs := &OpalCBServer{
		token:       "xxx",
		triggerChan: make(chan bool, 1024),
	}
	if err := ocbs.registerCallback(); err != nil {
		return nil, err
	}
	if err := ocbs.serveOPALCallback(); err != nil {
		return nil, err
	}
	return ocbs, nil
}