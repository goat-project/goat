package wrapper

import (
	"strconv"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// The following functions provide a conversion of wrapped values.

func u64(wr *wrappers.UInt64Value) *uint64 {
	if wr == nil {
		return nil
	}
	result := new(uint64)
	*result = wr.GetValue()

	return result
}

func u32(wr *wrappers.UInt32Value) *uint32 {
	if wr == nil {
		return nil
	}
	result := new(uint32)
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

func s(wr *wrappers.StringValue) *string {
	if wr == nil {
		return nil
	}

	result := new(string)
	*result = wr.GetValue()

	return result
}

func st(wr *timestamp.Timestamp) *string {
	if wr == nil {
		return nil
	}

	result := new(string)
	*result = strconv.FormatInt(wr.Seconds, 10)

	return result
}

func sd(wr *duration.Duration) *string {
	if wr == nil {
		return nil
	}

	result := new(string)
	*result = strconv.FormatInt(wr.Seconds, 10)

	return result
}

func b(wr string) byte {
	if wr == "IPv4" {
		return 4
	}

	if wr == "IPv6" {
		return 6
	}

	return 0
}
