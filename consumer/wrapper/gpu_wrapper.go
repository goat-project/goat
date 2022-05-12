package wrapper

import (
	goat_grpc "github.com/goat-project/goat-proto-go"
)

type gpuWrapper struct {
	gpu goat_grpc.GPURecord
}

type gpuJSON struct {
	MeasurementMonth     uint64
	MeasurementYear      uint64
	AssociatedRecordType string
	AssociatedRecord     string
	GlobalUserName       *string
	FQAN                 string
	SiteName             string
	Count                float32
	Cores                *uint32
	ActiveDuration       *uint64
	AvailableDuration    uint64
	BenchmarkType        *string
	Benchmark            *float32
	Type                 string
	Model                *string
}

// NewGPUWrapper wraps given gpu in a RecordWrapper
func NewGPUWrapper(gpu goat_grpc.GPURecord) RecordWrapper {
	return gpuWrapper{
		gpu: gpu,
	}

}

func (gw gpuWrapper) Filename() string {
	// TODO
	return ""
}

func (gw gpuWrapper) AsXML() (interface{}, error) {
	return nil, ErrNotImplemented
}

func (gw gpuWrapper) AsJSON() (interface{}, error) {
	return gpuJSON{ // TODO turn attributes
		MeasurementMonth:     gw.gpu.GetMeasurementMonth(),
		MeasurementYear:      gw.gpu.GetMeasurementYear(),
		AssociatedRecordType: gw.gpu.GetAssociatedRecordType(),
		AssociatedRecord:     gw.gpu.GetAssociatedRecord(),
		GlobalUserName:       s(gw.gpu.GetGlobalUserName()),
		FQAN:                 gw.gpu.GetFqan(),
		SiteName:             gw.gpu.GetSiteName(),
		Count:                gw.gpu.GetCount(),
		Cores:                u32(gw.gpu.GetCores()),
		ActiveDuration:       u64(gw.gpu.GetActiveDuration()),
		AvailableDuration:    gw.gpu.GetAvailableDuration(),
		BenchmarkType:        s(gw.gpu.GetBenchmarkType()),
		Benchmark:            f32(gw.gpu.GetBenchmark()),
		Type:                 gw.gpu.GetType(),
		Model:                s(gw.gpu.GetModel()),
	}, nil
}

func (gw gpuWrapper) AsTemplate() (interface{}, error) {
	return nil, ErrNotImplemented
}
