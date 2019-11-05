package importer

import (
	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/golang/protobuf/ptypes/empty"
)

var emptyMessage = empty.Empty{}

// WrapVms wraps gRPC stream of VMs in a RecordStream
func WrapVms(vms goat_grpc.AccountingService_ProcessVmsServer) RecordStream {
	return grpcStreamWrapper{
		recv: func() (identifiable, error) {
			return vms.Recv()
		},
		close: func() error {
			return vms.SendAndClose(&emptyMessage)
		},
	}
}

// WrapIps wraps gRPC stream of IPs in a RecordStream
func WrapIps(ips goat_grpc.AccountingService_ProcessIpsServer) RecordStream {
	return grpcStreamWrapper{
		recv: func() (identifiable, error) {
			return ips.Recv()
		},
		close: func() error {
			return ips.SendAndClose(&emptyMessage)
		},
	}
}

// WrapStorages wraps gRPC stream of Storage records in a RecordStream
func WrapStorages(storages goat_grpc.AccountingService_ProcessStoragesServer) RecordStream {
	return grpcStreamWrapper{
		recv: func() (identifiable, error) {
			return storages.Recv()
		},
		close: func() error {
			return storages.SendAndClose(&emptyMessage)
		},
	}
}

type identifiable interface {
	GetIdentifier() string
}

type grpcStreamWrapper struct {
	recv  func() (identifiable, error)
	close func() error
}

func (gsw grpcStreamWrapper) ReceiveIdentifier() (string, error) {
	received, err := gsw.recv()
	if err != nil {
		return "", err
	}

	identifier := received.(identifiable).GetIdentifier()
	if identifier == "" {
		return "", ErrFirstClientIdentifier
	}

	return identifier, nil
}

func (gsw grpcStreamWrapper) Receive() (wrapper.RecordWrapper, error) {
	data, err := gsw.recv()

	if err != nil {
		return nil, err
	}

	switch data.(type) {
	case *goat_grpc.IpData:
		ipd := data.(*goat_grpc.IpData).GetIp()
		if ipd == nil {
			return nil, ErrNonFirstClientIdentifier
		}

		return wrapper.WrapIP(*ipd), nil
	case *goat_grpc.VmData:
		vmd := data.(*goat_grpc.VmData).GetVm()
		if vmd == nil {
			return nil, ErrNonFirstClientIdentifier
		}

		return wrapper.WrapVM(*vmd), nil
	case *goat_grpc.StorageData:
		std := data.(*goat_grpc.StorageData).GetStorage()
		if std == nil {
			return nil, ErrNonFirstClientIdentifier
		}

		return wrapper.WrapStorage(*std), nil
	default:
		return nil, ErrUnknownMessageType
	}
}

func (gsw grpcStreamWrapper) Close() error {
	return gsw.close()
}
