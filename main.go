package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func transparentProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = target.Host
	}
	return proxy
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func maybeFail(failureRate float64, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() < failureRate {
			log.Print("ERROR: failing request with 502 Bad Gateway")
			w.WriteHeader(http.StatusBadGateway)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\tbadgateway [flags] targetURL\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	fail := flag.Float64("fail", 0.5, "rate at which requests fail")
	port := flag.Int64("port", 8080, "listen port")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	target := args[0]

	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	log.Printf("Starting proxy on %d with failure rate %f...\n", *port, *fail)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logRequest(maybeFail(*fail, transparentProxy(targetUrl)))))
}
