package centos

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Centos", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, mockFile string, details *api.ArtifactDetails, envVariables map[string]string, metadata *api.Metadata) {
			c := New(release, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("centos:7-2009", "7-2009", "testdata/centos7.checksum",
			&api.ArtifactDetails{
				SHA256Sum:         "e38bab0475cc6d004d2e17015969c659e5a308111851b0e2715e84646035bdd3",
				DownloadURL:       "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2009.qcow2",
				ImageArchitecture: "amd64",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.small",
				common.DefaultPreferenceEnv:   "centos.7",
			},
			&api.Metadata{
				Name:        "centos",
				Version:     "7-2009",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "centos",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.small",
					common.DefaultPreferenceEnv:   "centos.7",
				},
			},
		),
		Entry("centos:7-1809", "7-1809", "testdata/centos7.checksum",
			&api.ArtifactDetails{
				SHA256Sum:         "42c062df8a8c36991ec0282009dd52ac488461a3f7ee114fc21a765bfc2671c2",
				DownloadURL:       "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-1809.qcow2",
				ImageArchitecture: "amd64",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.small",
				common.DefaultPreferenceEnv:   "centos.7",
			},
			&api.Metadata{
				Name:        "centos",
				Version:     "7-1809",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "centos",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.small",
					common.DefaultPreferenceEnv:   "centos.7",
				},
			},
		),
	)
})

func TestCentos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Centos Suite")
}
