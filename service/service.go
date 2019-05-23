package service

import (
	"fmt"
	"net"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer"
	"github.com/goat-project/goat/importer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Serve starts grpc server on ip:port, optionally using tls. If *tls == true, then *certFile and
// *keyFile must be != null
func Serve(ip *string, port *uint, tls *bool, certFile *string, keyFile *string, outDir *string,
	templatesDir *string, vmPerFile *uint64) error {
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)

	wr := consumer.NewTemplateGroupWriter(*outDir, *templatesDir, *vmPerFile)
	goat_grpc.RegisterAccountingServiceServer(grpcServer, importer.NewAccountingServiceImpl(wr, wr, wr))

	return grpcServer.Serve(server)
}
