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

	// AsAPEL returns a structure that can be serialized via template
	AsTemplate() interface{}
}

func WrapVm(vm goat_grpc.VmRecord) RecordWrapper {
	return NewVmWrapper(vm)
}

func WrapIp(ip goat_grpc.IpRecord) RecordWrapper {
	return NewIpWrapper(ip)
}

func WrapStorage(st goat_grpc.StorageRecord) RecordWrapper {
	return NewStorageWrapper(st)
}
