package main

import (
	"context"
	_ "expvar"
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/henvic/trigram/server"
)

var p server.Params
var debug bool

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	if debug {
		go profiler()
	}

	ctx := context.TODO()
	// TODO(henvic): use ctxsignal.WithTermination instead (I wrote it)
	// Reference: https://github.com/henvic/ctxsignal
	// ctx, cancel := ctxsignal.WithTermination(context.Background())
	// defer cancel()

	if err := server.Run(ctx, p); err != nil {
		log.Fatal(err)
	}
}

func profiler() {
	// let expvar and pprof be exposed here indirectly through http.DefaultServeMux
	log.Println("Exposing expvar and pprof on localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}

func init() {
	flag.StringVar(&p.Address, "addr", "127.0.0.1:8080", "Serving address")
	flag.BoolVar(&debug, "expose-debug", false, "Expose debugging tools over HTTP (on port 8081)")
}
