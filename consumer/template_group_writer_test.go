package consumer

import (
	"bytes"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/google/uuid"
)

var testOutputDirPath = "test-output-directory"
var testCountPerFile = uint64(1)

var _ = Describe("Template group writer tests", func() {
	var (
		tgw  TemplateGroupWriter
		hook *test.Hook
	)

	hook = test.NewGlobal()

	testTemplatesDir := "../templates"

	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	uuid3 := uuid.New().String()

	rec1 := testVMRecord(uuid1)
	rec2 := testVMRecord(uuid2)
	rec3 := testVMRecord(uuid3)

	length := 5

	testRecords := make([]interface{}, length)
	testRecords[0] = rec1
	testRecords[1] = rec2
	testRecords[2] = rec3

	JustAfterEach(func() {
		tgw = TemplateGroupWriter{}
	})

	Describe("outputDir", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					outDir: testOutputDirPath,
				}
			})

			It("should return non-empty string", func() {
				Expect(tgw.outputDir()).To(Equal(testOutputDirPath))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty string", func() {
				Expect(tgw.outputDir()).To(BeEmpty())
			})
		})
	})

	Describe("countPerFile", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					count: testCountPerFile,
				}
			})

			It("should return non-empty value", func() {
				Expect(tgw.countPerFile()).To(Equal(testCountPerFile))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty int64", func() {
				Expect(tgw.countPerFile()).To(Equal(uint64(0)))
			})
		})
	})

	Describe("records", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					recs: testRecords,
				}
			})

			It("should return non-empty slice", func() {
				Expect(tgw.records()).To(HaveLen(length))

				x := tgw.records()[0].(goat_grpc.VmRecord)
				y := tgw.records()[1].(goat_grpc.VmRecord)
				z := tgw.records()[2].(goat_grpc.VmRecord)

				Expect(x.GetVmUuid()).To(Equal(uuid1))
				Expect(y.GetVmUuid()).To(Equal(uuid2))
				Expect(z.GetVmUuid()).To(Equal(uuid3))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty slice", func() {
				Expect(tgw.records()).To(BeEmpty())
			})
		})
	})

	Describe("save", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					recs: testRecords,
				}
			})

			It("should return non-empty slice", func() {
				tgw.save(rec1, 4)

				Expect(tgw.records()).To(HaveLen(length))

				x := tgw.records()[4].(goat_grpc.VmRecord)

				Expect(x.GetVmUuid()).To(Equal(uuid1))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should log error", func() {
				tgw.save(rec1, 2)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("index out of range"))
			})
		})

		Context("when index is out of range", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					recs: testRecords,
				}
			})

			It("should log error", func() {
				tgw.save(rec1, 5)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("index out of range"))
			})
		})
	})

	Describe("wrap", func() {
		Context("when incoming data is correct", func() {
			It("should NOT return an error", func() {
				_, err := tgw.wrap(wrapper.WrapVM(rec1))

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when incoming data is correct", func() {
			It("should return an error", func() {
				_, err := tgw.wrap(wrapper.WrapIP(goat_grpc.IpRecord{}))

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("convertAndWrite", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				wrappedTestRecords := make([]interface{}, length)
				var err error

				for i, v := range testRecords {
					if v == nil {
						continue
					}

					wrappedTestRecords[i], err = tgw.wrap(wrapper.WrapVM(v.(goat_grpc.VmRecord)))
					if err != nil {
						continue
					}
				}

				tgw = TemplateGroupWriter{
					recs:         wrappedTestRecords,
					templatesDir: testTemplatesDir,
				}
			})

			It("should not return an error", func() {
				var b bytes.Buffer

				Expect(tgw.convertAndWrite(&b, uint64(3))).NotTo(HaveOccurred())
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return an return", func() {
				var b bytes.Buffer

				Expect(tgw.convertAndWrite(&b, uint64(2))).To(HaveOccurred())
			})
		})

		Context("when buffer is nil", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					recs:         testRecords,
					templatesDir: testTemplatesDir,
				}
			})

			It("should return an return", func() {
				Expect(tgw.convertAndWrite(nil, uint64(1))).To(HaveOccurred())
			})
		})

		Context("when no records to write", func() {
			BeforeEach(func() {
				tgw = TemplateGroupWriter{
					recs:         nil,
					templatesDir: testTemplatesDir,
				}
			})

			It("should not return an return", func() {
				var b bytes.Buffer

				Expect(tgw.convertAndWrite(&b, uint64(3))).NotTo(HaveOccurred())
			})
		})
	})
})

func testVMRecord(uuid string) goat_grpc.VmRecord {
	return goat_grpc.VmRecord{
		VmUuid:              uuid,
		SiteName:            "test-site-name",
		CloudComputeService: wrapStringValue("test-cloud-compute-service"),
		MachineName:         "machine-name",
		LocalUserId:         wrapStringValue("test-local-user-id"),
		LocalGroupId:        wrapStringValue("test-local-group-id"),
		GlobalUserName:      wrapStringValue("test-global-user-name"),
		Fqan:                wrapStringValue("test-fqan"),
		Status:              wrapStringValue("test-status"),
		StartTime:           wrapTimestamp(int64(1573746846)),
		EndTime:             wrapTimestamp(int64(1573746846)),
		SuspendDuration:     wrapDuration(int64(1573746846)),
		WallDuration:        wrapDuration(int64(1573746846)),
		CpuDuration:         wrapDuration(int64(1573746846)),
		CpuCount:            uint32(4),
		NetworkType:         wrapStringValue("test-network-type"),
		NetworkInbound:      wrapUint64Value(uint64(1024)),
		NetworkOutbound:     wrapUint64Value(uint64(1024)),
		PublicIpCount:       wrapUint64Value(uint64(8)),
		Memory:              wrapUint64Value(uint64(1024)),
		Disk:                wrapUint64Value(uint64(4)),
		BenchmarkType:       wrapStringValue("test-benchmark-type"),
		Benchmark:           wrapFloatValue(float32(1024)),
		StorageRecordId:     nil,
		ImageId:             wrapStringValue("test-image-id"),
		CloudType:           wrapStringValue("test-cloud-type"),
	}
}

func wrapStringValue(s string) *wrappers.StringValue {
	return &wrappers.StringValue{Value: s}
}

func wrapTimestamp(s int64) *timestamp.Timestamp {
	return &timestamp.Timestamp{Seconds: s}
}

func wrapDuration(s int64) *duration.Duration {
	return &duration.Duration{Seconds: s}
}

func wrapUint64Value(v uint64) *wrappers.UInt64Value {
	return &wrappers.UInt64Value{Value: v}
}

func wrapFloatValue(f float32) *wrappers.FloatValue {
	return &wrappers.FloatValue{Value: f}
}
