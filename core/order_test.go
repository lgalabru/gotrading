package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order", func() {

	var (
		order Order
	)

	BeforeEach(func() {
		order = Order{}
	})

	Describe("Getting started with ginko", func() {
		Context("With 101 test", func() {
			It("should work", func() {
				Expect(1).To(Equal(1))
			})
		})
	})

})
