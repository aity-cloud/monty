package urn_test

import (
	"github.com/aity-cloud/monty/pkg/urn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("URN", Label("unit"), func() {
	When("URN doesn't start with urn", func() {
		It("should return error", func() {
			_, err := urn.ParseString("noturn:foo:bar:baz:bat")
			Expect(err).To(MatchError(urn.ErrInvalidURN))
		})
	})
	When("URN doesn't have 5 parts", func() {
		It("should return error", func() {
			_, err := urn.ParseString("urn:foo:bar:baz")
			Expect(err).To(MatchError(urn.ErrInvalidURN))
			_, err = urn.ParseString("urn:foo:bar:baz:bat:qux")
			Expect(err).To(MatchError(urn.ErrInvalidURN))
		})
	})
	When("URN namespace is not monty", func() {
		It("should return error", func() {
			_, err := urn.ParseString("urn:foo:bar:baz:bat")
			Expect(err).To(MatchError(urn.ErrInvalidURN))
			Expect(err).To(MatchError(ContainSubstring("invalid namespace: foo")))
		})
	})
	When("URN is valid", func() {
		It("should parse successfully", func() {
			u, err := urn.ParseString("urn:monty:plugin:foo:bar")
			Expect(err).NotTo(HaveOccurred())
			Expect(u.Namespace).To(Equal("monty"))
			Expect(u.Type).To(Equal(urn.Plugin))
			Expect(u.Strategy).To(Equal("foo"))
			Expect(u.Component).To(Equal("bar"))
		})
	})
	Context("String construction", func() {
		Specify("should return a correct string", func() {
			u := urn.NewMontyURN(urn.Agent, "foo", "bar")
			Expect(u.String()).To(Equal("urn:monty:agent:foo:bar"))
		})
	})
})
