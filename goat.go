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
	ipName           = "listen-ip"
	portName         = "port"
	tlsName          = "tls"
	certFileName     = "cert-file"
	keyFileName      = "key-file"
	outDirName       = "out-dir"
	templatesDirName = "templates-dir"
	vmPerFileName    = "vm-per-file"
	ipPerFileName    = "ip-per-file"
	stPerFileName    = "storage-per-file"

	logPathName = "log-path"
	debugName   = "debug"
)

var allFlags = []string{ipName, portName, tlsName, certFileName, keyFileName, outDirName, templatesDirName,
	vmPerFileName, ipPerFileName, stPerFileName, logPathName, debugName}

// CLI option values
var (
	ip           = flag.String(ipName, "127.0.0.1", "IP address to bind to")
	port         = flag.Uint(portName, 9623, "port to bind to")
	tls          = flag.Bool(tlsName, false, "True uses TLS, false uses plaintext TCP")
	certFile     = flag.String(certFileName, "server.pem", "server certificate file")
	keyFile      = flag.String(keyFileName, "server.key", "server key file")
	outDir       = flag.String(outDirName, "", "output directory")
	templatesDir = flag.String(templatesDirName, "", "templates directory")
	vmPerFile    = flag.Uint64(vmPerFileName, 500, "number of VMs per template file")
	ipPerFile    = flag.Uint64(ipPerFileName, 500, "number of IPs per json file")
	stPerFile    = flag.Uint64(stPerFileName, 500, "number of storages per xml file")

	logPath = flag.String(logPathName, "", "log path")
	debug   = flag.Bool(debugName, false, "True for debug mode, false otherwise")
)

func checkRequired(required []string) error {
	for _, req := range required {
		f := flag.Lookup(req)
		if f != nil {
			if f.Value.String() == "" {
				return fmt.Errorf("missing %s, please specify -%s", f.Usage, f.Name)
			}
		}
	}

	return nil
}

func logFlags(flags []string) {
	for _, f := range flags {
		fl := flag.Lookup(f)
		if fl != nil {
			logrus.WithFields(logrus.Fields{"name": fl.Name, "value": fl.Value}).Debug("flag initialized")
		}
	}
}

func main() {
	flag.Parse()

	logger.Init(*logPath, *debug)

	required := []string{outDirName, templatesDirName}

	if *tls {
		required = append(required, []string{certFileName, keyFileName}...)
	}

	err := checkRequired(required)
	if err != nil {
		logrus.WithField("error", err).Fatal("missing required argument")
		return
	}

	logFlags(allFlags)

	err = service.Serve(fmt.Sprintf("%s:%d", *ip, *port), *tls, *certFile, *keyFile, *outDir, *templatesDir,
		*vmPerFile, *ipPerFile, *stPerFile)
	if err != nil {
		logrus.WithField("error", err).Fatal("fatal error run goat")
	}
}
