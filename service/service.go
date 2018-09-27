package service

import (
	"fmt"
	"github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/importer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

const (
	vmBufferSize      = 32
	ipBufferSize      = 32
	storageBufferSize = 32
)

// Serve starts grpc server on ip:port, optionally using tls. If *tls == true, then *certFile and *keyFile must be != null
func Serve(ip *string, port *uint, tls *bool, certFile *string, keyFile *string) error {
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
	vms := make(chan goat_grpc.VmRecord, vmBufferSize)
	ips := make(chan goat_grpc.IpRecord, ipBufferSize)
	storages := make(chan goat_grpc.StorageRecord, storageBufferSize)

	goat_grpc.RegisterAccountingServiceServer(grpcServer, importer.NewAccountingServiceImpl(vms, ips, storages))

	return grpcServer.Serve(server)
}
