package importer

import (
	"errors"
	"github.com/goat-project/goat-proto-go"
	"github.com/golang/protobuf/ptypes/empty"
	"io"
)

var (
	// ErrFirstClientIdentifier indicates that the first message of the stream is not client identifier
	ErrFirstClientIdentifier = errors.New("First message in the stream must be client identifier")
	// ErrNonFirstClientIdentifier indicates that client identifier was found as a non-first message of the stream
	ErrNonFirstClientIdentifier = errors.New("Client identifier found as a non-first message of the stream")
	// ErrUnknownMessageType indicates that an unknown type has arrived as part of data stream
	ErrUnknownMessageType = errors.New("Unhandled message type received")
)

// AccountingServiceImpl implements goat_grpc.AccountingService server
type AccountingServiceImpl struct {
	vmConsumer      chan<- goat_grpc.VmRecord
	ipConsumer      chan<- goat_grpc.IpRecord
	storageConsumer chan<- goat_grpc.StorageRecord
}

// NewAccountingServiceImpl creates a grpc server that sends received data to given channels and uses clientIdentifierValidator to validate client identifiers
func NewAccountingServiceImpl(vms chan<- goat_grpc.VmRecord, ips chan<- goat_grpc.IpRecord, storages chan<- goat_grpc.StorageRecord) AccountingServiceImpl {
	return AccountingServiceImpl{
		vmConsumer:      vms,
		ipConsumer:      ips,
		storageConsumer: storages,
	}
}

func (asi AccountingServiceImpl) receiveIdentifier(stream goat_grpc.AccountingService_ProcessServer) (string, error) {
	id, err := stream.Recv()
	if err != nil {
		return "", err
	}

	switch id.Data.(type) {
	case *goat_grpc.AccountingData_Identifier:
		return id.GetIdentifier(), nil
	default:
		return "", ErrFirstClientIdentifier
	}
}

// Process is a GRPC call -- do not use!
func (asi AccountingServiceImpl) Process(stream goat_grpc.AccountingService_ProcessServer) error {
	// TODO: use the first return value!
	_, err := asi.receiveIdentifier(stream)
	if err != nil {
		return err
	}

	for {
		data, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}

		if err != nil {
			return err
		}

		switch data.Data.(type) {
		case *goat_grpc.AccountingData_Identifier:
			return ErrNonFirstClientIdentifier
		case *goat_grpc.AccountingData_Vm:
			asi.vmConsumer <- *data.GetVm()
		case *goat_grpc.AccountingData_Ip:
			asi.ipConsumer <- *data.GetIp()
		case *goat_grpc.AccountingData_Storage:
			asi.storageConsumer <- *data.GetStorage()
		default:
			return ErrUnknownMessageType
		}
	}
}
