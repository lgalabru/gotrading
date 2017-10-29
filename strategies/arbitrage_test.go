package strategies

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Arbitrage", func() {

	var (
		arbitrage Arbitrage
	)

	BeforeEach(func() {
		arbitrage = Arbitrage{}
	})

	Describe("Getting started with ginko", func() {
		Context("With 101 test", func() {
			It("should work", func() {
				Expect(1).To(Equal(1))
			})
		})
	})

})
