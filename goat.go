package main

import (
	"flag"
	"fmt"
	"github.com/goat-project/goat/service"
	"log"
)

// CLI option names
var (
	ip           = flag.String("listen-ip", "127.0.0.1", "IP address to bind to")
	port         = flag.Uint("port", 9623, "port to bind to")
	tls          = flag.Bool("tls", false, "True uses TLS, false uses plaintext TCP")
	certFile     = flag.String("cert-file", "server.pem", "server certificate file")
	keyFile      = flag.String("key-file", "server.key", "server key file")
	outDir       = flag.String("out-dir", "", "output directory")
	templatesDir = flag.String("templates-dir", "", "templates directory")
)

func checkArgs() error {
	if *tls {
		if *certFile == "" {
			return fmt.Errorf("Please specify a -cert-file")
		}
		if *keyFile == "" {
			return fmt.Errorf("Please specify a -key-file")
		}
	}

	if *outDir == "" {
		return fmt.Errorf("Please specify an -out-dir")
	}

	if *templatesDir == "" {
		return fmt.Errorf("Please specify a -templates-dir")
	}

	return nil
}

func main() {
	flag.Parse()

	err := checkArgs()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = service.Serve(ip, port, tls, certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
}
