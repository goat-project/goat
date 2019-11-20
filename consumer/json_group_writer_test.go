package consumer

import (
	"bytes"
	"time"

	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("JSON group writer tests", func() {
	var (
		jgw  JSONGroupWriter
		hook *test.Hook
	)

	hook = test.NewGlobal()

	user1 := "user1"
	user2 := "user2"
	user3 := "user3"

	now := time.Now().Unix()

	rec1 := goat_grpc.IpRecord{LocalUser: user1, MeasurementTime: wrapTimestamp(now)}
	rec2 := goat_grpc.IpRecord{LocalUser: user2, MeasurementTime: wrapTimestamp(now)}
	rec3 := goat_grpc.IpRecord{LocalUser: user3, MeasurementTime: wrapTimestamp(now)}

	length := 5

	testRecords := make([]interface{}, length)
	testRecords[0] = rec1
	testRecords[1] = rec2
	testRecords[2] = rec3

	JustAfterEach(func() {
		jgw = JSONGroupWriter{}
	})

	Describe("outputDir", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					outDir: testOutputDirPath,
				}
			})

			It("should return non-empty string", func() {
				Expect(jgw.outputDir()).To(Equal(testOutputDirPath))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty string", func() {
				Expect(jgw.outputDir()).To(BeEmpty())
			})
		})
	})

	Describe("countPerFile", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					count: testCountPerFile,
				}
			})

			It("should return non-empty value", func() {
				Expect(jgw.countPerFile()).To(Equal(testCountPerFile))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty int64", func() {
				Expect(jgw.countPerFile()).To(Equal(uint64(0)))
			})
		})
	})

	Describe("records", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					recs: testRecords,
				}
			})

			It("should return non-empty slice", func() {
				Expect(jgw.records()).To(HaveLen(length))

				x := jgw.records()[0].(goat_grpc.IpRecord)
				y := jgw.records()[1].(goat_grpc.IpRecord)
				z := jgw.records()[2].(goat_grpc.IpRecord)

				Expect(x.GetLocalUser()).To(Equal(user1))
				Expect(y.GetLocalUser()).To(Equal(user2))
				Expect(z.GetLocalUser()).To(Equal(user3))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should return empty slice", func() {
				Expect(jgw.records()).To(BeEmpty())
			})
		})
	})

	Describe("save", func() {
		Context("when structure is set correctly", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					recs: testRecords,
				}
			})

			It("should return non-empty slice", func() {
				jgw.save(rec1, 4)

				Expect(jgw.records()).To(HaveLen(length))

				x := jgw.records()[4].(goat_grpc.IpRecord)

				Expect(x.GetLocalUser()).To(Equal(user1))
			})
		})

		Context("when structure is NOT set correctly", func() {
			It("should log error", func() {
				jgw.save(rec1, 2)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("index out of range"))
			})
		})

		Context("when index is out of range", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					recs: testRecords,
				}
			})

			It("should log error", func() {
				jgw.save(rec1, 5)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("index out of range"))
			})
		})
	})

	Describe("wrap", func() {
		Context("when incoming data is correct", func() {
			It("should NOT return an error", func() {
				_, err := jgw.wrap(wrapper.WrapIP(rec1))

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when incoming data is correct", func() {
			It("should return an error", func() {
				_, err := jgw.wrap(wrapper.WrapVM(goat_grpc.VmRecord{}))

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

					wrappedTestRecords[i], err = jgw.wrap(wrapper.WrapIP(v.(goat_grpc.IpRecord)))
					if err != nil {
						continue
					}
				}

				jgw = JSONGroupWriter{
					recs: wrappedTestRecords,
				}
			})

			It("should not return an error", func() {
				var b bytes.Buffer

				Expect(jgw.convertAndWrite(&b, uint64(3))).NotTo(HaveOccurred())
			})
		})

		Context("when buffer is nil", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					recs: testRecords,
				}
			})

			It("should return an return", func() {
				Expect(jgw.convertAndWrite(nil, uint64(2))).To(HaveOccurred())
			})
		})

		Context("when no records to write", func() {
			BeforeEach(func() {
				jgw = JSONGroupWriter{
					recs: nil,
				}
			})

			It("should not return an return", func() {
				var b bytes.Buffer

				Expect(jgw.convertAndWrite(&b, uint64(1))).NotTo(HaveOccurred())
			})
		})
	})
})
