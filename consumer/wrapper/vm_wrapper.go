package wrapper

import (
	"github.com/goat-project/goat-proto-go"
)

type vmWrapper struct {
	vm goat_grpc.VmRecord
}

type vmTemplate struct {
	// TODO
}

// NewVMWrapper wraps given vm in a RecordWrapper
func NewVMWrapper(vm goat_grpc.VmRecord) RecordWrapper {
	return vmWrapper{
		vm: vm,
	}
}

func (vw vmWrapper) Filename() string {
	return vw.vm.GetVmUuid()
}

func (vw vmWrapper) AsJSON() interface{} {
	return vw.vm
}

func (vw vmWrapper) AsXML() interface{} {
	return vw.vm
}

func (vw vmWrapper) AsTemplate() interface{} {
	// TODO
	return vmTemplate{}
}
