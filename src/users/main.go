package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/karthequian/ecommerce/src/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	helloCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_calls",
			Help: "Number of authentication calls.",
		},
		[]string{"url"},
	)
	authFailCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_failure",
			Help: "Number of authentication calls.",
		},
		[]string{"url"},
	)
	authSuccessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_success",
			Help: "Number of authentication calls.",
		},
		[]string{"url"},
	)

)

var tokenMap map[string]common.User

func init() {
	tokenMap = make(map[string]common.User)

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(helloCounter)
	prometheus.MustRegister(authFailCounter)
	prometheus.MustRegister(authSuccessCounter)

	for _, user := range common.Userlist {
		tokenMap[user.Token] = user
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	helloCounter.With(prometheus.Labels{"url": "/hello"}).Inc()
	fmt.Fprintf(w, "Hello from Users API")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	helloCounter.With(prometheus.Labels{"url": "/login"}).Inc()
	u := r.URL.Query().Get("u")
	p := r.URL.Query().Get("p")

	for _, user := range common.Userlist {
		if user.Username == u {
			if p == user.Password {
				returnUser := user
				returnUser.Password = ""
				jsonuser, _ := json.Marshal(returnUser)
				fmt.Fprintf(w, string(jsonuser))
				authSuccessCounter.With(prometheus.Labels{"url": "/login"}).Inc()
				return
			}
		}
	}
	authFailCounter.With(prometheus.Labels{"url": "/login"}).Inc()
	log.Debugf("Login information was invalid")
	http.Error(w, "Invalid username or password", 401)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Version handler was called")
	helloCounter.With(prometheus.Labels{"url": "/version"}).Inc()
	fmt.Fprintf(w, "{'version':'1.0'}")
}
func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Status handler was called")
	helloCounter.With(prometheus.Labels{"url": "/status"}).Inc()
	fmt.Fprintf(w, "{'status':'ok'}")
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	helloCounter.With(prometheus.Labels{"url": "/token"}).Inc()
	t := r.URL.Query().Get("t")

	foundUser := tokenMap[t]
	if len(foundUser.Token) > 0 {
		foundUser.Password = ""
		jsonuser, _ := json.Marshal(foundUser)
		fmt.Fprintf(w, string(jsonuser))
		return
	}
	log.Debugf("Invalid token was passed")
	http.Error(w, "Invalid token", 401)

}

func newHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("newhandler was called")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Welcome to the Users API. Valid endpoints are users/login, users/token, /version, /status, /metrics")

	helloCounter.With(prometheus.Labels{"url": "/list"}).Inc()
}

func main() {
	log.Info(os.Environ())
	port := os.Getenv("PORT")
	log.Infof("Port: %v", port)
	if len(port) == 0 {
		log.Fatalf("Port wasn't passed. An env variable for port must be passed")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", newHandler)
	r.HandleFunc("/users/login", loginHandler)
	r.HandleFunc("/users/token", TokenHandler)
	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/status", statusHandler)
	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Starting up server")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
