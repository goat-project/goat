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
	stPerFile, gpuPerFile uint64) error {
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

	vmWriter := consumer.CreateConsumer(consumer.NewTemplateGroupWriter(outDir, templatesDir, vmPerFile))
	ipWriter := consumer.CreateConsumer(consumer.NewJSONGroupWriter(outDir, ipPerFile, consumer.IP))
	stWriter := consumer.CreateConsumer(consumer.NewXMLGroupWriter(outDir, stPerFile))
	gpuWriter := consumer.CreateConsumer(consumer.NewJSONGroupWriter(outDir, gpuPerFile, consumer.GPU))
	goat_grpc.RegisterAccountingServiceServer(grpcServer, importer.NewAccountingServiceImpl(vmWriter, ipWriter,
		stWriter, gpuWriter))

	logrus.WithField("address", address).Info("gRPC server listening at")

	return grpcServer.Serve(server)
}
