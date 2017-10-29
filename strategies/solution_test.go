package strategies

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Solution", func() {

	var (
		solution Solution
	)

	BeforeEach(func() {
		solution = Solution{}
	})

	Describe("Getting started with ginko", func() {
		Context("With 101 test", func() {
			It("should work", func() {
				Expect(1).To(Equal(1))
			})
		})
	})

})
