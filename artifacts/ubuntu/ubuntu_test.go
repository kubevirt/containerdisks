package ubuntu

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Ubuntu", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, arch, mockFile string, details *api.ArtifactDetails, envVariables map[string]string, metadata *api.Metadata) {
			c := New(release, arch, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("ubuntu:22.04 x86_64", "22.04", "x86_64", "testdata/SHA256SUM",
			&api.ArtifactDetails{
				SHA256Sum:         "de5e632e17b8965f2baf4ea6d2b824788e154d9a65df4fd419ec4019898e15cd",
				DownloadURL:       "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img",
				ImageArchitecture: "amd64",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "ubuntu",
			},
			&api.Metadata{
				Name:        "ubuntu",
				Version:     "22.04",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "ubuntu",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "ubuntu",
				},
			},
		),
		Entry("ubuntu:22.04 aarch64", "22.04", "aarch64", "testdata/SHA256SUM",
			&api.ArtifactDetails{
				SHA256Sum:         "66224c7fed99ff5a5539eda406c87bbfefe8af6ff6b47d92df3187832b5b5d4f",
				DownloadURL:       "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-arm64.img",
				ImageArchitecture: "arm64",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "ubuntu",
			},
			&api.Metadata{
				Name:        "ubuntu",
				Version:     "22.04",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "ubuntu",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "ubuntu",
				},
			},
		),
		Entry("ubuntu:22.04 s390x", "22.04", "s390x", "testdata/SHA256SUM",
			&api.ArtifactDetails{
				SHA256Sum:         "192c18a58917622e12a3bb6aaf246fcc6a76d9562eb9f49d34df81fbc59610af",
				DownloadURL:       "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-s390x.img",
				ImageArchitecture: "s390x",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "ubuntu",
			},
			&api.Metadata{
				Name:        "ubuntu",
				Version:     "22.04",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "ubuntu",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "ubuntu",
				},
			},
		),
	)
})

func TestUbuntu(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ubuntu Suite")
}
