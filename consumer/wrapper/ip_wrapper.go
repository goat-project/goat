package wrapper

import (
	"github.com/goat-project/goat-proto-go"
)

type ipWrapper struct {
	ip goat_grpc.IpRecord
}

type ipJSON struct {
	// TODO
}

// NewIPWrapper wraps given ip in a RecordWrapper
func NewIPWrapper(ip goat_grpc.IpRecord) RecordWrapper {
	return ipWrapper{
		ip: ip,
	}

}

func (iw ipWrapper) Filename() string {
	// TODO
	return ""
}

func (iw ipWrapper) AsXML() interface{} {
	return iw.ip
}

func (iw ipWrapper) AsJSON() interface{} {
	// TODO
	return ipJSON{}
}

func (iw ipWrapper) AsTemplate() interface{} {
	return iw.ip
}
