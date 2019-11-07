package importer

import (
	"errors"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/golang/protobuf/ptypes/empty"
)

var emptyMessage = empty.Empty{}

// errors
var (
	// errFirstClientIdentifier indicates that the first message of the stream is not client identifier
	errFirstClientIdentifier = errors.New("first message in the stream must be client identifier")
	// errNonFirstClientIdentifier indicates that client identifier was found as a non-first message of the stream
	errNonFirstClientIdentifier = errors.New("client identifier found as a non-first message of the stream")
	// errUnknownMessageType indicates that an unknown type has arrived as part of data stream
	errUnknownMessageType = errors.New("unhandled message type received")
)

// WrapVms wraps GRPC stream of VMs in a RecordStream
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
		return "", errFirstClientIdentifier
	}

	return identifier, nil
}

func (gsw grpcStreamWrapper) Receive() (wrapper.RecordWrapper, error) {
	data, err := gsw.recv()

	if err != nil {
		return nil, err
	}

	switch data := data.(type) {
	case *goat_grpc.IpData:
		ipd := data.GetIp()
		if ipd == nil {
			return nil, errNonFirstClientIdentifier
		}

		return wrapper.WrapIP(*ipd), nil
	case *goat_grpc.VmData:
		vmd := data.GetVm()
		if vmd == nil {
			return nil, errNonFirstClientIdentifier
		}

		return wrapper.WrapVM(*vmd), nil
	case *goat_grpc.StorageData:
		std := data.GetStorage()
		if std == nil {
			return nil, errNonFirstClientIdentifier
		}

		return wrapper.WrapStorage(*std), nil
	default:
		return nil, errUnknownMessageType
	}
}

func (gsw grpcStreamWrapper) Close() error {
	return gsw.close()
}
