package trackers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"optoggles/config"
	"optoggles/utils"
	"strings"
)
const UpdateEndpoint string = "/opa_update"

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

type OpalClient struct {
	config config.OpalConfig
	triggerChan chan bool
}

func (oc *OpalClient) getOpaClient(ctx context.Context) (opaClient *OpaClient, err error) {
	log.Printf("Querying opa details from opal")

	var responseData []byte
	if responseData, err = utils.HttpGetBody(ctx, oc.config.Url, "/policy-store/config", oc.config.Token);
	err != nil {
		return
	}

	responseJson := struct {
		Url string `json:"url"`
		Token string `json:"token"`
	}{}
	if err = json.Unmarshal(responseData, &responseJson); err != nil {
		log.Printf("Can't parse OPA query result: " + err.Error())
		return
	}
	log.Printf("OPA details. url: %s, token: %s", responseJson.Url, responseJson.Token)

	return NewOpaClient(responseJson.Url, responseJson.Token), nil
}

func (oc *OpalClient) registerCallback(ctx context.Context) (err error) {
	token, err := callbackServer.AddClient(oc.config.Id, oc.triggerChan)
	if err != nil {
		return err
	}

	// TODO: Use struct for readability?
	callbackRegister := fmt.Sprintf(
		`{"key": "%s", "url": "%s", "config": {"method": "post", "headers": {"Authorization": "Bearer %s"}}}`,
		oc.config.AdvertisedAddress,
		"http://" + oc.config.AdvertisedAddress + UpdateEndpoint,
		token)

	return utils.HttpPost(ctx, oc.config.Url, "/callbacks", oc.config.Token, []byte(callbackRegister))
}

func (oc *OpalClient) waitForTrigger(ctx context.Context) error {
	select {
	case <-oc.triggerChan:
		// Empty the trigger channel, avoiding redundant updates
		for i := 0; i < len(oc.triggerChan); i++ {
			_ = <-oc.triggerChan
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func NewOpalClient(config config.OpalConfig) *OpalClient {
	return &OpalClient{
		config: config,
		triggerChan: make(chan bool, 1024),
	}
}

type CallbackServer struct {
	triggers map[string]chan bool
	key []byte
}

func NewCallbackServer() (*CallbackServer, error) {
	key := make([]byte, 20)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return &CallbackServer{
		triggers: make(map[string]chan bool),
		key: key,
	}, nil
}

func (cs *CallbackServer) AddClient(id string, triggerChan chan bool) (string, error) {
	cs.triggers[id] = triggerChan

	claims := make(jwt.MapClaims)
	claims["id"] = id
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = claims

	return token.SignedString(cs.key)
}

func (cs *CallbackServer) ParseToken(auth string) (string, error) {
	signedToken := strings.TrimPrefix(auth, "Bearer ")

	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) { return cs.key, nil } )
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return "", errors.New("invalid claims")
	} else {
		return claims["id"].(string), nil
	}
}

func (cs *CallbackServer) ServeForever() {
	http.HandleFunc(UpdateEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sourceTrackerId, err := cs.ParseToken(r.Header.Get("Authorization"))
		if err != nil {
			log.Printf("Got unauthorized callback from: %s", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if trigger, ok := callbackServer.triggers[sourceTrackerId]; ok {
			log.Printf("Got callback from opal source %s", sourceTrackerId)
			trigger <- true
		} else {
			log.Printf("Got update callback of unknown source: %s, from: %s", sourceTrackerId, r.RemoteAddr)
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// TODO: Use context, Port should be configurable(?)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var callbackServer *CallbackServer

func init() {
	var err error
	callbackServer, err = NewCallbackServer()
	if err != nil {
		log.Fatalf("Failed to create opal callback server: %s", err.Error())
	}

	// Execute singleton http callback server
	go callbackServer.ServeForever()
}