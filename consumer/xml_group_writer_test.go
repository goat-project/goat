package consumer

import (
	"bytes"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("XML group writer tests", func() {
	var (
		xgw  XMLGroupWriter
		hook *test.Hook
	)

	hook = test.NewGlobal()

	id1 := "1"
	id2 := "2"
	id3 := "3"

	rec1 := goat_grpc.StorageRecord{RecordID: id1, CreateTime: wrapTimestamp(int64(0)),
		StartTime: wrapTimestamp(int64(0)), EndTime: wrapTimestamp(int64(0))}
	rec2 := goat_grpc.StorageRecord{RecordID: id2, CreateTime: wrapTimestamp(int64(0)),
		StartTime: wrapTimestamp(int64(0)), EndTime: wrapTimestamp(int64(0))}
	rec3 := goat_grpc.StorageRecord{RecordID: id3, CreateTime: wrapTimestamp(int64(0)),
		StartTime: wrapTimestamp(int64(0)), EndTime: wrapTimestamp(int64(0))}

	length := 5

	testRecords := make([]interface{}, length)
	testRecords[0] = rec1
	testRecords[1] = rec2
	testRecords[2] = rec3

	JustAfterEach(func() {
		xgw = XMLGroupWriter{}
	})

	Describe("outputDir", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					outDir: testOutputDirPath,
				}
			})

			It("should return non-empty string", func() {
				Expect(xgw.outputDir()).To(Equal(testOutputDirPath))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty string", func() {
				Expect(xgw.outputDir()).To(BeEmpty())
			})
		})
	})

	Describe("countPerFile", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					count: testCountPerFile,
				}
			})

			It("should return non-empty value", func() {
				Expect(xgw.countPerFile()).To(Equal(testCountPerFile))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty int64", func() {
				Expect(xgw.countPerFile()).To(Equal(uint64(0)))
			})
		})
	})

	Describe("records", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					recs: testRecords,
				}
			})

			It("should return non-empty slice", func() {
				Expect(xgw.records()).To(HaveLen(length))

				x := xgw.records()[0].(goat_grpc.StorageRecord)
				y := xgw.records()[1].(goat_grpc.StorageRecord)
				z := xgw.records()[2].(goat_grpc.StorageRecord)

				Expect(x.GetRecordID()).To(Equal(id1))
				Expect(y.GetRecordID()).To(Equal(id2))
				Expect(z.GetRecordID()).To(Equal(id3))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty slice", func() {
				Expect(xgw.records()).To(BeEmpty())
			})
		})
	})

	Describe("save", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					recs: testRecords,
				}
			})

			It("should return non-empty slice", func() {
				xgw.save(rec1, 4)

				Expect(xgw.records()).To(HaveLen(length))

				x := xgw.records()[4].(goat_grpc.StorageRecord)

				Expect(x.GetRecordID()).To(Equal(id1))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should log error", func() {
				xgw.save(rec1, 2)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("index out of range"))
			})
		})

		Context("when index is out of range", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					recs: testRecords,
				}
			})

			It("should log error", func() {
				xgw.save(rec1, 5)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("index out of range"))
			})
		})
	})

	Describe("wrap", func() {
		Context("when incoming data is correct", func() {
			It("should NOT return an error", func() {
				_, err := xgw.wrap(wrapper.WrapStorage(rec1))

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when incoming data is correct", func() {
			It("should return an error", func() {
				_, err := xgw.wrap(wrapper.WrapIP(goat_grpc.IpRecord{}))

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

					wrappedTestRecords[i], err = xgw.wrap(wrapper.WrapStorage(v.(goat_grpc.StorageRecord)))
					if err != nil {
						continue
					}
				}

				xgw = XMLGroupWriter{
					recs: wrappedTestRecords,
				}
			})

			It("should not return an error", func() {
				var b bytes.Buffer

				Expect(xgw.convertAndWrite(&b, uint64(3))).NotTo(HaveOccurred())
			})
		})

		Context("when buffer is nil", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					recs: testRecords,
				}
			})

			It("should return an return", func() {
				Expect(xgw.convertAndWrite(nil, uint64(2))).To(HaveOccurred())
			})
		})

		Context("when no records to write", func() {
			BeforeEach(func() {
				xgw = XMLGroupWriter{
					recs: nil,
				}
			})

			It("should not return an return", func() {
				var b bytes.Buffer

				Expect(xgw.convertAndWrite(&b, uint64(1))).NotTo(HaveOccurred())
			})
		})
	})
})
