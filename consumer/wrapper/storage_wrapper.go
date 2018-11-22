package wrapper

import (
	"github.com/goat-project/goat-proto-go"
)

type storageWrapper struct {
	st goat_grpc.StorageRecord
}

type storageXML struct {
}

func NewStorageWrapper(st goat_grpc.StorageRecord) RecordWrapper {
	return storageWrapper{
		st: st,
	}
}

func (sw storageWrapper) Filename() string {
	return sw.st.GetRecordID()
}

func (sw storageWrapper) AsXML() interface{} {
	// TODO
	return storageXML{}
}

func (sw storageWrapper) AsJSON() interface{} {
	return sw.st
}

func (sw storageWrapper) AsTemplate() interface{} {
	return sw.st

}
