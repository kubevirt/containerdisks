package centos

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Centos", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, mockFile string, details *api.ArtifactDetails, metadata *api.Metadata) {
			c := New(release)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("centos:8.4", "8.4", "testdata/centos8.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "3510fc7deb3e1939dbf3fe6f65a02ab1efcc763480bc352e4c06eca2e4f7c2a2",
				DownloadURL:          "https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-GenericCloud-8.4.2105-20210603.0.x86_64.qcow2",
				Compression:          "",
				AdditionalUniqueTags: []string{"8.4.2105-20210603.0", "8.4.2105"},
			},
			&api.Metadata{
				Name:        "centos",
				Version:     "8.4",
				Description: description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{
					Username: "centos",
				}),
			},
		),
		Entry("centos:8.3", "8.3", "testdata/centos8.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "7ec97062618dc0a7ebf211864abf63629da1f325578868579ee70c495bed3ba0",
				DownloadURL:          "https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-GenericCloud-8.3.2011-20201204.2.x86_64.qcow2",
				Compression:          "",
				AdditionalUniqueTags: []string{"8.3.2011-20201204.2", "8.3.2011"},
			},
			&api.Metadata{
				Name:        "centos",
				Version:     "8.3",
				Description: description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{
					Username: "centos",
				}),
			},
		),
		Entry("centos:7-2009", "7-2009", "testdata/centos7.checksum",
			&api.ArtifactDetails{
				SHA256Sum:   "e38bab0475cc6d004d2e17015969c659e5a308111851b0e2715e84646035bdd3",
				DownloadURL: "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2009.qcow2",
			},
			&api.Metadata{
				Name:        "centos",
				Version:     "7-2009",
				Description: description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{
					Username: "centos",
				}),
			},
		),
		Entry("centos:7-1809", "7-1809", "testdata/centos7.checksum",
			&api.ArtifactDetails{
				SHA256Sum:   "42c062df8a8c36991ec0282009dd52ac488461a3f7ee114fc21a765bfc2671c2",
				DownloadURL: "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-1809.qcow2",
			},
			&api.Metadata{
				Name:        "centos",
				Version:     "7-1809",
				Description: description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{
					Username: "centos",
				}),
			},
		),
	)
})

func TestCentos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Centos Suite")
}
