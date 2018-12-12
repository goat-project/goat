package wrapper

import (
	"github.com/goat-project/goat-proto-go"
)

type storageWrapper struct {
	st goat_grpc.StorageRecord
}

type storageXML struct {
}

// NewStorageWrapper wraps given storage in a RecordWrapper
func NewStorageWrapper(st goat_grpc.StorageRecord) RecordWrapper {
	return storageWrapper{
		st: st,
	}
}

func (sw storageWrapper) Filename() string {
	return sw.st.GetRecordID()
}

func (sw storageWrapper) AsXML() (interface{}, error) {
	// TODO
	return storageXML{}, nil
}

func (sw storageWrapper) AsJSON() (interface{}, error) {
	return nil, ErrNotImplemented
}

func (sw storageWrapper) AsTemplate() (interface{}, error) {
	return nil, ErrNotImplemented

}
