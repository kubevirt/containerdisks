package ubuntu

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/api/instancetype"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Ubuntu", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, mockFile string, details *api.ArtifactDetails, additionalLabels map[string]string, metadata *api.Metadata) {
			c := New(release, additionalLabels)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("ubuntu:22.04", "22.04", "testdata/SHA256SUM",
			&api.ArtifactDetails{
				SHA256Sum:   "de5e632e17b8965f2baf4ea6d2b824788e154d9a65df4fd419ec4019898e15cd",
				DownloadURL: "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img",
			},
			map[string]string{
				instancetype.DefaultInstancetypeLabel: "u1.small",
				instancetype.DefaultPreferenceLabel:   "ubuntu",
			},
			&api.Metadata{
				Name:        "ubuntu",
				Version:     "22.04",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "ubuntu",
				},
				AdditionalLabels: map[string]string{
					instancetype.DefaultInstancetypeLabel: "u1.small",
					instancetype.DefaultPreferenceLabel:   "ubuntu",
				},
			},
		),
	)
})

func TestUbuntu(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ubuntu Suite")
}
