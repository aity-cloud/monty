package oci_test

import (
	"fmt"

	"github.com/aity-cloud/monty/pkg/oci"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/go-digest"
)

const (
	testSHAValue = "sha256:15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225"
)

var _ = Describe("OCI", Label("unit"), func() {
	When("the image string has a localhost server", func() {
		It("should parse the image correctly", func() {
			image := "localhost/monty/minimal:v1.0.0"
			parsed, err := oci.Parse(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsed.Registry).To(Equal("localhost"))
			Expect(parsed.Repository).To(Equal("monty/minimal"))
			Expect(parsed.Tag).To(Equal("v1.0.0"))
		})
	})
	When("image string does not contain registry", func() {
		When("image tag is present", func() {
			It("should parse the image correctly", func() {
				image := "monty/minimal:v1.0.0"
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal(""))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Tag).To(Equal("v1.0.0"))
			})
		})
		When("image digest is present", func() {
			It("should parse the image correctly", func() {
				image := fmt.Sprintf("monty/minimal@%s", testSHAValue)
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal(""))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Digest.String()).To(Equal(testSHAValue))
			})
		})
		When("neither tag nor digest is present", func() {
			It("should parse the image correctly", func() {
				image := "monty/minimal"
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal(""))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Digest.String()).To(Equal(""))
				Expect(parsed.Tag).To(Equal(""))
			})
		})
		When("both tag and digest are present", func() {
			It("should parse the image correctly", func() {
				image := fmt.Sprintf("monty/minimal:v1.0.0@%s", testSHAValue)
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal(""))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Digest.String()).To(Equal(testSHAValue))
				Expect(parsed.Tag).To(Equal("v1.0.0"))
			})
		})
	})
	When("image string contains registry", func() {
		When("image tag is present", func() {
			It("should parse the image correctly", func() {
				image := "quay.io/monty/minimal:v1.0.0"
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal("quay.io"))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Tag).To(Equal("v1.0.0"))
			})
		})
		When("image digest is present", func() {
			It("should parse the image correctly", func() {
				image := fmt.Sprintf("quay.io/monty/minimal@%s", testSHAValue)
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal("quay.io"))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Digest.String()).To(Equal(testSHAValue))
			})
		})
		When("neither tag nor digest is present", func() {
			It("should parse the image correctly", func() {
				image := "quay.io/monty/minimal"
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal("quay.io"))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Tag).To(Equal(""))
				Expect(parsed.Digest.String()).To(Equal(""))
			})
		})
		When("both tag and digest are present", func() {
			It("should parse the image correctly", func() {
				image := fmt.Sprintf("docker.io/monty/minimal:v1.0.0@%s", testSHAValue)
				parsed, err := oci.Parse(image)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsed.Registry).To(Equal("docker.io"))
				Expect(parsed.Repository).To(Equal("monty/minimal"))
				Expect(parsed.Digest.String()).To(Equal(testSHAValue))
				Expect(parsed.Tag).To(Equal("v1.0.0"))
			})
		})
	})
	When("image string contains registry and port", func() {
		It("should parse the image correctly", func() {
			image := "quay.io:5000/monty/minimal:v1.0.0"
			parsed, err := oci.Parse(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsed.Registry).To(Equal("quay.io:5000"))
			Expect(parsed.Repository).To(Equal("monty/minimal"))
			Expect(parsed.Tag).To(Equal("v1.0.0"))
		})
	})
	When("image string contains registry and a simple repository", func() {
		It("should parse the image correctly", func() {
			image := "quay.io/monty"
			parsed, err := oci.Parse(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsed.Registry).To(Equal("quay.io"))
			Expect(parsed.Repository).To(Equal("monty"))
			Expect(parsed.Tag).To(Equal(""))
			Expect(parsed.Digest.String()).To(Equal(""))
		})
	})
	When("image string contains a simple repository", func() {
		It("should parse the image correctly", func() {
			image := "monty"
			parsed, err := oci.Parse(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsed.Registry).To(Equal(""))
			Expect(parsed.Repository).To(Equal("monty"))
			Expect(parsed.Digest.String()).To(Equal(""))
			Expect(parsed.Tag).To(Equal(""))
		})
	})
	When("image string contains a simple repository and a tag", func() {
		It("should parse the image correctly", func() {
			image := "monty:v1.0.0"
			parsed, err := oci.Parse(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsed.Registry).To(Equal(""))
			Expect(parsed.Repository).To(Equal("monty"))
			Expect(parsed.Tag).To(Equal("v1.0.0"))
		})
	})

	When("image has a registry and no tag/digest", func() {
		image := &oci.Image{
			Registry:   "quay.io",
			Repository: "monty/minimal",
		}
		It("should return the correct image strings", func() {
			Expect(image.Path()).To(Equal("quay.io/monty/minimal"))
			Expect(image.String()).To(Equal("quay.io/monty/minimal"))
		})
	})
	When("image has a registry and a tag", func() {
		image := &oci.Image{
			Registry:   "quay.io",
			Repository: "monty/minimal",
			Tag:        "v1.0.0",
		}
		It("should return the correct image strings", func() {
			Expect(image.Path()).To(Equal("quay.io/monty/minimal"))
			Expect(image.String()).To(Equal("quay.io/monty/minimal:v1.0.0"))
		})
	})
	When("image has a registry and a digest", func() {
		var image *oci.Image
		BeforeEach(func() {
			digest, err := digest.Parse(testSHAValue)
			Expect(err).ToNot(HaveOccurred())
			image = &oci.Image{
				Registry:   "quay.io",
				Repository: "monty/minimal",
				Digest:     digest,
			}
		})

		It("should return the correct image strings", func() {
			Expect(image.Path()).To(Equal("quay.io/monty/minimal"))
			Expect(image.String()).To(Equal(fmt.Sprintf("quay.io/monty/minimal@%s", testSHAValue)))
		})
	})
	When("image has no registry and no tag/digest", func() {
		image := oci.Image{
			Registry:   "",
			Repository: "monty/minimal",
		}
		It("should return the correct image strings", func() {
			Expect(image.Path()).To(Equal("monty/minimal"))
			Expect(image.String()).To(Equal("monty/minimal"))
		})
	})
	When("image has no registry and a tag", func() {
		image := oci.Image{
			Registry:   "",
			Repository: "monty/minimal",
			Tag:        "v1.0.0",
		}
		It("should return the correct image strings", func() {
			Expect(image.Path()).To(Equal("monty/minimal"))
			Expect(image.String()).To(Equal("monty/minimal:v1.0.0"))
		})
	})
	When("image has a registry and a digest", func() {
		var image *oci.Image
		BeforeEach(func() {
			digest, err := digest.Parse(testSHAValue)
			Expect(err).ToNot(HaveOccurred())
			image = &oci.Image{
				Registry:   "",
				Repository: "monty/minimal",
				Digest:     digest,
			}
		})
		It("should return the correct image strings", func() {
			Expect(image.Path()).To(Equal("monty/minimal"))
			Expect(image.String()).To(Equal(fmt.Sprintf("monty/minimal@%s", testSHAValue)))
		})
	})

	When("updating the reference", func() {
		var ref string
		When("the image has a valid tag", func() {
			BeforeEach(func() {
				ref = "v0.1.1"
			})
			It("should successfuilly update the tag", func() {
				image := &oci.Image{}
				err := image.UpdateDigestOrTag(ref)
				Expect(err).ToNot(HaveOccurred())
				Expect(image.Tag).To(Equal(ref))
			})
		})
		When("the image has an invalid tag", func() {
			BeforeEach(func() {
				ref = "v0.1.1~foo"
			})
			It("should error", func() {
				image := &oci.Image{}
				err := image.UpdateDigestOrTag(ref)
				Expect(err).To(MatchError(oci.ErrInvalidReferenceFormat))
			})
		})
		When("the image has a valid digest", func() {
			BeforeEach(func() {
				ref = testSHAValue
			})
			It("should successfuilly update the tag", func() {
				image := &oci.Image{}
				err := image.UpdateDigestOrTag(ref)
				Expect(err).ToNot(HaveOccurred())
				Expect(image.Digest.String()).To(Equal(ref))
			})
		})
		When("the image has a valid digest", func() {
			BeforeEach(func() {
				ref = "sha256:123456789"
			})
			It("should error", func() {
				image := &oci.Image{}
				err := image.UpdateDigestOrTag(ref)
				Expect(err).To(MatchError(oci.ErrInvalidReferenceFormat))
			})
		})
		When("the ref is an empty string", func() {
			It("should error", func() {
				image := &oci.Image{}
				err := image.UpdateDigestOrTag(ref)
				Expect(err).To(MatchError(oci.ErrInvalidReferenceFormat))
			})
		})
	})
})
