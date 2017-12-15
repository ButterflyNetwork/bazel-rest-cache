package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/ButterflyNetwork/bazel-rest-cache/app"
)

func main() {
	port := flag.Int(
		"port",
		8080,
		"The port the HTTP server listens on",
	)
	redisAddr := flag.String(
		"redis_addr",
		"",
		"The address of the redis server (host:port).",
	)
	flag.Parse()

	if len(*redisAddr) == 0 {
		flag.Usage()
		return
	}

	c := app.NewRedisBazelCache(*redisAddr)
	r := mux.NewRouter()
	app.NewCacheApp(r, c)

	main := alice.New(
		compressHandler,
		logHandler,
	).Then(r)

	p := ":" + strconv.Itoa(*port)
	fmt.Println("Starting bazel-rest-cache on localhost" + p + " @ redis://" + *redisAddr)
	log.Fatal(http.ListenAndServe(p, main))
}

func compressHandler(h http.Handler) http.Handler {
	return handlers.CompressHandlerLevel(h, gzip.BestCompression)
}

func logHandler(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}