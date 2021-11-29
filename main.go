package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func main() {

	hosts := mustHosts()
	proxy := &httputil.ReverseProxy{Director: director(hosts)}
	listener := mustListener(hosts)
	go func() {
		log.Fatal(http.Serve(listener, proxy))
	}()

	// Redirect HTTP requests to secure protocol.
	redirect := http.Server{
		Handler:      http.HandlerFunc(redirectHandler),
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
		IdleTimeout:  time.Second * 60,
		Addr:         ":http",
	}
	println("Proxy listening on port 80 and port 443")
	log.Fatal(redirect.ListenAndServe())
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusMovedPermanently
	url := "https://" + r.Host + r.RequestURI
	http.Redirect(w, r, url, code)
}

type host struct {
	Host   string
	Port   string
	Scheme string
	Cert   string
	Key    string
}

func mustHosts() (hosts []host) {
	bb, err := os.ReadFile("./hosts.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bb, &hosts); err != nil {
		log.Fatal(err)
	}
	return hosts
}

func mustListener(hosts []host) net.Listener {
	var certs []tls.Certificate
	for _, h := range hosts {
		kp, err := tls.LoadX509KeyPair(h.Cert, h.Key)
		if err != nil {
			log.Fatal(err)
		}
		certs = append(certs, kp)
	}
	tlsConfig := &tls.Config{Certificates: certs}
	listener, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	return listener
}

func director(hosts []host) func(*http.Request) {
	return func(r *http.Request) {
		for _, h := range hosts {
			if r.Host == h.Host {
				r.URL.Host = "localhost:" + h.Port
				r.URL.Scheme = h.Scheme
				return
			}
		}
		fmt.Printf("unknown hostname: %q", r.Host)
	}
}
