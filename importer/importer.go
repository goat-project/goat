package importer

import (
	"github.com/goat-project/goat-proto-go"
)

// AccountingServiceImpl implements goat_grpc.AccountingService server
type AccountingServiceImpl struct {
}

// NewAccountingServiceImpl creates a grpc server
func NewAccountingServiceImpl() AccountingServiceImpl {
	return AccountingServiceImpl{}
}

func (asi *AccountingServiceImpl) ProcessVms(stream goat_grpc.AccountingService_ProcessVmsServer) error {
	return nil
}

func (asi *AccountingServiceImpl) ProcessIps(stream goat_grpc.AccountingService_ProcessIpsServer) error {
	return nil
}

func (asi *AccountingServiceImpl) ProcessStorage(stream goat_grpc.AccountingService_ProcessStorageServer) error {
	return nil
}
