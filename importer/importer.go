package importer

import (
	"github.com/goat-project/goat-proto-go"
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

func (asi *AccountingServiceImpl) ProcessVms(stream goat_grpc.AccountingService_ProcessVmsServer) error {
	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&goat_grpc.Confirmation{
				Accepted: true,
				Msg:      "ok",
			})
		}

		if err != nil {
			return err
		}

		// TODO process data
	}
}

func (asi *AccountingServiceImpl) ProcessIps(stream goat_grpc.AccountingService_ProcessIpsServer) error {
	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&goat_grpc.Confirmation{
				Accepted: true,
				Msg:      "ok",
			})
		}

		if err != nil {
			return err
		}

		// TODO process data
	}
}

func (asi *AccountingServiceImpl) ProcessStorage(stream goat_grpc.AccountingService_ProcessStorageServer) error {
	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&goat_grpc.Confirmation{
				Accepted: true,
				Msg:      "ok",
			})
		}

		if err != nil {
			return err
		}

		// TODO process data
	}
}
