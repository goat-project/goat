package importer

import (
	"context"
	"io"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer"
	"github.com/goat-project/goat/consumer/wrapper"
)

// AccountingServiceImpl implements goat_grpc.AccountingService server
type AccountingServiceImpl struct {
	vmConsumer                                     consumer.Interface
	ipConsumer                                     consumer.Interface
	storageConsumer                                consumer.Interface
	gpuConsumer                                    consumer.Interface
	goat_grpc.UnimplementedAccountingServiceServer // warning: struct of size 72 bytes could be of size 64 bytes
}

// NewAccountingServiceImpl creates a grpc server that sends received data to given channels and
// uses clientIdentifierValidator to validate client identifiers
func NewAccountingServiceImpl(vmConsumer, ipConsumer, storageConsumer,
	gpuConsumer consumer.Interface) AccountingServiceImpl {
	return AccountingServiceImpl{
		vmConsumer:      vmConsumer,
		ipConsumer:      ipConsumer,
		storageConsumer: storageConsumer,
		gpuConsumer:     gpuConsumer,
	}
}

func (asi AccountingServiceImpl) processStream(stream RecordStream, consumer consumer.Interface) error {
	id, err := stream.ReceiveIdentifier()
	if err != nil {
		return err
	}

	consumerContext, cancelConsumers := context.WithCancel(context.Background())

	defer func() {
		cancelConsumers()
		e := stream.Close()
		if err == nil {
			err = e
		}
	}()

	endOfWriting := make(chan bool)
	records := make(chan wrapper.RecordWrapper)

	results, err := consumer.Consume(consumerContext, id, endOfWriting, records)
	if err != nil {
		return err
	}

	for {
		record, err := stream.Receive()
		if err == io.EOF {
			close(records)
			// caller should not be informed that an error occurred if the stream just ended.
			err = nil

			// It should wait until consumer finishes the work - until it writes all data to the file.
			<-endOfWriting

			return err
		}

		if err != nil {
			return err
		}

	inner:
		for {
			select {
			case records <- record:
				break inner
			// see if there is an error
			case res := <-results:
				if !res.IsOK() {
					return res.Error()
				}
			}
		}
	}
}

// ProcessVms is a GRPC call -- do not use
func (asi AccountingServiceImpl) ProcessVms(vms goat_grpc.AccountingService_ProcessVmsServer) error {
	return asi.processStream(WrapVms(vms), asi.vmConsumer)
}

// ProcessIps is a GRPC call -- do not use
func (asi AccountingServiceImpl) ProcessIps(ips goat_grpc.AccountingService_ProcessIpsServer) error {
	return asi.processStream(WrapIps(ips), asi.ipConsumer)
}

// ProcessStorages is a GRPC call -- do not use
func (asi AccountingServiceImpl) ProcessStorages(storages goat_grpc.AccountingService_ProcessStoragesServer) error {
	return asi.processStream(WrapStorages(storages), asi.storageConsumer)
}

// ProcessGPUs is a GRPC call -- do not use
func (asi AccountingServiceImpl) ProcessGPUs(gpus goat_grpc.AccountingService_ProcessGPUsServer) error {
	return asi.processStream(WrapGPUs(gpus), asi.gpuConsumer)
}
