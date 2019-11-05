package service

import (
	"net"

	"github.com/sirupsen/logrus"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer"
	"github.com/goat-project/goat/importer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Serve starts grpc server on ip:port, optionally using tls. If *tls == true, then *certFile and
// *keyFile must be != null
func Serve(address string, tls bool, certFile, keyFile, outDir, templatesDir string, vmPerFile, ipPerFile,
	stPerFile uint64) error {
	server, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	if tls {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)

	vmWriter := consumer.NewTemplateGroupWriter(outDir, templatesDir, vmPerFile)
	ipWriter := consumer.NewJSONGroupWriter(outDir, ipPerFile)
	stWriter := consumer.NewXMLGroupWriter(outDir, stPerFile)
	goat_grpc.RegisterAccountingServiceServer(grpcServer, importer.NewAccountingServiceImpl(vmWriter, ipWriter, stWriter))

	logrus.WithField("address", address).Debug("gRPC server listening at")

	return grpcServer.Serve(server)
}
