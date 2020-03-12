package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

var listenAddr = flag.String("listen", ":1969", "The address to listen on for http requests")
var listenAddrTLS = flag.String("listentls", ":1968", "The address to listen on for https requests")
var serverCert = flag.String("cert", "server.crt", "The file name containing a TLS server certificate")
var serverKey = flag.String("key", "server.key", "The file name containing a TLS server key")
var checkAddr = flag.String("check", ":1936", "The address to listen on for healty probes.")

func main() {

	// parse flags
	flag.Parse()

	// graceful server shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// mux
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(reply))     // mux.HandleFunc("/", reply)
	mux.Handle("/path", http.HandlerFunc(reply)) // mux.HandleFunc("/path", reply)

	// HTTP server
	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}

	// HTTPS server
	tlsSrv := &http.Server{
		Addr:    *listenAddrTLS,
		Handler: mux,
	}

	go serve(srv)
	go serveTLS(tlsSrv, *serverCert, *serverKey)
	go serveChecks(*checkAddr, healthz)

	<-quit
	log.Println("Shutting down gracefully the server ...")
	// Gracefully server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

}

func reply(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Fprintf(w, "Server Name = %q\n", hostname)
	fmt.Fprintf(w, "Client Addr = %q\n", r.RemoteAddr)
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "Method = %q\n", r.Method)
	fmt.Fprintf(w, "URL = %q\n", r.URL)
	fmt.Fprintf(w, "Protocol = %q\n", r.Proto)

	keys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintf(w, "+++ Request Headers: +++\n")
	for _, v := range keys {
		fmt.Fprintf(w, "Header[%q] = %q\n", v, r.Header[v])
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func serve(s *http.Server) {
	log.Printf("Server HTTP started at %s ...", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Printf("Starting server failed: %s", err)
	}
}

func serveTLS(s *http.Server, serverCert, serverKey string) {
	log.Printf("Server HTTPS started at %s ...", s.Addr)
	if err := s.ListenAndServeTLS(serverCert, serverKey); err != nil {
		log.Printf("Starting server failed: %s", err)
	}
}

func serveChecks(addr string, probe func(w http.ResponseWriter, r *http.Request)) {
	log.Printf("Serving checks at %s ...", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(probe)); err != nil {
		log.Printf("Starting checks listener failed")
	}
}
