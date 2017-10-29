package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExchangeMashup", func() {

	var (
		mashup ExchangeMashup
	)

	BeforeEach(func() {
		mashup = ExchangeMashup{}
	})

	Describe("Getting started with ginko", func() {
		Context("With 101 test", func() {
      It("should work", func() {
				Expect(1).To(Equal(1))
			})
		})
	})

})
