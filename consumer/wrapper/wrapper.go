package wrapper

import (
	"github.com/goat-project/goat-proto-go"
)

// RecordWrapper is an interface to wrap record types
type RecordWrapper interface {
	// Filename returns name of file this should be saved in WITHOUT extension
	Filename() string

	// AsJson returns an annotated structure that can be serialized to JSON
	AsJSON() interface{}

	// AsXml returns an annotated structure that can be serialized to XML
	AsXML() interface{}

	// AsTemplate returns a structure that can be serialized via template
	AsTemplate() interface{}
}

// WrapVM wraps given vm in a RecordWrapper
func WrapVM(vm goat_grpc.VmRecord) RecordWrapper {
	return NewVMWrapper(vm)
}

// WrapIP wraps given ip in a RecordWrapper
func WrapIP(ip goat_grpc.IpRecord) RecordWrapper {
	return NewIPWrapper(ip)
}

// WrapStorage wraps given storage in a RecordWrapper
func WrapStorage(st goat_grpc.StorageRecord) RecordWrapper {
	return NewStorageWrapper(st)
}
