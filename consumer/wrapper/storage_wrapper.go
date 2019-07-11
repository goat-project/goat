package wrapper

import (
	"time"

	goat_grpc "github.com/goat-project/goat-proto-go"
)

type storageWrapper struct {
	st goat_grpc.StorageRecord
}

type storageXML struct {
	RecordID                  string    `xml:"RECORD_ID"`
	CreateTime                time.Time `xml:"CREATE_TIME"`
	StorageSystem             string    `xml:"STORAGE_SYSTEM"`
	Site                      *string   `xml:"SITE"`
	StorageShare              *string   `xml:"STORAGE_SHARE"`
	StorageMedia              *string   `xml:"STORAGE_MEDIA"`
	StorageClass              *string   `xml:"STORAGE_CLASS"`
	FileCount                 *string   `xml:"FILE_COUNT"`
	DirectoryPath             *string   `xml:"DIRECTORY_PATH"`
	LocalUser                 *string   `xml:"LOCAL_USER"`
	LocalGroup                *string   `xml:"LOCAL_GROUP"`
	UserIdentity              *string   `xml:"USER_IDENTITY"`
	Group                     *string   `xml:"GROUP"`
	GroupAttribute            *string   `xml:"GROUP_ATTRIBUTE"`
	GroupAttributeType        *string   `xml:"GROUP_ATTRIBUTE_TYPE"`
	StartTime                 time.Time `xml:"START_TIME"`
	EndTime                   time.Time `xml:"END_TIME"`
	ResourceCapacityUsed      uint64    `xml:"RESOURCE_CAPACITY_USED"`
	LogicalCapacityUsed       *uint64   `xml:"LOGICAL_CAPACITY_USED"`
	ResourceCapacityAllocated *uint64   `xml:"RESOURCE_CAPACITY_ALLOCATED"`
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
		RecordID:                  sw.st.GetRecordID(),
		CreateTime:                time.Unix(sw.st.GetCreateTime().Seconds, 0),
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
