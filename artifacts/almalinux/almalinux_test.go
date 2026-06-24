package almalinux

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("AlmaLinux", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, arch, mockFile string, details *api.ArtifactDetails,
			exampleUserData *docs.UserData, envVariables map[string]string, metadata *api.Metadata,
		) {
			a := New(release, arch, exampleUserData, envVariables)
			a.getter = testutil.NewMockGetter(mockFile)
			got, err := a.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got.ChecksumHash).ToNot(BeNil())
			Expect(got.Checksum).To(Equal(details.Checksum))
			Expect(got.DownloadURL).To(Equal(details.DownloadURL))
			Expect(got.AdditionalUniqueTags).To(Equal(details.AdditionalUniqueTags))
			Expect(got.ImageArchitecture).To(Equal(details.ImageArchitecture))
			Expect(got.Compression).To(Equal(details.Compression))
			Expect(a.Metadata()).To(Equal(metadata))
		},
		Entry("almalinux:9 x86_64", "9", "x86_64", "testdata/almalinux9-x86_64.checksum",
			&api.ArtifactDetails{
				Checksum:             "c397eed7023e92c841155831b1f47e26300e5bef0f0256c129322307c897a251",
				DownloadURL:          "https://repo.almalinux.org/almalinux/9/cloud/x86_64/images/AlmaLinux-9-GenericCloud-9.8-20260526.x86_64.qcow2",
				AdditionalUniqueTags: []string{"9.8-20260526"},
				ImageArchitecture:    "amd64",
			},
			&docs.UserData{
				Username: "almalinux",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "rhel.9",
			},
			&api.Metadata{
				Name:        "almalinux",
				Version:     "9",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "almalinux",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "rhel.9",
				},
				Arch: "x86_64",
			},
		),
		Entry("almalinux:9 aarch64", "9", "aarch64", "testdata/almalinux9-aarch64.checksum",
			&api.ArtifactDetails{
				Checksum:             "b5d883c5f84c68a9828fbd3aac863f9a723b43f8965a32ff7a1c198301f42a29",
				DownloadURL:          "https://repo.almalinux.org/almalinux/9/cloud/aarch64/images/AlmaLinux-9-GenericCloud-9.8-20260526.aarch64.qcow2",
				AdditionalUniqueTags: []string{"9.8-20260526"},
				ImageArchitecture:    "arm64",
			},
			&docs.UserData{
				Username: "almalinux",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "rhel.9",
			},
			&api.Metadata{
				Name:        "almalinux",
				Version:     "9",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "almalinux",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "rhel.9",
				},
				Arch: "aarch64",
			},
		),
		Entry("almalinux:9 s390x", "9", "s390x", "testdata/almalinux9-s390x.checksum",
			&api.ArtifactDetails{
				Checksum:             "772eacf66540673b947b927e7f2c00aa1a9697d3416f3964291379976e8e3b76",
				DownloadURL:          "https://repo.almalinux.org/almalinux/9/cloud/s390x/images/AlmaLinux-9-GenericCloud-9.8-20260526.s390x.qcow2",
				AdditionalUniqueTags: []string{"9.8-20260526"},
				ImageArchitecture:    "s390x",
			},
			&docs.UserData{
				Username: "almalinux",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "rhel.9",
			},
			&api.Metadata{
				Name:        "almalinux",
				Version:     "9",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "almalinux",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "rhel.9",
				},
				Arch: "s390x",
			},
		),
		Entry("almalinux:10 x86_64", "10", "x86_64", "testdata/almalinux10-x86_64.checksum",
			&api.ArtifactDetails{
				Checksum:             "47f2218668dd4776be140dd92fa3bea700be1766e2c7d88bdfd6a4b50f477b4d",
				DownloadURL:          "https://repo.almalinux.org/almalinux/10/cloud/x86_64/images/AlmaLinux-10-GenericCloud-10.2-20260526.0.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"10.2-20260526.0"},
				ImageArchitecture:    "amd64",
			},
			&docs.UserData{
				Username: "almalinux",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "rhel.10",
			},
			&api.Metadata{
				Name:        "almalinux",
				Version:     "10",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "almalinux",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "rhel.10",
				},
				Arch: "x86_64",
			},
		),
		Entry("almalinux:10 aarch64", "10", "aarch64", "testdata/almalinux10-aarch64.checksum",
			&api.ArtifactDetails{
				Checksum:             "336c6861fe0ba9115af00f557ed1b09385a3525612dd0cb9cce7e0486f8e74a4",
				DownloadURL:          "https://repo.almalinux.org/almalinux/10/cloud/aarch64/images/AlmaLinux-10-GenericCloud-10.2-20260526.0.aarch64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"10.2-20260526.0"},
				ImageArchitecture:    "arm64",
			},
			&docs.UserData{
				Username: "almalinux",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "rhel.10",
			},
			&api.Metadata{
				Name:        "almalinux",
				Version:     "10",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "almalinux",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "rhel.10",
				},
				Arch: "aarch64",
			},
		),
		Entry("almalinux:10 s390x", "10", "s390x", "testdata/almalinux10-s390x.checksum",
			&api.ArtifactDetails{
				Checksum:             "db1f7e8247e150e6da902da348f6894960bfd3590001207d0bf3732992d98ac7",
				DownloadURL:          "https://repo.almalinux.org/almalinux/10/cloud/s390x/images/AlmaLinux-10-GenericCloud-10.2-20260526.0.s390x.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"10.2-20260526.0"},
				ImageArchitecture:    "s390x",
			},
			&docs.UserData{
				Username: "almalinux",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "rhel.10",
			},
			&api.Metadata{
				Name:        "almalinux",
				Version:     "10",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "almalinux",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "rhel.10",
				},
				Arch: "s390x",
			},
		),
	)
})

func TestAlmaLinux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AlmaLinux Suite")
}
