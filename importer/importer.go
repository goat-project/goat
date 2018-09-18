package importer

import (
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

// ProcessVms is a GRPC call -- do not use!
func (asi AccountingServiceImpl) ProcessVms(stream goat_grpc.AccountingService_ProcessVmsServer) error {
	for {
		data, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}

		if err != nil {
			return err
		}

		for _, vm := range data.Vms {
			asi.vmConsumer <- vm
		}
	}
}

// ProcessIps is a GRPC call -- do not use!
func (asi AccountingServiceImpl) ProcessIps(stream goat_grpc.AccountingService_ProcessIpsServer) error {
	for {
		data, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}

		if err != nil {
			return err
		}

		for _, ip := range data.Ips {
			asi.ipConsumer <- ip
		}
	}
}

// ProcessStorage is a GRPC call -- do not use!
func (asi AccountingServiceImpl) ProcessStorage(stream goat_grpc.AccountingService_ProcessStorageServer) error {
	for {
		data, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}

		if err != nil {
			return err
		}

		for _, storage := range data.Storages {
			asi.storageConsumer <- storage
		}
	}
}
