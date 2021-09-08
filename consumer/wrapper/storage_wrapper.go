package wrapper

import (
	"time"

	goat_grpc "github.com/goat-project/goat-proto-go"
)

type storageWrapper struct {
	st goat_grpc.StorageRecord
}

type storageXML struct {
	RecordID                  recordIdentity `xml:"RecordIdentity"`
	StorageSystem             string         `xml:"StorageSystem"`
	Site                      *string        `xml:"Site"`
	StorageShare              *string        `xml:"StorageShare"`
	StorageMedia              *string        `xml:"StorageMedia"`
	StorageClass              *string        `xml:"StorageClass"`
	FileCount                 *string        `xml:"FileCount"`
	DirectoryPath             *string        `xml:"DirectoryPath"`
	LocalUser                 *string        `xml:"LocalUser"`
	LocalGroup                *string        `xml:"LocalGroup"`
	UserIdentity              *string        `xml:"UserIdentity"`
	Group                     *string        `xml:"Group"`
	GroupAttribute            *string        `xml:"GroupAttribute"`
	GroupAttributeType        *string        `xml:"GroupAttributeType"`
	StartTime                 time.Time      `xml:"StartTime"`
	EndTime                   time.Time      `xml:"EndTime"`
	ResourceCapacityUsed      uint64         `xml:"ResourceCapacityUsed"`
	LogicalCapacityUsed       *uint64        `xml:"LogicalCapacityUsed"`
	ResourceCapacityAllocated *uint64        `xml:"ResourceCapacityAllocated"`
}

type recordIdentity struct {
	RecordID   string    `xml:"sr:recordId,attr"`
	CreateTime time.Time `xml:"sr:createTime,attr"`
}

// NewStorageWrapper wraps given storage in a RecordWrapper
func NewStorageWrapper(st goat_grpc.StorageRecord) RecordWrapper {
	return storageWrapper{
		st: st,
	}
}

func (sw storageWrapper) Filename() string {
	return sw.st.GetRecordID()
}

func (sw storageWrapper) AsXML() (interface{}, error) {
	return storageXML{
		RecordID: recordIdentity{
			RecordID:   sw.st.GetRecordID(),
			CreateTime: time.Unix(sw.st.GetCreateTime().Seconds, 0),
		},
		StorageSystem:             sw.st.GetStorageSystem(),
		Site:                      s(sw.st.GetSite()),
		StorageShare:              s(sw.st.GetStorageShare()),
		StorageMedia:              s(sw.st.GetStorageMedia()),
		StorageClass:              s(sw.st.GetStorageClass()),
		FileCount:                 s(sw.st.GetFileCount()),
		DirectoryPath:             s(sw.st.GetDirectoryPath()),
		LocalUser:                 s(sw.st.GetLocalUser()),
		LocalGroup:                s(sw.st.GetLocalGroup()),
		UserIdentity:              s(sw.st.GetUserIdentity()),
		Group:                     s(sw.st.GetGroup()),
		GroupAttribute:            s(sw.st.GetGroupAttribute()),
		GroupAttributeType:        s(sw.st.GetGroupAttributeType()),
		StartTime:                 time.Unix(sw.st.GetStartTime().Seconds, 0),
		EndTime:                   time.Unix(sw.st.GetEndTime().Seconds, 0),
		ResourceCapacityUsed:      sw.st.GetResourceCapacityUsed(),
		LogicalCapacityUsed:       u64(sw.st.GetLogicalCapacityUsed()),
		ResourceCapacityAllocated: u64(sw.st.GetResourceCapacityAllocated()),
	}, nil
}

func (sw storageWrapper) AsJSON() (interface{}, error) {
	return nil, ErrNotImplemented
}

func (sw storageWrapper) AsTemplate() (interface{}, error) {
	return nil, ErrNotImplemented

}
