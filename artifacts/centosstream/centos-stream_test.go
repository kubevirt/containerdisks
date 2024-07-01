package centosstream

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/instancetype"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("CentosStream", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, arch, mockFile string, details *api.ArtifactDetails,
			exampleUserData *docs.UserData, envVariables map[string]string, metadata *api.Metadata,
		) {
			c := New(release, arch, exampleUserData, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("centos-stream:8 x86_64", "8", "x86_64", "testdata/centos-stream8-x86_64.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "8e22e67687b81e38c7212fc30c47cb24cbc4935c0f2459ed139f498397d1e7cd",
				DownloadURL:          "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-20210603.0.x86_64.qcow2",
				AdditionalUniqueTags: []string{"8-20210603.0"},
				ImageArchitecture:    "amd64",
			},
			&docs.UserData{
				Username: "centos",
			},
			map[string]string{
				instancetype.DefaultInstancetypeEnv: "u1.medium",
				instancetype.DefaultPreferenceEnv:   "centos.stream8",
			},
			&api.Metadata{
				Name:        "centos-stream",
				Version:     "8",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "centos",
				},
				EnvVariables: map[string]string{
					instancetype.DefaultInstancetypeEnv: "u1.medium",
					instancetype.DefaultPreferenceEnv:   "centos.stream8",
				},
			},
		),
		Entry("centos-stream:8 aarch64", "8", "aarch64", "testdata/centos-stream8-aarch64.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "5879673f9f7a992e57548dcd99a1f3d0e30c7250a725b8e92f105ac489208ba0",
				DownloadURL:          "https://cloud.centos.org/centos/8-stream/aarch64/images/CentOS-Stream-GenericCloud-8-latest.aarch64.qcow2",
				AdditionalUniqueTags: []string{"8-latest"},
				ImageArchitecture:    "arm64",
			},
			&docs.UserData{
				Username: "centos",
			},
			map[string]string{
				instancetype.DefaultInstancetypeEnv: "u1.medium",
				instancetype.DefaultPreferenceEnv:   "centos.stream8",
			},
			&api.Metadata{
				Name:        "centos-stream",
				Version:     "8",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "centos",
				},
				EnvVariables: map[string]string{
					instancetype.DefaultInstancetypeEnv: "u1.medium",
					instancetype.DefaultPreferenceEnv:   "centos.stream8",
				},
			},
		),
		Entry("centos-stream:9 x86_64", "9", "x86_64", "testdata/centos-stream9-x86_64.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "bcebdc00511d6e18782732570056cfbc7cba318302748bfc8f66be9c0db68142",
				DownloadURL:          "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20211222.0.x86_64.qcow2",
				AdditionalUniqueTags: []string{"9-20211222.0"},
				ImageArchitecture:    "amd64",
			},
			&docs.UserData{
				Username: "cloud-user",
			},
			map[string]string{
				instancetype.DefaultInstancetypeEnv: "u1.medium",
				instancetype.DefaultPreferenceEnv:   "centos.stream9",
			},
			&api.Metadata{
				Name:        "centos-stream",
				Version:     "9",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "cloud-user",
				},
				EnvVariables: map[string]string{
					instancetype.DefaultInstancetypeEnv: "u1.medium",
					instancetype.DefaultPreferenceEnv:   "centos.stream9",
				},
			},
		),
		Entry("centos-stream:9 aarch64", "9", "aarch64", "testdata/centos-stream9-aarch64.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "66dd927b7aa643b18ad21a9368571c6ef57cc381b4febc8934397b137f14b995",
				DownloadURL:          "https://cloud.centos.org/centos/9-stream/aarch64/images/CentOS-Stream-GenericCloud-9-latest.aarch64.qcow2",
				AdditionalUniqueTags: []string{"9-latest"},
				ImageArchitecture:    "arm64",
			},
			&docs.UserData{
				Username: "cloud-user",
			},
			map[string]string{
				instancetype.DefaultInstancetypeEnv: "u1.medium",
				instancetype.DefaultPreferenceEnv:   "centos.stream9",
			},
			&api.Metadata{
				Name:        "centos-stream",
				Version:     "9",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "cloud-user",
				},
				EnvVariables: map[string]string{
					instancetype.DefaultInstancetypeEnv: "u1.medium",
					instancetype.DefaultPreferenceEnv:   "centos.stream9",
				},
			},
		),
	)
})

func TestCentosStream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CentosStream Suite")
}
