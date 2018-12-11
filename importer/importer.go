package importer

import (
	"context"
	"errors"
	"github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer"
	"github.com/goat-project/goat/consumer/wrapper"
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
	vmConsumer      consumer.Consumer
	ipConsumer      consumer.Consumer
	storageConsumer consumer.Consumer
}

// NewAccountingServiceImpl creates a grpc server that sends received data to given channels and uses clientIdentifierValidator to validate client identifiers
func NewAccountingServiceImpl(vmConsumer consumer.Consumer, ipConsumer consumer.Consumer, storageConsumer consumer.Consumer) AccountingServiceImpl {
	return AccountingServiceImpl{
		vmConsumer:      vmConsumer,
		ipConsumer:      ipConsumer,
		storageConsumer: storageConsumer,
	}
}

func (asi AccountingServiceImpl) processStream(stream RecordStream, consumer consumer.Consumer) error {

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

	records := make(chan wrapper.RecordWrapper)
	results, err := consumer.Consume(consumerContext, id, records)
	if err != nil {
		return err
	}

	for {
		record, err := stream.Receive()
		if err == io.EOF {
			close(records)
			// caller should not be informed that an error occured if the stream just ended.
			err = nil
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
