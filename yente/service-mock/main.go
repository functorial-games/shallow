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
	coreBinary := flag.String("core", "", "absolute path to existing Idriç yente-conversation executable")
	stateRoot := flag.String("state-root", "", "absolute path to trusted fictional Idriç journal root")
	materials := flag.String("materials", "", "absolute path to original Yente conversation materials")
	routeTokensPath := flag.String("fixture-tokens", "", "private TSV of fictional route names and distinct bearer tokens")
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
	if !loopback && *routeTokensPath == "" && len(token) < 24 {
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
	gateway := newGateway(routes, token)
	if *coreBinary != "" || *stateRoot != "" || *materials != "" {
		if *coreBinary == "" || *stateRoot == "" || *materials == "" {
			log.Fatal("-core, -state-root and -materials must be supplied together")
		}
		bridge, err := newCoreBridge(*coreBinary, *stateRoot, *materials)
		if err != nil {
			log.Fatal(err)
		}
		gateway.Core = bridge
	}
	if *routeTokensPath != "" {
		if gateway.Core == nil {
			log.Fatal("-fixture-tokens requires the configured Idriç core")
		}
		credentials, err := os.Open(*routeTokensPath)
		if err != nil {
			log.Fatal(err)
		}
		info, err := credentials.Stat()
		if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 {
			_ = credentials.Close()
			log.Fatal("-fixture-tokens must be a private regular file")
		}
		gateway.RouteTokens, err = loadRouteTokens(credentials, routes)
		_ = credentials.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	server := &http.Server{
		Addr: *listen, Handler: gateway,
		ReadHeaderTimeout: 5 * time.Second, MaxHeaderBytes: 8192,
	}
	if gateway.statefulFixtures() {
		fmt.Fprintln(os.Stderr, "Yente fictional HTTPS transport: CLI, FastChat protocol and simulated SMS connected to Idriç journal")
	} else if gateway.Core != nil {
		fmt.Fprintln(os.Stderr, "Yente fictional HTTPS transport: stateful CLI connected to Idriç journal; FastChat/SMS remain transport-only")
	} else {
		fmt.Fprintln(os.Stderr, "Yente fictional HTTPS transport: acknowledgements only (no Idriç journal)")
	}
	log.Fatal(server.ListenAndServeTLS(*certPath, *keyPath)) // TLS negotiates HTTP/2 via Go's standard server
}
