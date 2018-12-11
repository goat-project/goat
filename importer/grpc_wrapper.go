package importer

import (
	"github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/golang/protobuf/ptypes/empty"
)

var (
	emptyMessage = empty.Empty{}
)

// WrapVms wraps GRPC stream of VMs in a RecordStream
func WrapVms(vms goat_grpc.AccountingService_ProcessVmsServer) RecordStream {
	return grpcStreamWrapper{stream: vms.(grpcStream)}
}

// WrapIps wraps GRPC stream of IPs in a RecordStream
func WrapIps(ips goat_grpc.AccountingService_ProcessIpsServer) RecordStream {
	return grpcStreamWrapper{stream: ips.(grpcStream)}
}

// WrapStorages wraps GRPC stream of Storage records in a RecordStream
func WrapStorages(storages goat_grpc.AccountingService_ProcessStoragesServer) RecordStream {
	return grpcStreamWrapper{stream: storages.(grpcStream)}
}

type identifiable interface {
	GetIdentifier() string
}

type grpcStream interface {
	Recv() (identifiable, error)
	SendAndClose(*empty.Empty) error
}

type grpcStreamWrapper struct {
	stream grpcStream
}

func (gsw grpcStreamWrapper) ReceiveIdentifier() (string, error) {
	identifiable, err := gsw.stream.Recv()
	if err != nil {
		return "", err
	}

	identifier := identifiable.GetIdentifier()
	if identifier == "" {
		return "", ErrFirstClientIdentifier
	}

	return identifier, nil
}

func (gsw grpcStreamWrapper) Receive() (wrapper.RecordWrapper, error) {
	identifiable, err := gsw.stream.Recv()

	if err != nil {
		return nil, err
	}

	switch identifiable.(type) {
	case *goat_grpc.IpData:
		ipd := identifiable.(*goat_grpc.IpData).GetIp()
		if ipd == nil {
			return nil, ErrNonFirstClientIdentifier
		}

		return wrapper.WrapIP(*ipd), nil
	case *goat_grpc.VmData:
		vmd := identifiable.(*goat_grpc.VmData).GetVm()
		if vmd == nil {
			return nil, ErrNonFirstClientIdentifier
		}

		return wrapper.WrapVM(*vmd), nil
	case *goat_grpc.StorageData:
		std := identifiable.(*goat_grpc.StorageData).GetStorage()
		if std == nil {
			return nil, ErrNonFirstClientIdentifier
		}

		return wrapper.WrapStorage(*std), nil
	default:
		return nil, ErrUnknownMessageType
	}
}

func (gsw grpcStreamWrapper) Close() error {
	return gsw.stream.SendAndClose(&emptyMessage)
}
