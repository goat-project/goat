package importer

import (
	"errors"
	"github.com/goat-project/goat-proto-go"
	"github.com/golang/protobuf/ptypes/empty"
	"io"
)

// AccountingServiceImpl implements goat_grpc.AccountingService server
type AccountingServiceImpl struct {
	vmConsumer      chan<- *goat_grpc.VmRecord
	ipConsumer      chan<- *goat_grpc.IpRecord
	storageConsumer chan<- *goat_grpc.StorageRecord
}

// NewAccountingServiceImpl creates a grpc server that sends received data to given channels and uses clientIdentifierValidator to validate client identifiers
func NewAccountingServiceImpl(vms chan<- *goat_grpc.VmRecord, ips chan<- *goat_grpc.IpRecord, storages chan<- *goat_grpc.StorageRecord) AccountingServiceImpl {
	return AccountingServiceImpl{
		vmConsumer:      vms,
		ipConsumer:      ips,
		storageConsumer: storages,
	}
}

// Process is a GRPC call -- do not use!
func (asi AccountingServiceImpl) Process(stream goat_grpc.AccountingService_ProcessServer) error {
	id, err := stream.Recv()
	if err != nil {
		return err
	}

	clientID := id.GetIdentifier()
	if clientID == "" {
		return errors.New("First message in the stream must be client identifier")
	}

	for {
		data, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}

		if err != nil {
			return err
		}

		vm := data.GetVm()
		ip := data.GetIp()
		storage := data.GetStorage()

		if vm != nil {
			asi.vmConsumer <- vm
		}

		if ip != nil {
			asi.ipConsumer <- ip
		}

		if storage != nil {
			asi.storageConsumer <- storage
		}

		if data.GetIdentifier() != "" {
			return errors.New("Client identifier found as a non-first message of the stream")
		}
	}
}
