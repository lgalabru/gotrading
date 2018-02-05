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

	Describe("Testing Order", func() {
		Context("Considering orders with fees", func() {
			It("should return the correct volume out", func() {
				o := Order{}
				o.InitBid(0.00000561, 23.73927493)
				o.updateVolumesInOut()
				Expect(o.BaseVolumeOut).To(Equal(23.679926742675))
			})
		})
	})

})
