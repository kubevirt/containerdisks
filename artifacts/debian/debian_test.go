package debian

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Debian", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, versionName, arch, mockFile string, details *api.ArtifactDetails,
			exampleUserData *docs.UserData, envVariables map[string]string, metadata *api.Metadata,
		) {
			c := New(release, versionName, arch, exampleUserData, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got.ChecksumHash).ToNot(BeNil())
			Expect(got.Checksum).To(Equal(details.Checksum))
			Expect(got.DownloadURL).To(Equal(details.DownloadURL))
			Expect(got.AdditionalUniqueTags).To(Equal(details.AdditionalUniqueTags))
			Expect(got.ImageArchitecture).To(Equal(details.ImageArchitecture))
			Expect(got.Compression).To(Equal(details.Compression))
			Expect(c.Metadata()).To(Equal(metadata))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("debian:11 x86_64", "11", "bullseye", "x86_64", "testdata/debian-11-genericcloud-amd64.json",
			&api.ArtifactDetails{
				Checksum: "3c08356d6860f987089c14b45953fb1f266d1b1b50dd086744925e2ed4113b804e848a8b1b46614febc48cd" +
					"e759f18e824b76bfb02618ed6b3d06ed15ea99283",
				DownloadURL:          "https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-genericcloud-amd64.qcow2",
				ImageArchitecture:    "amd64",
				AdditionalUniqueTags: []string{"11-20250303-2040"},
			},
			&docs.UserData{
				Username: "debian",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "debian",
			},
			&api.Metadata{
				Name:        "debian",
				Version:     "11",
				Arch:        "x86_64",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "debian",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "debian",
				},
			},
		),
		Entry("debian:11 aarch64", "11", "bullseye", "aarch64", "testdata/debian-11-genericcloud-arm64.json",
			&api.ArtifactDetails{
				Checksum: "c1a1645cf37ce628a8734bb25dce09fcd0858865302635ce0ae88b2da23bb615da43d483984709d743cd6b6" +
					"b45d56d88e9f6800f0b3110ba1b09c01b990342f3",
				DownloadURL:          "https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-genericcloud-arm64.qcow2",
				ImageArchitecture:    "arm64",
				AdditionalUniqueTags: []string{"11-20250303-2040"},
			},
			&docs.UserData{
				Username: "debian",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "debian",
			},
			&api.Metadata{
				Name:        "debian",
				Version:     "11",
				Arch:        "aarch64",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "debian",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "debian",
				},
			},
		),
		Entry("debian:12 x86_64", "12", "bookworm", "x86_64", "testdata/debian-12-genericcloud-amd64.json",
			&api.ArtifactDetails{
				Checksum: "a58d86525d75fd8e139a2302531ce5d2ab75ef0273cfe78f9d53aada4b23efd45f8433b4806fa4570cfe981" +
					"c8fae26f5e5e855cbd66ba2198862f28125fd2d45",
				DownloadURL:          "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2",
				ImageArchitecture:    "amd64",
				AdditionalUniqueTags: []string{"12-20250210-2019"},
			},
			&docs.UserData{
				Username: "debian",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "debian",
			},
			&api.Metadata{
				Name:        "debian",
				Version:     "12",
				Arch:        "x86_64",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "debian",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "debian",
				},
			},
		),
		Entry("debian:12 aarch64", "12", "bookworm", "aarch64", "testdata/debian-12-genericcloud-arm64.json",
			&api.ArtifactDetails{
				Checksum: "a17a462acbc3412ef195390fb60dffba2134fef1a276d500ca50a06036c488035657409fcd02f2f70d1e7a9" +
					"1776ca4249cfbceabeb90e74cb123b9971381c72a",
				DownloadURL:          "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-arm64.qcow2",
				ImageArchitecture:    "arm64",
				AdditionalUniqueTags: []string{"12-20250210-2019"},
			},
			&docs.UserData{
				Username: "debian",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "debian",
			},
			&api.Metadata{
				Name:        "debian",
				Version:     "12",
				Arch:        "aarch64",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "debian",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "debian",
				},
			},
		),

		Entry("debian:13 x86_64", "13", "trixie", "x86_64", "testdata/debian-13-genericcloud-amd64.json",
			&api.ArtifactDetails{
				Checksum: "d76122c87c940d1ab9334f4307c98c01dc42f0b49a20cddf278d59b92d34ab63d05ac1f40dffda3d2d32e38" +
					"1f097706eee6ccbf79a596bfb2cbb3d83c635ae35",
				DownloadURL:          "https://cloud.debian.org/images/cloud/trixie/latest/debian-13-genericcloud-amd64.qcow2",
				ImageArchitecture:    "amd64",
				AdditionalUniqueTags: []string{"13-20250814-2204"},
			},
			&docs.UserData{
				Username: "debian",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "debian",
			},
			&api.Metadata{
				Name:        "debian",
				Version:     "13",
				Arch:        "x86_64",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "debian",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "debian",
				},
			},
		),
		Entry("debian:13 aarch64", "13", "trixie", "aarch64", "testdata/debian-13-genericcloud-arm64.json",
			&api.ArtifactDetails{
				Checksum: "e36d98d9ee1f09fc9d7748f0f0f703e6424a2637f075de1aa9f06a4a58039e88086ee14bdaf" +
					"9483bbc259633248d28b89a32190da959a767f4732c088bdd30d0",
				DownloadURL:          "https://cloud.debian.org/images/cloud/trixie/latest/debian-13-genericcloud-arm64.qcow2",
				ImageArchitecture:    "arm64",
				AdditionalUniqueTags: []string{"13-20250814-2204"},
			},
			&docs.UserData{
				Username: "debian",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "debian",
			},
			&api.Metadata{
				Name:        "debian",
				Version:     "13",
				Arch:        "aarch64",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "debian",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "debian",
				},
			},
		),
	)
})

func TestDebian(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Debian Suite")
}
