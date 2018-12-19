package wrapper

import (
	"fmt"

	goat_grpc "github.com/goat-project/goat-proto-go"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

type vmWrapper struct {
	vm goat_grpc.VmRecord
}

type vmTemplate struct {
	// naming of these fields must be consistent with vm.tmpl

	VMUUID              string
	SiteName            string
	CloudComputeService *string
	MachineName         string
	LocalUserID         *string
	LocalGroupID        *string
	GlobalUserName      *string
	Fqan                *string
	Status              *string
	StartTime           *string
	EndTime             *string
	SuspendDuration     *string
	WallDuration        *string
	CPUDuration         *string
	CPUCount            uint32
	NetworkType         *string
	NetworkInbound      *uint64
	NetworkOutbound     *uint64
	PublicIPCount       *uint64
	Memory              *uint64
	Disk                *uint64
	BenchmarkType       *string
	Benchmark           *float32
	StorageRecordID     *string
	ImageID             *string
	CloudType           *string
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

func (vw vmWrapper) AsJSON() (interface{}, error) {
	return nil, ErrNotImplemented
}

func (vw vmWrapper) AsXML() (interface{}, error) {
	return nil, ErrNotImplemented
}

func u64(wr *wrappers.UInt64Value) *uint64 {
	if wr == nil {
		return nil
	}
	result := new(uint64)
	*result = wr.GetValue()

	return result
}

func f32(wr *wrappers.FloatValue) *float32 {
	if wr == nil {
		return nil
	}
	result := new(float32)
	*result = wr.GetValue()

	return result
}

func s(wr fmt.Stringer) *string {
	if wr == nil {
		return nil
	}

	result := new(string)
	*result = wr.String()
	return result
}

func (vw vmWrapper) AsTemplate() (interface{}, error) {
	return vmTemplate{
		VMUUID:              vw.vm.GetVmUuid(),
		SiteName:            vw.vm.GetSiteName(),
		CloudComputeService: s(vw.vm.GetCloudComputeService()),
		MachineName:         vw.vm.GetMachineName(),
		LocalUserID:         s(vw.vm.GetLocalUserId()),
		LocalGroupID:        s(vw.vm.GetLocalGroupId()),
		GlobalUserName:      s(vw.vm.GetGlobalUserName()),
		Fqan:                s(vw.vm.GetFqan()),
		Status:              s(vw.vm.GetStatus()),
		StartTime:           s(vw.vm.GetStartTime()),
		EndTime:             s(vw.vm.GetEndTime()),
		SuspendDuration:     s(vw.vm.GetSuspendDuration()),
		CPUDuration:         s(vw.vm.GetCpuDuration()),
		WallDuration:        s(vw.vm.GetWallDuration()),
		CPUCount:            vw.vm.GetCpuCount(),
		NetworkType:         s(vw.vm.GetNetworkType()),
		NetworkInbound:      u64(vw.vm.GetNetworkInbound()),
		NetworkOutbound:     u64(vw.vm.GetNetworkOutbound()),
		PublicIPCount:       u64(vw.vm.GetPublicIpCount()),
		Memory:              u64(vw.vm.GetMemory()),
		Disk:                u64(vw.vm.GetDisk()),
		BenchmarkType:       s(vw.vm.GetBenchmarkType()),
		Benchmark:           f32(vw.vm.GetBenchmark()),
		StorageRecordID:     s(vw.vm.GetStorageRecordId()),
		ImageID:             s(vw.vm.GetImageId()),
		CloudType:           s(vw.vm.GetCloudType()),
	}, nil
}
