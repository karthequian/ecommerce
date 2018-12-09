package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"time"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/karthequian/ecommerce/src/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/opentracing/opentracing-go/ext"
	opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"

	log "github.com/sirupsen/logrus"
)

var (
	helloCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_calls",
			Help: "Number of calls to view a cart.",
		},
		[]string{"url"},
	)
)

var tokenMap map[string]common.User
var closer io.Closer
var myTracer opentracing.Tracer

func init() {
	tokenMap = make(map[string]common.User)

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(helloCounter)
	for _, user := range common.Userlist {
		tokenMap[user.Token] = user
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Version handler was called")
	helloCounter.With(prometheus.Labels{"url": "/version"}).Inc()
	fmt.Fprintf(w, "{'version':'1.0'}")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Status handler was called")
	helloCounter.With(prometheus.Labels{"url": "/status"}).Inc()
	fmt.Fprintf(w, "{'status':'ok'}")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("List handler was called")

	//Create a span for tracing
	spanCtx, _ := myTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := myTracer.StartSpan("cart-list", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	vars := mux.Vars(r)
	fmt.Fprintf(w, "Cart: %v\n", common.CartList)
	// Increment prometheus count
	helloCounter.With(prometheus.Labels{"url": "/cart"}).Inc()
	
	time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)

	// Jaeger
	span.LogFields(
		otlog.String("event", "cart-list"),
		otlog.String("value", vars["user"]),
	)

}

func newHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("newhandler was called")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Welcome to the Cart API. Valid endpoints are /cart/{user}, /version, /status, /metrics")

	helloCounter.With(prometheus.Labels{"url": "/list"}).Inc()
}


func main() {

	myTracer, closer = common.InitTracer("cart")
	defer closer.Close()

	log.Info(os.Environ())
	port := os.Getenv("PORT")
	log.Infof("Port: %v", port)
	if len(port) == 0 {
		log.Fatalf("Port wasn't passed. An env variable for port must be passed")
	}

	r := mux.NewRouter()

	r.HandleFunc("/cart/{user}", listHandler)
	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/status", statusHandler)
	r.HandleFunc("/", newHandler)
	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Starting up server")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
