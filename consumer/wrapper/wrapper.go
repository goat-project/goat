package wrapper

import (
	"errors"

	goat_grpc "github.com/goat-project/goat-proto-go"
)

var (
	// ErrNotImplemented signals that the called method is not implemented
	ErrNotImplemented = errors.New("Not implemented")
)

// RecordWrapper is an interface to wrap record types
type RecordWrapper interface {
	// Filename returns name of file this should be saved in WITHOUT extension
	Filename() string

	// AsJSON returns an annotated structure that can be serialized to JSON.
	// ErrNotImplemented is returned if the operation is not implemented
	AsJSON() (interface{}, error)

	// AsXML returns an annotated structure that can be serialized to XML.
	// ErrNotImplemented is returned if the operation is not implemented
	AsXML() (interface{}, error)

	// AsTemplate returns a structure that can be serialized via template.
	// ErrNotImplemented is returned if the operation is not implemented
	AsTemplate() (interface{}, error)
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
