package main

import (
	"flag"
	"fmt"

	"github.com/goat-project/goat/logger"

	"github.com/goat-project/goat/service"
	"github.com/sirupsen/logrus"
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
	vmPerFile    = flag.Uint64("vm-per-file", 500, "number of VMs per template file")
	ipPerFile    = flag.Uint64("ip-per-file", 500, "number of IPs per json file")
	stPerFile    = flag.Uint64("storage-per-file", 500, "number of storages per xml file")

	logPath = flag.String("log-path", "", "log path")
	debug   = flag.Bool("debug", false, "True for debug mode, false otherwise")
)

func checkArgs() error {
	if *tls {
		if *certFile == "" {
			return fmt.Errorf("please specify a -cert-file")
		}
		if *keyFile == "" {
			return fmt.Errorf("please specify a -key-file")
		}
	}

	if *outDir == "" {
		return fmt.Errorf("please specify an -out-dir")
	}

	if *templatesDir == "" {
		return fmt.Errorf("please specify a -templates-dir")
	}

	return nil
}

func main() {
	flag.Parse()

	logger.Init(*logPath, *debug)

	err := checkArgs()
	if err != nil {
		logrus.WithField("error", err).Fatal("missing required argument")
		return
	}

	err = service.Serve(ip, port, tls, certFile, keyFile, outDir, templatesDir, vmPerFile, ipPerFile, stPerFile)
	if err != nil {
		logrus.WithField("error", err).Fatal("fatal error serve")
	}
}
