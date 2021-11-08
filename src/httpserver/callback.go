package httpserver

import (
	"crypto/rand"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

const UpdateEndpoint string = "/opa_update"

type CallbackServer struct {
	triggers map[string]chan bool
	key      []byte
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

	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) { return cs.key, nil })
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

func (cs *CallbackServer) RegisterHandler() {
	http.HandleFunc(UpdateEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sourceTrackerId, err := cs.ParseToken(r.Header.Get("Authorization"))
		if err != nil {
			log.Printf("Got unauthorized callback from: %s", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if trigger, ok := cs.triggers[sourceTrackerId]; ok {
			log.Printf("Got callback from opal source %s", sourceTrackerId)
			trigger <- true
		} else {
			log.Printf("Got update callback of unknown source: %s, from: %s", sourceTrackerId, r.RemoteAddr)
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func NewCallbackServer() (*CallbackServer, error) {
	key := make([]byte, 20)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return &CallbackServer{
		triggers: make(map[string]chan bool),
		key:      key,
	}, nil
}

var CBServer *CallbackServer

func init() {
	var err error
	if CBServer, err = NewCallbackServer(); err != nil {
		log.Fatalf("Failed to create opal callback server: %s", err.Error())
	}
	CBServer.RegisterHandler()
}
