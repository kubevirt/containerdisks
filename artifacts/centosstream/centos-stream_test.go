package centosstream

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("CentosStream", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, mockFile string, details *api.ArtifactDetails, metadata *api.Metadata) {
			c := New(release)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("centos-stream:8", "8", "testdata/centos-stream8.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "8e22e67687b81e38c7212fc30c47cb24cbc4935c0f2459ed139f498397d1e7cd",
				DownloadURL:          "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-20210603.0.x86_64.qcow2",
				AdditionalUniqueTags: []string{"8-20210603.0"},
			},
			&api.Metadata{
				Name:                   "centos-stream",
				Version:                "8",
				Description:            description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{}),
			},
		),
		Entry("centos-stream:9", "9", "testdata/centos-stream9.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "bcebdc00511d6e18782732570056cfbc7cba318302748bfc8f66be9c0db68142",
				DownloadURL:          "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20211222.0.x86_64.qcow2",
				AdditionalUniqueTags: []string{"9-20211222.0"},
			},
			&api.Metadata{
				Name:                   "centos-stream",
				Version:                "9",
				Description:            description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{}),
			},
		),
	)
})

func TestCentosStream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CentosStream Suite")
}
