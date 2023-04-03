package rhcos

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Rhcos", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, mockFile string, details *api.ArtifactDetails, metadata *api.Metadata) {
			c := New(release, true)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("rhcos:4.9", "4.9", "testdata/rhcos-4.9.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8",
				DownloadURL:          "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.9/latest/rhcos-openstack.x86_64.qcow2.gz",
				Compression:          "gzip",
				AdditionalUniqueTags: []string{"3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8"},
			},
			&api.Metadata{
				Name:                   "rhcos",
				Version:                "4.9",
				Description:            description,
				ExampleUserDataPayload: docs.Ignition(&docs.UserData{}),
			},
		),
		Entry("rhcos:4.8", "4.8", "testdata/rhcos-4.8.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "99da4ed945b391d452e55a3a7809c799e4c74f69acbee1ecaec78f368c4e369e",
				DownloadURL:          "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.8/latest/rhcos-openstack.x86_64.qcow2.gz",
				Compression:          "gzip",
				AdditionalUniqueTags: []string{"99da4ed945b391d452e55a3a7809c799e4c74f69acbee1ecaec78f368c4e369e"},
			},
			&api.Metadata{
				Name:                   "rhcos",
				Version:                "4.8",
				Description:            description,
				ExampleUserDataPayload: docs.Ignition(&docs.UserData{}),
			},
		),
	)
})

func TestRhcos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rhcos Suite")
}
