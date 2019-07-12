package wrapper

import (
	goat_grpc "github.com/goat-project/goat-proto-go"
)

type ipWrapper struct {
	ip goat_grpc.IpRecord
}

type ipJSON struct {
	MeasurementTime     int64
	SiteName            string
	CloudComputeService *string
	CloudType           string
	LocalUser           string
	LocalGroup          string
	GlobalUserName      string
	FQAN                string
	IPVersion           byte
	IPCount             int
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

func (iw ipWrapper) AsXML() (interface{}, error) {
	return nil, ErrNotImplemented
}

func (iw ipWrapper) AsJSON() (interface{}, error) {
	return ipJSON{
		MeasurementTime:     iw.ip.GetMeasurementTime().Seconds,
		SiteName:            iw.ip.GetSiteName(),
		CloudComputeService: s(iw.ip.GetCloudComputeService()),
		CloudType:           iw.ip.GetCloudType(),
		LocalUser:           iw.ip.GetLocalUser(),
		LocalGroup:          iw.ip.GetLocalGroup(),
		GlobalUserName:      iw.ip.GetGlobalUserName(),
		FQAN:                iw.ip.GetFqan(),
		IPVersion:           b(iw.ip.GetIpType()),
		IPCount:             int(iw.ip.GetIpCount()),
	}, nil
}

func (iw ipWrapper) AsTemplate() (interface{}, error) {
	return nil, ErrNotImplemented
}
