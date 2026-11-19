package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8443", "TLS address (loopback by default)")
	routesPath := flag.String("routes", "routes.synthetic.tsv", "fictional route table")
	certPath := flag.String("cert", "", "TLS certificate (with SAN for selected host)")
	keyPath := flag.String("key", "", "TLS private key")
	flag.Parse()
	if *certPath == "" || *keyPath == "" {
		log.Fatal("TLS -cert and -key are required; plaintext mode is not provided")
	}
	host, _, err := net.SplitHostPort(*listen)
	if err != nil {
		log.Fatal(err)
	}
	ip := net.ParseIP(host)
	loopback := host == "localhost" || (ip != nil && ip.IsLoopback())
	token := os.Getenv("YENTE_MOCK_TOKEN")
	if !loopback && len(token) < 24 {
		log.Fatal("non-loopback bind requires YENTE_MOCK_TOKEN (at least 24 characters); this is still NOT production auth")
	}
	file, err := os.Open(*routesPath)
	if err != nil {
		log.Fatal(err)
	}
	routes, err := loadRoutes(file)
	_ = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr: *listen, Handler: newGateway(routes, token),
		ReadHeaderTimeout: 5 * time.Second, MaxHeaderBytes: 8192,
	}
	fmt.Fprintln(os.Stderr, "Yente mock HTTPS transport listening on", *listen, "(no model, no Idriç journal, no real SMS)")
	log.Fatal(server.ListenAndServeTLS(*certPath, *keyPath)) // TLS negotiates HTTP/2 via Go's standard server
}
