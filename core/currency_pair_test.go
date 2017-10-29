package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CurrencyPair", func() {

	var (
		currencyPair CurrencyPair
	)

	BeforeEach(func() {
		currencyPair = CurrencyPair{}
	})

	Describe("Getting started with ginko", func() {
		Context("With 101 test", func() {
			It("should work", func() {
				Expect(1).To(Equal(1))
			})
		})
	})
})
